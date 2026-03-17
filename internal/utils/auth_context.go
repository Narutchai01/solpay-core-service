package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetUserIDFromLocals parses user/account ID from Fiber locals.
func GetUserIDFromLocals(c *fiber.Ctx) (int64, bool) {
	userID := c.Locals("userID")
	if userID == nil {
		userID = c.Locals("accountID")
	}

	switch value := userID.(type) {
	case int64:
		return value, value > 0
	case uint:
		return int64(value), value > 0
	case int:
		return int64(value), value > 0
	case float64:
		return int64(value), value > 0
	case string:
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed <= 0 {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}
