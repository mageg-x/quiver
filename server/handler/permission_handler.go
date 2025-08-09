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

type PermissionHandler struct {
	permService *services.PermissionService
}

func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{
		permService: services.NewPermissionService(),
	}
}

func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	var perm models.Permission
	if err := c.BodyParser(&perm); err != nil {
		logger.GetLogger("quiver").Errorf("invalid request body: %v", err)
		return utils.BadRequest(c, "invalid request body")
	}
	perm.UserID = userID
	logger.GetLogger("quiver").Infof("create permission: %v", perm)

	if err := h.permService.CreatePermission(env, &perm); err != nil {
		logger.GetLogger("quiver").Errorf("create permission failed: %v", err)
		return utils.InternalError(c, err.Error())
	}

	// 存入 cache
	if data, err := msgpack.Marshal(perm); err == nil && len(data) > 0 {
		_ = cache.Set(perm.CacheKey(env), data, 300*time.Second)
		logger.GetLogger("quiver").Infof("set permission %d to cache success", perm.ID)
	}

	return utils.Success(c, 0, "success", perm)
}

func (h *PermissionHandler) ListPermission(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	perms, err := h.permService.ListPermission(env, userID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, 0, "success", perms)
}

func (h *PermissionHandler) GetPermission(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	permissionId, err := c.ParamsInt("permission_id")
	if err != nil {
		return utils.BadRequest(c, "invalid permission_id")
	}

	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	// 先从cache中读取
	_perm := &models.Permission{UserID: userID, ID: uint64(permissionId)}
	data, ok, err := cache.Get(_perm.CacheKey(env))
	if ok && len(data) > 0 {
		perm := &models.Permission{}
		if err := msgpack.Unmarshal(data, perm); err != nil {
			// 数据已经被破坏
			_ = cache.Delete(_perm.CacheKey(env))
			logger.GetLogger("quiver").Warnf("error unmarshaling permission: %d, %v", permissionId, err)
		}
		logger.GetLogger("quiver").Infof("get permission %d from cache success", permissionId)
		return utils.Success(c, 0, "success", perm)
	}

	perm, err := h.permService.GetPermission(env, userID, uint64(permissionId))
	if err != nil {
		return utils.NotFound(c, "permission not found")
	}

	// 存入 cache
	if data, err := msgpack.Marshal(perm); err == nil && len(data) > 0 {
		_ = cache.Set(perm.CacheKey(env), data, 300*time.Second)
		logger.GetLogger("quiver").Infof("set permission %d to cache success", perm.ID)
	}

	return utils.Success(c, 0, "success", perm)
}

func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	permissionId, err := c.ParamsInt("permission_id")
	if err != nil || permissionId <= 0 {
		return utils.BadRequest(c, "invalid permission_id")
	}

	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	var update models.Permission
	if err := c.BodyParser(&update); err != nil {
		logger.GetLogger("quiver").Errorf("invalid request body: %v", err)
		return utils.BadRequest(c, "invalid request body")
	}

	if update.ResourceID == 0 || update.ResourceType == "" || update.Action == "" {
		logger.GetLogger("quiver").Errorf("required params not exist")
		return utils.BadRequest(c, "required params not exist")
	}

	update.UserID = userID
	update.ID = uint64(permissionId)
	logger.GetLogger("quiver").Infof("update permission: %+v", update)
	perm, err := h.permService.UpdatePermission(env, &update)
	if err != nil {
		if err.Error() == "permission not found" {
			return utils.NotFound(c, "permission not found")
		}
		return utils.InternalError(c, err.Error())
	}

	// 存入 cache
	if data, err := msgpack.Marshal(perm); err == nil && len(data) > 0 {
		_ = cache.Set(perm.CacheKey(env), data, 300*time.Second)
		logger.GetLogger("quiver").Infof("set permission %d to cache success", perm.ID)
	}

	return utils.Success(c, 0, "success", perm)
}

func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	permissionId, err := c.ParamsInt("permission_id")
	if err != nil {
		return utils.BadRequest(c, "invalid permission_id")
	}

	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	// 从缓存中删除
	_perm := models.Permission{UserID: userID, ID: uint64(permissionId)}
	if err := cache.Delete(_perm.CacheKey(env)); err != nil {
		logger.GetLogger("quiver").Errorf("delete permission from cache failed: %v", err)
	}

	if err := h.permService.DeletePermission(env, userID, uint64(permissionId)); err != nil {
		if err.Error() == "permission not found" {
			return utils.NotFound(c, "permission not found")
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, 0, "success", fiber.Map{"permission_id": permissionId})
}
