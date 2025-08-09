// middleware/env.go
package middleware

import (
	"quiver/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// EnvMiddleware 用于处理和验证 :env 参数
func EnvMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		env := c.Params("env")
		if env == "" {
			return utils.BadRequest(c, "env is required")
		}

		env = strings.ToLower(env)
		if !utils.ValidateEnv(env) {
			return utils.NotFound(c, "invalid env")
		}
		c.Locals("env", env)
		return c.Next()
	}
}
