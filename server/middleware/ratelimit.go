package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10000,           // 最大请求数
		Expiration: 1 * time.Minute, // 时间窗口
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // 基于IP限流
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"code":    429,
				"message": "Too many requests",
				"data":    nil,
			})
		},
	})
}
