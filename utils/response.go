package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.JSON(Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest 400错误
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, 400, message, nil)
}

// NotFound 404错误
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, 404, message, nil)
}

// Conflict 409错误
func Conflict(c *fiber.Ctx, message string) error {
	return Error(c, 409, message, nil)
}

// InternalError 500错误
func InternalError(c *fiber.Ctx, message string) error {
	return Error(c, 500, message, nil)
}

// Timeout 408超时
func Timeout(c *fiber.Ctx) error {
	return Error(c, 408, "timeout", nil)
}
