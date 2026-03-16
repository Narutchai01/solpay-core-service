package utils

import "errors"

func VerifyWithSlippage(quoteRate float64, currentMarketRate float64, maxSlippage float64) error {
	// คำนวณส่วนต่างราคา
	diff := ((currentMarketRate - quoteRate) / quoteRate) * 100

	// ถ้าเรทเปลี่ยนไป (ในทางที่ User ขาดทุน) เกินกว่าที่ตั้งไว้
	if diff > maxSlippage {
		return errors.New("Slippage exceeded: Market price changed too much")
	}
	return nil
}
