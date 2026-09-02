# QR code payments — design notes

Covers issue [#7](https://github.com/darivnu/CartePro/issues/7) `POST /clients/me/qrcode`
and issue [#11](https://github.com/darivnu/CartePro/issues/11) `POST /partners/me/validate`.
Backend only — QR image rendering (turning `qr_payload` into an actual scannable
image) is a frontend concern (issue #43); the backend only ever deals with strings.

## Flow

1. Client app calls `POST /clients/me/qrcode` (authenticated as client).
   Backend creates a short-lived, single-use `QrToken` row and returns a signed
   `qr_payload` string.
2. Client app renders `qr_payload` as a QR code on screen (frontend).
3. Partner scans it, their terminal calls `POST /partners/me/validate`
   (authenticated as partner) with `{qr_payload, amount, idempotency_key}` —
   the amount is decided by the partner terminal, like a card reader, not
   embedded in the QR code.
4. Backend verifies the payload, atomically debits the client and credits the
   partner, creates a `Transaction`, and returns it.

## Existing schema (already migrated, no change needed)

```go
type QrToken struct {
    ID        uint
    ClientID  uint
    Token     string // unique
    CreatedAt time.Time
    ExpiresAt time.Time
    UsedAt    *time.Time // nil until redeemed
}

type Transaction struct {
    ID, ClientID, PartnerID, EmployerID, QrTokenID uint/*ptr*/
    Amount    int64 // cents
    Type      TransactionType // "debit" | "topup"
    CreatedAt time.Time
}
```

## Schema gap: idempotency

`Transaction` has no field to store `idempotency_key`, but issue #11 requires
one as input. Needs a migration:

```go
type Transaction struct {
    ...
    IdempotencyKey *string `gorm:"uniqueIndex"` // nil for topups/admin-created rows
}
```

Nullable + globally unique is enough: only the partner-transaction path sets
it, and it's expected to be a random string (UUID) generated per logical
purchase attempt by the partner terminal.

## `qr_payload` format

The issue explicitly says "signed", so don't just hand out the raw opaque
`QrToken.Token`. Sign a small payload with a server-side HMAC secret:

```
payload  = "<token>.<expires_unix>"
sig      = hex(HMAC-SHA256(payload, QR_SIGNING_SECRET))
qr_payload = base64url(payload) + "." + sig
```

- `token`: opaque random string, same generation pattern as session tokens
  (`crypto/rand`, 32 bytes, hex-encoded) — this is what actually gets stored
  in `QrToken.Token` and is the DB lookup key.
- On redemption, the partner endpoint recomputes the HMAC and rejects on
  mismatch **before** touching the DB (cheap tamper/garbage rejection).
- The DB row remains the source of truth for expiry and single-use, never
  the embedded `expires_unix` alone — a leaked/rotated secret or clock skew
  shouldn't be the only thing standing between an expired code and a payout.
- New env var: `QR_SIGNING_SECRET` (random ≥32 bytes), loaded the same way as
  the DB credentials (`godotenv`). New config constant: TTL, suggest 60–120s
  (recommend starting at 90s) — short enough to limit the replay window if a
  code is screenshotted, long enough for a human to actually present it.

The `token` field returned alongside `qr_payload` in the response is the raw
opaque token — useful for the client app to reference/poll its own pending
code later if ever needed; `qr_payload` is the only thing that goes into the
actual QR image.

## `POST /clients/me/qrcode`

Auth: session cookie, role must be `client` (same pattern as
`HandleClientBalance`).

1. Resolve `Client` from session.
2. Generate token, compute `expires_at = now + TTL`.
3. Insert `QrToken{ClientID, Token, ExpiresAt}`.
4. Sign payload as above.
5. `201 {"token", "qr_payload", "expires_at"}` (RFC3339 for the timestamp,
   matching the rest of the API).

Old unused/unexpired tokens for the same client are simply left alone —
each is independently single-use and short-lived, so allowing several
outstanding codes at once is harmless and avoids extra bookkeeping.

## `POST /partners/me/validate`

Auth: session cookie, role must be `partner`. Also gate on
`Partner.Status == StatusApproved` — a pending/rejected partner shouldn't be
able to take payments; this isn't spelled out in the issue but follows from
the existing `PartnerStatus` model.

1. Parse `{qr_payload, amount, idempotency_key}`; reject empty/`amount <= 0`.
2. **Idempotency check first**, before any mutation: look up
   `Transaction` by `(PartnerID, IdempotencyKey)`. If found, return that
   existing transaction (201) instead of reprocessing — this is what makes
   network retries from the partner terminal safe.
3. Verify HMAC signature on `qr_payload`, decode `token` + embedded expiry.
   Bad signature/malformed payload → 400.
4. Everything from here happens in a single DB transaction with row locks
   (`SELECT ... FOR UPDATE`, Postgres supports it, GORM via
   `clause.Locking{Strength: "UPDATE"}`) to prevent a double-scan race from
   redeeming the same code twice or racing another concurrent balance change:
   - Lock and load `QrToken` by `Token`.
     - Not found → 404.
     - `UsedAt != nil` → 409 (already used).
     - `ExpiresAt` in the past → 410 (expired).
   - Lock and load the `Client` referenced by `QrToken.ClientID`.
     - `Client.Balance < amount` → 402/400 (insufficient funds).
   - `Client.Balance -= amount`, `Partner.Balance += amount`.
   - `QrToken.UsedAt = now`.
   - Create `Transaction{ClientID, PartnerID, QrTokenID, Amount, Type: debit, IdempotencyKey}`.
   - Commit.
5. `201 {"transaction": {...}}`.

## Error cases summary

| Condition | Status |
|---|---|
| Not authenticated | 401 |
| Authenticated but not a partner / partner not approved | 403 |
| Malformed body / bad signature / amount <= 0 | 400 |
| Token doesn't exist | 404 |
| Token already used | 409 |
| Token expired | 410 |
| Insufficient client balance | 402 (or 400, pick one and stay consistent) |
| Idempotency key already seen (same partner) | 201, replay existing transaction |

## Implementation sketch

Token generation and signing (`crypto/rand` for the opaque token, same as
`generateSessionToken` in `server/session.go`; HMAC-SHA256 for the signature):

```go
func generateQrToken() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil // stored as QrToken.Token
}

func signQrPayload(token string, expiresAt time.Time) string {
    payload := fmt.Sprintf("%s.%d", token, expiresAt.Unix())
    mac := hmac.New(sha256.New, []byte(os.Getenv("QR_SIGNING_SECRET")))
    mac.Write([]byte(payload))
    sig := hex.EncodeToString(mac.Sum(nil))
    // base64 first so the embedded "." in payload doesn't collide with our separator
    return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + sig
}

func verifyQrPayload(qrPayload string) (token string, expiresAt time.Time, err error) {
    parts := strings.SplitN(qrPayload, ".", 2)
    if len(parts) != 2 {
        return "", time.Time{}, errors.New("malformed qr payload")
    }
    payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
    if err != nil {
        return "", time.Time{}, errors.New("malformed qr payload")
    }

    mac := hmac.New(sha256.New, []byte(os.Getenv("QR_SIGNING_SECRET")))
    mac.Write(payloadBytes)
    sigBytes, err := hex.DecodeString(parts[1])
    if err != nil || !hmac.Equal(mac.Sum(nil), sigBytes) { // constant-time compare
        return "", time.Time{}, errors.New("invalid signature")
    }

    payloadParts := strings.SplitN(string(payloadBytes), ".", 2)
    expiresUnix, err := strconv.ParseInt(payloadParts[1], 10, 64)
    if err != nil {
        return "", time.Time{}, errors.New("malformed qr payload")
    }
    return payloadParts[0], time.Unix(expiresUnix, 0), nil
}
```

Redemption — everything from the row lock onward runs inside one
`database.DB.Transaction(...)` closure, using `clause.Locking{Strength:
"UPDATE"}` (Postgres `SELECT ... FOR UPDATE`) so a second concurrent
redemption of the same code blocks on the lock and then fails the `UsedAt !=
nil` check instead of racing:

```go
err = database.DB.Transaction(func(tx *gorm.DB) error {
    var qrToken database.QrToken
    if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        Where("token = ?", token).First(&qrToken).Error; err != nil {
        return &httpError{http.StatusNotFound, "QR code not found"}
    }
    if qrToken.UsedAt != nil {
        return &httpError{http.StatusConflict, "QR code already used"}
    }
    if qrToken.ExpiresAt.Before(time.Now()) {
        return &httpError{http.StatusGone, "QR code expired"}
    }

    var client database.Client
    if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&client, qrToken.ClientID).Error; err != nil {
        return &httpError{http.StatusNotFound, "Client not found"}
    }
    if client.Balance < req.Amount {
        return &httpError{http.StatusPaymentRequired, "Insufficient balance"}
    }

    client.Balance -= req.Amount
    partner.Balance += req.Amount
    now := time.Now()
    qrToken.UsedAt = &now

    if err := tx.Save(&client).Error; err != nil { return err }
    if err := tx.Save(&partner).Error; err != nil { return err }
    if err := tx.Save(&qrToken).Error; err != nil { return err }

    transaction = database.Transaction{
        ClientID: client.ID, PartnerID: &partner.ID, QrTokenID: &qrToken.ID,
        Amount: req.Amount, Type: database.TransactionTypeDebit,
        IdempotencyKey: &req.IdempotencyKey,
    }
    return tx.Create(&transaction).Error
})
```

`httpError` is a small local `{status int; message string}` type implementing
`error`, used only so the closure can carry a specific HTTP status back out
through the handler's single `if err != nil` after `.Transaction(...)`
returns.

The idempotency lookup (`Where("idempotency_key = ?", req.IdempotencyKey)`)
must happen *before* this transaction even opens — it's what makes a retried
request from the partner terminal return the original transaction instead of
attempting to redeem an already-used token.

## Suggested file layout

- `backend/users/qrcode.go` — token generation, HMAC sign/verify helpers,
  `HandleClientQrCode`.
- `backend/users/transactions.go` — `HandlePartnerTransaction` (new file
  rather than piling onto `partners.go`, since it's a distinct concern with
  its own locking logic).
- Router additions in `backend/routing/router.go`:
  ```go
  mux.HandleFunc("POST /clients/me/qrcode", users.HandleClientQrCode)
  mux.HandleFunc("POST /partners/me/validate", users.HandlePartnerTransaction)
  ```

## Open decisions before implementing

- Exact TTL value (suggested 90s).
- 402 vs 400 for insufficient balance (repo has no precedent yet either way).
- Whether idempotency uniqueness should be scoped to `(PartnerID, IdempotencyKey)` or just global on `IdempotencyKey` — global is simpler and fine since the key is expected to be a random UUID anyway.
