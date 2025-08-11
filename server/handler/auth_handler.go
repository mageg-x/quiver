package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vmihailenco/msgpack/v5"
	"quiver/cache"
	"quiver/logger"
	"quiver/models"
	"quiver/services"
	"quiver/utils"
	"time"
)

// LoginRequest 登录请求体
type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// LoginResponse 响应结构
type LoginResponse struct {
	UserName  string    `json:"user_name"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言

	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user := services.NewUserService()
	u, err := user.GetUserByName(env, req.UserName)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.NotFound(ctx, "user not found")
		}
		return utils.InternalError(ctx, err.Error())
	}

	if !u.CheckPassword(req.Password) {
		return utils.Unauthorized(ctx, "invalid password")
	}

	// 生成临时access token
	s := services.NewAccessKeyService()
	ak, sk, err := s.GenerateKeys()

	accessKey := &models.AccessKey{
		UserID:     u.UserID,
		AccessKey:  ak,
		SecretKey:  sk,
		ExpireAt:   utils.Ptr(time.Now().Add(time.Hour * 24)),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	// 把 accessKey 存入 cache
	if data, err := msgpack.Marshal(&accessKey); err == nil && len(data) > 0 {
		_ = cache.Set(accessKey.CacheKey(env), data, time.Hour*24)
		logger.GetLogger("quiver").Infof("write accesskey %s to cache", accessKey.AccessKey)
	}

	return utils.Success(ctx, 0, "success", LoginResponse{
		ExpiresAt: *accessKey.ExpireAt,
		Token:     accessKey.AccessKey,
		UserName:  u.UserName,
	})
}

func (h *AuthHandler) RefreshToken(ctx *fiber.Ctx) error {
	return utils.Success(ctx, 0, "success", nil)
}

func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	return utils.Success(ctx, 0, "success", nil)
}
