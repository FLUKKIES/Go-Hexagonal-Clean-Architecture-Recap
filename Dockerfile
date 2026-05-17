FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# ก็อปปี้ไฟล์ go.mod และ go.sum มาโหลด Dependencies ก่อน
# มันจะอ่านไฟล์ go.mod และ go.sum แล้วดาวน์โหลด Dependencies ทุกตัวตามเวอร์ชันที่ระบุไว้เป๊ะๆ ลงเครื่อง
COPY go.mod go.sum ./

# ดาวน์โหลด dependencies
RUN go mod download 

# ก็อปปี้ไฟล์ source code ทั้งหมดมาลงใน image
COPY . .

# compile source code
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api/main.go

# ==========================================
# Stage 2: Production (เอาแค่ไฟล์ที่ Build เสร็จแล้วมาใช้ Image จะเล็กมาก)
# ==========================================
# Alpine เป็น Image ที่มีขนาดเล็กมาก (แค่ประมาณ 5MB) ทำให้ Container ของเรามีขนาดเล็กตามไปด้วย
FROM alpine:latest AS runner

WORKDIR /app

# เอาไฟล์ที่ Build เสร็จแล้ว (main) จาก Stage แรก มาใส่ในนี้
COPY --from=builder /app/main .

# ระบุ Port ที่จะเปิด
EXPOSE 4001

# คำสั่งที่จะให้ container รันเมื่อถูกสร้างขึ้นมา
CMD [ "./main" ]