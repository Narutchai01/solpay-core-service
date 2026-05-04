package utils

import (
	"bytes"
	"fmt"
	"image/png"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

type SlipOffchain struct {
	Address       string
	Amount        float64
	TransactionID string
	PromptPayID   string
	CreatedAt     string
}

func GetSlipOFFCHAINInformation(data SlipOffchain) ([]byte, error) {
	cfg := config.LoadConfig()
	// 1. โหลดรูป Template พื้นหลัง
	im, err := gg.LoadImage("internal/asstes/template_offchain.png")
	if err != nil {
		return nil, fmt.Errorf("cannot load template: %w", err)
	}

	// สร้าง Context จากรูป Template
	dc := gg.NewContextForImage(im)

	// 2. ตั้งค่าสีฟอนต์ (สีเทาเข้ม/ดำ)
	dc.SetHexColor("#4A4A4A")

	// โหลดฟอนต์ (ต้องมีไฟล์ฟอนต์ .ttf ในโฟลเดอร์)
	// ปรับขนาดฟอนต์ (เช่น 20) ตามความเหมาะสมของขนาดรูป
	if err := dc.LoadFontFace("internal/asstes/Kanit-Regular.ttf", 14); err != nil {
		return nil, fmt.Errorf("cannot load font: %w", err)
	}

	// 3. เริ่มเขียนข้อความลงบนรูป (ต้องกะระยะแกน X, Y ให้ตรงกับรูปของคุณ)
	// dc.DrawString("ข้อความ", แกนX, แกนY)
	dc.DrawString(data.TransactionID, 21, 65)
	dc.DrawString(data.CreatedAt, 21, 90)

	displayAddress := data.Address
	if len(displayAddress) > 12 {
		displayAddress = displayAddress[:6] + "..." + displayAddress[len(displayAddress)-6:]
	}
	dc.DrawString(displayAddress, 85, 144) // เลขบัญชีผู้โอน

	// โหลดฟอนต์แบบหนา (Bold) สำหรับชื่อ
	dc.LoadFontFace("internal/asstes/Kanit-Regular.ttf", 14)
	dc.DrawString(data.PromptPayID, 85, 280)

	// จำนวนเงิน (ชิดขวา)
	dc.SetHexColor("#000000")
	dc.LoadFontFace("internal/asstes/Kanit-Regular.ttf", 14)
	// ตัวเลข 1.0 คือให้ anchor (จุดศูนย์กลางการจัดหน้า) อยู่ด้านขวา
	dc.DrawStringAnchored(fmt.Sprintf("%.2f THB", data.Amount), float64(dc.Width()-40), 370, 1.0, 0.5)

	// 4. สร้าง QR Code
	qrData := fmt.Sprintf(cfg.SLIP_INFORMATION, data.TransactionID) // ข้อมูล PromptPay หรือ URL
	qrBytes, err := qrcode.Encode(qrData, qrcode.Medium, 72)        // 180 คือขนาดพิกเซล กว้างxยาว
	if err != nil {
		return nil, fmt.Errorf("cannot generate QR: %w", err)
	}

	// แปลง bytes เป็น image.Image
	qrImg, err := png.Decode(bytes.NewReader(qrBytes))
	if err != nil {
		return nil, fmt.Errorf("cannot decode QR image: %w", err)
	}

	// 5. แปะรูป QR Code ลงบน Template
	// กำหนดจุด X, Y ที่จะวาง QR Code (เช่นวางไว้ขวากลางๆ)
	qrX := dc.Width() - 97
	qrY := 248
	dc.DrawImage(qrImg, qrX, qrY)

	// 6. เข้ารหัสรูปภาพเป็น byte slice
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, fmt.Errorf("cannot encode output image: %w", err)
	}

	return buf.Bytes(), nil
}
