package payment

// ประกาศ Custom Type เป็น BankBrand
type BankBrand string

// ประกาศ Enum ค่าคงที่สำหรับธนาคารต่างๆ
const (
	BankBBL   BankBrand = "bbl"
	BankKBank BankBrand = "kbank"
	BankKTB   BankBrand = "ktb"
	BankSCB   BankBrand = "scb"
	BankBAY   BankBrand = "bay"
	BankTTB   BankBrand = "ttb"
	BankGSB   BankBrand = "gsb"
)

// ฟังก์ชันเสริม: เอาไว้ Validate ว่า Frontend ส่งค่าที่รองรับมาให้เราหรือไม่
func (b BankBrand) IsValid() bool {
	switch b {
	case BankBBL, BankKBank, BankKTB, BankSCB, BankBAY, BankTTB, BankGSB:
		return true
	}
	return false
}
