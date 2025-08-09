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
	// 写入 cache
	if data, err := msgpack.Marshal(&ak); err == nil && len(data) > 0 {
		_ = cache.Set(ak.CacheKey(env), data, 300*time.Second)
		logger.GetLogger("quiver").Infof("write accesskey %s to cache", ak.AccessKey)
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

	// 先从cache中读取
	_ak := &models.AccessKey{AccessKey: accessKey, UserID: userID}
	data, ok, err := cache.Get(_ak.CacheKey(env))
	//logger.GetLogger("quiver").Infof(" %v, %v, %v", ok, err, data)
	if ok && len(data) > 0 {
		ak := &models.AccessKey{}
		if err := msgpack.Unmarshal(data, ak); err != nil {
			// 数据已经被破坏
			_ = cache.Delete(_ak.CacheKey(env))
			logger.GetLogger("quiver").Warnf("error unmarshaling access key: %s, %v", accessKey, err)
		}
		logger.GetLogger("quiver").Infof("get access key  %s from cache success", accessKey)
		return utils.Success(c, 0, "success", ak)
	}

	ak, err := h.akService.GetAccessKey(env, userID, accessKey)
	if err != nil {
		logger.GetLogger("quiver").Errorf("get accesskey failed: %v", err)
		return utils.InternalError(c, "failed to get accesskey")
	}

	// 写入 cache
	if data, err := msgpack.Marshal(&ak); err == nil && len(data) > 0 {
		_ = cache.Set(ak.CacheKey(env), data, 300*time.Second)
		logger.GetLogger("quiver").Infof("set access key %s to cache success", ak.AccessKey)
	}

	// 不返回 secret_key
	//ak.SecretKey = ""
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

	// 从cache中删除 accesskey
	_ak := &models.AccessKey{AccessKey: accessKey, UserID: userID}
	if err = cache.Delete(_ak.CacheKey(env)); err != nil {
		logger.GetLogger("quiver").Errorf("error deleting accesskey from cache: %v", err)
	} else {
		logger.GetLogger("quiver").Infof("accesskey deleted from cache: %s", accessKey)
	}

	if err := h.akService.DeleteAccessKey(env, userID, accessKey); err != nil {
		if err.Error() == "accesskey not found" {
			return utils.NotFound(c, "accesskey not found")
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, 0, "success", fiber.Map{"accesskey": accessKey})
}
