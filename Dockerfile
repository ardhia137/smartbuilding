# Gunakan image base untuk Go
FROM golang:1.21-alpine AS builder

# Set work directory
WORKDIR /app

# Copy semua file ke dalam container
COPY . .

# Unduh dependencies
RUN go mod tidy

# Build aplikasi
RUN go build -o smartbuilding

# Gunakan image lebih kecil untuk menjalankan aplikasi
FROM alpine:latest

# Set work directory
WORKDIR /root/

# Copy file yang sudah dikompilasi dari tahap sebelumnya
COPY --from=builder /app/smartbuilding .

# Ekspos port yang digunakan aplikasi
EXPOSE 1312

# Jalankan aplikasi
CMD ["./smartbuilding"]
