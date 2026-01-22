# ---------- Build stage ----------
FROM golang:1.24 AS builder

WORKDIR /app

# Enable Go toolchain auto-download (required for toolchain directive)
ENV GOTOOLCHAIN=auto

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o server ./cmd/server


# ---------- Runtime stage ----------
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates \
      tzdata \
      ghostscript \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Create non-root user
RUN useradd -r -u 10001 nonroot

# Copy binary
COPY --from=builder --chown=nonroot:nonroot /app/server /app/server

USER nonroot

EXPOSE 3001

ENTRYPOINT ["/app/server"]
