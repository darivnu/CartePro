# Frontend architecture

Stack: React 19 + TypeScript, built with Vite. Tailwind CSS v4 for styling, shadcn/radix-ui for
primitive components, React Router v7 for routing, TanStack React Query v5 for server state.
Auth is a session cookie (`credentials: 'include'`), no tokens stored client-side.

## Directory layout (`frontend/src`)

```
api/          fetch wrappers, one file per backend resource (auth.ts, partners.ts, admin.ts, clients.ts)
              client.ts holds the shared apiFetch() + ApiError
auth/         useAuth.ts — the auth React Query hooks (useAuth, useLogin, useLogout)
partners/     usePartnersData.ts — React Query hooks wrapping api/partners.ts, used by client+partner routes
admin/        useAdminPartners.ts / useAdminClients.ts — same pattern, admin-only
client/       useClientData.ts — same pattern, client-only
routes/       one file per page/screen, plus RoleGuard.tsx and RootRedirect.tsx
  routes/admin/    admin screens, nested under AdminLayout
  routes/client/   client screens
  routes/partner/  partner screens
components/   shared presentational components (Wordmark, Logo, QrScanner, SimulationNotice)
components/ui/   shadcn primitives (button, input, label, card, dialog) — generated, keep edits minimal
types/        one .ts per resource, plain interfaces matching backend JSON shapes
lib/          money.ts (cents <-> euros helpers), utils.ts (cn() for class merging)
```

There is no global state library (Redux/Zustand) — React Query's cache is the app's server-state
store, and local UI state lives in `useState` inside each route component.

## Routing (`App.tsx`)

Flat `<Routes>` tree in `App.tsx`, no route config file. Three role-scoped subtrees plus a few
public pages:

- `/client`, `/client/partners`, `/client/partners/:id` — wrapped in `<RoleGuard role="client">`
- `/partner`, `/partner/dashboard` — wrapped in `<RoleGuard role="partner">`
- `/admin/*` — wrapped in `<RoleGuard role="admin">`, nested under `<AdminLayout>` (its own
  `<Outlet>`-based sub-routing for `partners` / `clients` / `clients/:id`)
- `/partners`, `/partners/:id` — public read-only mirror of the client partner catalog
- `/login`, `/terms`, `/accessibility` — public
- `/` — `RootRedirect`, sends the user to their role's home or to `/login`

**There is currently no registration route or page for any role** — `/partners/register`,
`/clients/register`, `/admin/register` exist on the backend (`backend/routing/router.go`) but
nothing in `frontend/src/api` or `frontend/src/routes` calls them yet.

## Auth flow

`frontend/src/auth/useAuth.ts` wraps three React Query hooks around `frontend/src/api/auth.ts`:

- `useAuth()` — runs `GET /auth/me` under query key `['me']`. Returns
  `{ status: 'loading' | 'unauthenticated' | 'authenticated', user }`. `retry: false` so a 401
  resolves immediately to `unauthenticated` instead of retrying.
- `useLogin()` — `POST /auth/login`, on success seeds `['me']` with the returned user via
  `queryClient.setQueryData`, avoiding a refetch.
- `useLogout()` — `POST /auth/logout`, on success sets `['me']` to `null`.

`RoleGuard` (`routes/RoleGuard.tsx`) reads `useAuth()`: renders nothing while loading, redirects to
`/login` if unauthenticated, redirects to `/` if the logged-in user's role doesn't match the route,
otherwise renders `children`. `RootRedirect` does the loading/unauthenticated checks the same way,
then sends authenticated users to `roleHomePaths[user.role]`.

There is no registration-side equivalent yet (no `useRegisterPartner`-style hook, no unauthenticated
"pending approval" state) — see `docs/frontend/partner-registration.md`.

## Data fetching pattern

Every backend resource follows the same three-layer pattern; use it as the template for anything
new:

1. **`types/<resource>.ts`** — plain interfaces mirroring the backend JSON response exactly
   (including Go-style `PascalCase` field names where the backend hasn't wrapped them — e.g.
   `Partner.BusinessName`, `Partner.MinisterPick`).
2. **`api/<resource>.ts`** — one function per endpoint, calling `apiFetch<ResponseType>(path, opts)`
   from `api/client.ts`. Functions take plain params, build query strings by hand
   (`URLSearchParams`), and return the typed promise directly — no try/catch here.
3. **`<domain>/use<Resource>Data.ts`** — React Query hooks (`useQuery`/`useMutation`) wrapping the
   api functions. Query keys are arrays of the function name + params, e.g.
   `['partners', search, category, page, limit]`. Mutations that should update cached reads do so
   explicitly in `onSuccess` (see `useLogin`/`useLogout`), otherwise the caller invalidates/refetches.

`api/client.ts`'s `apiFetch`:
- Prefixes every call with `VITE_API_BASE_URL` (from `.env`).
- Always sends `credentials: 'include'` and `Content-Type: application/json`.
- JSON-encodes a `body` option if present.
- On a non-2xx response, tries to parse `{ message }` from the JSON body and throws `ApiError`
  (extends `Error`, carries `.status`); falls back to `response.statusText`.
- Returns `undefined` for `204 No Content`, otherwise parses the body as JSON.

Route components branch on query/mutation state directly (`isLoading`, `isError`, `isSuccess`,
`isPending`) rather than using a shared loading/error wrapper component. Error messages from
`ApiError` are often mapped to friendlier copy by status code in the route itself (see
`getCollectErrorMessage` in `routes/partner/PartnerHome.tsx`).

## Forms

There's no form library (no react-hook-form/formik/zod) — forms are plain controlled inputs with
one `useState` per field, a `handleSubmit` that calls `event.preventDefault()` then
`mutation.mutate(...)`, and inline validation done by hand before calling mutate (see
`handleConfirmAmount` in `PartnerHome.tsx` for a validate-then-mutate example, and `LoginPage.tsx`
for the simplest two-field form). Errors are rendered as a `<p>` reading `mutation.error.message`
below the fields, gated on `mutation.isError`.

## Styling / design system

Tailwind v4 config lives entirely in `frontend/src/index.css` (no `tailwind.config.js`) via
`@theme` blocks: custom color tokens (`--color-brand-blue`, `--color-ink-900`,
`--color-surface-*`, `--color-danger-600`, etc.), custom text styles bundling
size/line-height/letter-spacing/weight (`text-h1`, `text-body`, `text-label-caps`,
`text-button`, ...), custom radii (`rounded-card`, `rounded-panel`, `rounded-row`), and custom
shadows (`shadow-card-highlight`, `shadow-row-raised`, `shadow-key-secondary`). A second `:root`/
`.dark` block maps those onto shadcn's semantic tokens (`--primary`, `--border`, etc.) for the
`components/ui/*` primitives.

Conventions seen across every route:
- Page shell: `<div className="min-h-screen bg-surface-200 p-4">` with an inner
  `<div className="mx-auto flex max-w-sm flex-col gap-6">` — everything is designed mobile-first
  at a fixed narrow max width, no responsive breakpoints in use yet.
- Two font families: `font-display` (Marianne, for headings/labels/buttons, usually uppercase +
  letter-spaced) and `font-sans` (Spectral, a serif, for body copy) — note the names are swapped
  from what you'd expect.
- Headings use `text-h1`/`text-h2` + `font-display` + `text-brand-blue`; body text uses
  `text-body`/`text-caption` + `font-sans` + `text-ink-900`/`text-ink-600`; buttons/labels use
  `text-button`/`text-label-caps` + `font-display` + `uppercase`.
- List rows: `rounded-row bg-surface-050 px-3 py-3 shadow-row-raised`.
- Primary actions: `<Button>` (default variant) full-width, `h-11`, `rounded-panel`; secondary
  actions restyle the same `<Button>` with `bg-surface-300 ... shadow-key-secondary`.
  `components/ui/button.tsx` also has `outline`/`ghost`/`destructive`/`link` variants and
  `sm`/`lg`/`icon` sizes from its `cva` config if a plain secondary restyle isn't the right fit.
  `components/ui/card.tsx`, `input.tsx`, `label.tsx`, `dialog.tsx` are the other available shadcn
  primitives.
- `@/...` is aliased to `frontend/src` (see `vite.config.ts`); use it instead of relative
  `../../` paths in new files, matching most (not all — some older files use relative imports)
  existing code.
- `components/Wordmark.tsx` (logo + "CartePro" text) appears at the top of every full-page screen.
- `components/SimulationNotice.tsx` is a small disclaimer banner used on payment-related screens
  (this is a training/demo app, not real money) — has a `tone="dark"` variant for use over the
  dark card gradient.

## Types quick reference

- `types/auth.ts` — `Role = 'client' | 'partner' | 'admin'`, `User { id, email, role }`.
- `types/partner.ts` — `Partner { ID, BusinessName, Category, Region, Address, MinisterPick }`
  plus transaction/dashboard/QR-validation response shapes.
- `types/admin.ts` / `types/client.ts` — admin- and client-scoped resource shapes, including
  `PartnerStatus = 'pending' | 'approved' | 'rejected'` and `AdminPartner` (adds `Status`,
  `Balance`, `Siret` on top of the base `Partner` fields).

Note the inconsistency: `types/auth.ts` uses `camelCase` (a hand-shaped API), while
`types/partner.ts`/`types/admin.ts` use `PascalCase` (mirroring Go struct field JSON tags
verbatim, no `json:"..."` renaming on those endpoints). Match whichever casing the actual backend
response uses for any new resource — check the handler in `backend/users/*.go` rather than
assuming.
