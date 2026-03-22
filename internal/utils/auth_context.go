package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetUserIDFromLocals parses user/account ID from Fiber locals.
func GetUserIDFromLocals(c *fiber.Ctx) (uint, bool) {
	userID := c.Locals("userID")
	if userID == nil {
		userID = c.Locals("accountID")
	}

	switch value := userID.(type) {
	case int64:
		return uint(value), value > 0
	case uint:
		return value, value > 0
	case int:
		return uint(value), value > 0
	case float64:
		return uint(value), value > 0
	case string:
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed <= 0 {
			return 0, false
		}
		return uint(parsed), true
	default:
		return 0, false
	}
}
