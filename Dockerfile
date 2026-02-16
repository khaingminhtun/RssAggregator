# syntax=docker/dockerfile:1.6

# -------- Build Stage --------
FROM golang:1.25-alpine AS build

WORKDIR /app

# Install CA certificates (needed for HTTPS)
RUN apk add --no-cache ca-certificates git

# Create non-root user
RUN adduser -D -u 1001 nonroot

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./

# Download dependencies (BuildKit cache)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy all source code
COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -extldflags '-static'" \
    -tags netgo \
    -o api-golang ./cmd
    # <-- points to cmd/ folder where main.go exists

# -------- Runtime Stage --------
FROM scratch

# Copy CA certificates (for HTTPS)
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy non-root user
COPY --from=build /etc/passwd /etc/passwd

# Copy binary
COPY --from=build /app/api-golang /api-golang

# Copy .env file so app can read environment variables
COPY --from=build /app/.env .env

# Run as non-root
USER nonroot

EXPOSE 3000

CMD ["/api-golang"]
