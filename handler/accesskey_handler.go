package handler

import (
	"github.com/gofiber/fiber/v2"
	"quiver/logger"
	"quiver/services"
	"quiver/utils"
)

type AccessKeyHandler struct {
	akService *services.AccessKeyService
}

func NewAccessKeyHandler() *AccessKeyHandler {
	return &AccessKeyHandler{
		akService: services.NewAccessKeyService(),
	}
}

func (h *AccessKeyHandler) CreateAccessKey(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	ak, err := h.akService.CreateAccessKey(env, userID)
	if err != nil {
		logger.GetLogger("quiver").Errorf("create accesskey failed: %v", err)
		return utils.InternalError(c, err.Error())
	}

	// 注意：secret_key 只在此刻返回
	return utils.Success(c, 0, "success", ak)
}

func (h *AccessKeyHandler) ListAccessKey(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	aks, err := h.akService.ListAccessKey(env, userID)
	if err != nil {
		logger.GetLogger("quiver").Errorf("list access keys failed: %v", err)
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, 0, "success", aks)
}

func (h *AccessKeyHandler) GetAccessKey(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	accessKey := c.Params("accesskey")
	if err != nil {
		logger.GetLogger("quiver").Errorf("invalid accesskey: %v", err)
		return utils.BadRequest(c, "invalid accesskey")
	}

	ak, err := h.akService.GetAccessKey(env, userID, accessKey)
	if err != nil {
		logger.GetLogger("quiver").Errorf("get accesskey failed: %v", err)
		return utils.InternalError(c, "failed to get accesskey")
	}

	// 不返回 secret_key
	ak.SecretKey = ""
	return utils.Success(c, 0, "success", ak)
}

func (h *AccessKeyHandler) DeleteAccessKey(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	accessKey := c.Params("accesskey")
	if err != nil {
		logger.GetLogger("quiver").Errorf("invalid accesskey: %v", err)
		return utils.BadRequest(c, "invalid accesskey")
	}
	if err := h.akService.DeleteAccessKey(env, userID, accessKey); err != nil {
		if err.Error() == "accesskey not found" {
			return utils.NotFound(c, "accesskey not found")
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, 0, "success", fiber.Map{"accesskey": accessKey})
}
