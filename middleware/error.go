package middleware

import (
	"log"
	"quiver/utils"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler(c *fiber.Ctx, err error) error {
	// 默认状态码500
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// 如果是Fiber错误，获取状态码和消息
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// 记录错误日志
	log.Printf("Error [%d]: %s - Path: %s", code, message, c.Path())

	// 返回错误响应
	return utils.Error(c, code, message, nil)
}

// NotFoundHandler 404处理器
func NotFoundHandler(c *fiber.Ctx) error {
	return utils.NotFound(c, "Route not found")
}
