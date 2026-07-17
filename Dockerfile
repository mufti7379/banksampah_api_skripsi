# ── STAGE 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Download dependency dulu (memanfaatkan Docker cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary — CGO_ENABLED=0 agar bisa jalan di Alpine tanpa glibc
RUN CGO_ENABLED=0 GOOS=linux go build -o banksampah-api ./cmd/api/main.go

# ── STAGE 2: Runtime (image sekecil mungkin) ─────────────────────────────────
FROM alpine:3.19

WORKDIR /app

# Install timezone data (penting untuk DATETIME Asia/Jakarta)
RUN apk --no-cache add tzdata ca-certificates

COPY --from=builder /app/banksampah-api .

# Cloud Run mendengarkan di port 8080 secara default
EXPOSE 8080

CMD ["./banksampah-api"]
