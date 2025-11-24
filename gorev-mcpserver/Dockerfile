FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o gorev ./cmd/gorev

# Final stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates && adduser -D -g '' gorev

# Create directories
RUN mkdir -p /workspace /data && chown -R gorev:gorev /workspace /data

# Copy binary
COPY --from=builder /app/gorev /usr/local/bin/gorev
USER gorev
WORKDIR /workspace
EXPOSE 5082
VOLUME ["/data", "/workspace"]
CMD ["gorev", "daemon", "--detach"]
