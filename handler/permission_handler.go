package handler

import (
	"github.com/gofiber/fiber/v2"
	"quiver/logger"
	"quiver/models"
	"quiver/services"
	"quiver/utils"
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

	if err := h.permService.CreatePermission(env, &perm); err != nil {
		logger.GetLogger("quiver").Errorf("create permission failed: %v", err)
		return utils.InternalError(c, err.Error())
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
	id, err := c.ParamsInt("permission_id")
	if err != nil {
		return utils.BadRequest(c, "invalid permission_id")
	}

	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	perm, err := h.permService.GetPermission(env, userID, uint64(id))
	if err != nil {
		return utils.NotFound(c, "permission not found")
	}

	return utils.Success(c, 0, "success", perm)
}

func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	id, err := c.ParamsInt("permission_id")
	if err != nil {
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
	update.UserID = userID
	update.ID = uint64(id)

	perm, err := h.permService.UpdatePermission(env, &update)
	if err != nil {
		if err.Error() == "permission not found" {
			return utils.NotFound(c, "permission not found")
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, 0, "success", perm)
}

func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	id, err := c.ParamsInt("permission_id")
	if err != nil {
		return utils.BadRequest(c, "invalid permission_id")
	}

	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	if err := h.permService.DeletePermission(env, userID, uint64(id)); err != nil {
		if err.Error() == "permission not found" {
			return utils.NotFound(c, "permission not found")
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, 0, "success", fiber.Map{"permission_id": id})
}
