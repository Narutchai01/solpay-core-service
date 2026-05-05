package utils

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // Support JPEG decoding
	"image/png"
	"net/http"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

// downloadImage fetches an image from a given URL and decodes it.
func downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code when fetching image: %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image from URL: %w", err)
	}

	return img, nil
}

type SlipOffchain struct {
	Address       string
	Amount        float64
	TransactionID string
	PromptPayID   string
	CreatedAt     string
}

func GetSlipOFFCHAINInformation(data SlipOffchain) ([]byte, error) {
	cfg := config.LoadConfig()
	// 1. โหลดรูป Template พื้นหลัง (mock data from URL)
	templateURL := "https://cbsievyfzgizotkepdle.supabase.co/storage/v1/object/public/masterdata/template_offchain.png"
	im, err := downloadImage(templateURL)
	if err != nil {
		// Fallback to local file if download fails
		im, err = gg.LoadImage("internal/asstes/template_offchain.png")
		if err != nil {
			return nil, fmt.Errorf("cannot load template from URL or local path: %w", err)
		}
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

type SlipOnchain struct {
	Address       string
	THBAmount     float64
	USDTAmount    float64
	FreeAmount    float64
	TransactionID string
	PromptPayID   string
	CreatedAt     string
}

func GetSlipOnChain(data SlipOnchain) ([]byte, error) {

	cfg := config.LoadConfig()
	// 1. โหลดรูป Template พื้นหลัง
	templateURL := "https://cbsievyfzgizotkepdle.supabase.co/storage/v1/object/public/masterdata/template_onchain.png" // Mock data for onchain template
	im, err := downloadImage(templateURL)
	if err != nil {
		// Fallback to offchain URL if onchain fails
		fallbackURL := "https://cbsievyfzgizotkepdle.supabase.co/storage/v1/object/public/masterdata/template_offchain.png"
		im, err = downloadImage(fallbackURL)
		if err != nil {
			// Fallback to local files
			im, err = gg.LoadImage("internal/asstes/template_onChain.png")
			if err != nil {
				im, err = gg.LoadImage("internal/asstes/template_offchain.png")
				if err != nil {
					return nil, fmt.Errorf("cannot load template from URL or local paths: %w", err)
				}
			}
		}
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
	dc.DrawString(data.TransactionID, 21, 64)
	dc.DrawString(data.CreatedAt, 24, 88)

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
	dc.DrawStringAnchored(fmt.Sprintf("%.2f THB", data.THBAmount), float64(dc.Width()-11), 372, 1.0, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%.6f USDT", data.USDTAmount), float64(dc.Width()-11), 409, 1.0, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%.2f USDT", data.FreeAmount), float64(dc.Width()-11), 445, 1.0, 0.5)

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

	// --- ทำป้าย Category (ปุ่มสีม่วงๆ ด้านล่าง) ---
	// วาดกล่องสี่เหลี่ยมขอบมน
	dc.SetHexColor("#9B8BFF") // สีพื้นหลังปุ่ม
	dc.DrawRoundedRectangle(float64(dc.Width()-200), 800, 150, 40, 10)
	dc.Fill()

	// 6. เข้ารหัสรูปภาพเป็น byte slice
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, fmt.Errorf("cannot encode output image: %w", err)
	}

	return buf.Bytes(), nil
}
