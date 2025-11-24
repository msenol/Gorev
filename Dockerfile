FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache make git

# Copy go mod files from gorev-mcpserver directory
COPY gorev-mcpserver/go.mod gorev-mcpserver/go.sum ./
RUN go mod download

# Copy source code
COPY gorev-mcpserver/ .

# Build the binary (Web UI is pre-built and embedded)
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o gorev ./cmd/gorev

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk add --no-cache ca-certificates && adduser -D -g '' gorev

# Create directories for data
RUN mkdir -p /workspace /data && chown -R gorev:gorev /workspace /data

# Copy the binary from builder
COPY --from=builder /app/gorev /usr/local/bin/gorev

# Switch to non-root user
USER gorev

# Set working directory
WORKDIR /workspace

# Default port
EXPOSE 5082

# Volume for persistent data
VOLUME ["/data", "/workspace"]

# Run server in foreground (daemon mode not suitable for containers)
CMD ["gorev", "serve", "--api-port", "5082"]
