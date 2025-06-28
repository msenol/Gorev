# Build aşaması
FROM golang:1.23-alpine AS builder

# Build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev git

WORKDIR /app

# Go modüllerini önce kopyala (cache için)
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# Binary'yi derle
RUN CGO_ENABLED=1 GOOS=linux go build -mod=mod -buildvcs=false -a -installsuffix cgo -ldflags '-extldflags "-static"' -o gorev ./cmd/gorev

# Final aşama
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Binary'yi kopyala
COPY --from=builder /app/gorev .

# Migration dosyalarını kopyala
COPY --from=builder /app/internal/veri/migrations ./internal/veri/migrations

# Veri dizini oluştur
RUN mkdir -p /app/data

# Çalışma portu (MCP stdio kullandığı için port gerekmez)
# EXPOSE değeri yok

# Varsayılan komut
ENTRYPOINT ["./gorev"]
CMD ["serve"]