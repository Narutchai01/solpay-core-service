# --- Stage 1: Builder ---
# ใช้ Go version 1.25.5 ตามที่ระบุ
FROM golang:1.25.5-alpine AS builder

# ติดตั้ง git เผื่อจำเป็นสำหรับบาง module
# RUN apk add --no-cache git

WORKDIR /app

# 1. Copy ไฟล์ go.mod และ go.sum ก่อน เพื่อทำ cache layer
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy Source Code ทั้งหมด (รวมถึง internal, cmd, config)
COPY . .

# 3. Build โดยระบุ path ไปที่ cmd/main.go
# ตั้งชื่อไฟล์ binary ว่า 'server' (หรือชื่ออื่นตามต้องการ)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server cmd/main.go

# --- Stage 2: Runner ---
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Bangkok

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]