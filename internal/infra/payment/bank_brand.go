package payment

// BankBrand represents a Thai bank identifier for Omise.
type BankBrand string

// Supported bank brands.
const (
	BankBBL   BankBrand = "bbl"
	BankKBank BankBrand = "kbank"
	BankKTB   BankBrand = "ktb"
	BankSCB   BankBrand = "scb"
	BankBAY   BankBrand = "bay"
	BankTTB   BankBrand = "ttb"
	BankGSB   BankBrand = "gsb"
)

// IsValid reports whether b is a recognised bank brand.
func (b BankBrand) IsValid() bool {
	switch b {
	case BankBBL, BankKBank, BankKTB, BankSCB, BankBAY, BankTTB, BankGSB:
		return true
	}
	return false
}
