# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o backend .

# Runtime stage
FROM alpine:3.22

WORKDIR /app

# Install CA certificates for HTTPS requests
# RUN apk --no-cache add ca-certificates

COPY --from=builder /app/backend .

EXPOSE 8080

CMD ["./backend"]
