package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"quiver/logger"
	"quiver/models"
	"quiver/services"
	"quiver/utils"
	"strconv"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	env := c.Locals("env").(string)

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		logger.GetLogger("quiver").Error("invalid request body")
		return utils.BadRequest(c, "invalid request body")
	}

	if err := h.userService.CreateUser(env, &user); err != nil {
		logger.GetLogger("quiver").Errorf("create user failed: %v", err)
		if err.Error() == "username already exists" || err.Error() == "email already exists" || err.Error() == "phone already exists" {
			logger.GetLogger("quiver").Errorf("username %s already exists", user.UserName)
			return utils.Conflict(c, err.Error())
		}
		logger.GetLogger("quiver").Errorf("create user failed: %v", err)
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, 0, "success", user)
}

func (h *UserHandler) ListUser(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "100"))

	users, total, err := h.userService.ListUser(env, page, size)
	if err != nil {
		logger.GetLogger("quiver").Errorf("list user error: %v", err)
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, 0, "success", fiber.Map{
		"total": total,
		"page":  page,
		"size":  size,
		"users": users,
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	user, err := h.userService.GetUser(env, userID)
	if err != nil {
		logger.GetLogger("quiver").Errorf("user not found: %v", err)
		return utils.BadRequest(c, "user not found")
	}

	return utils.Success(c, 0, "success", user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		logger.GetLogger("quiver").Errorf("invalid request body: %v", err)
		return utils.BadRequest(c, "invalid request body")
	}

	user, err := h.userService.UpdateUser(env, userID, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.GetLogger("quiver").Errorf("user not found: %v", err)
			return utils.NotFound(c, "user not found")
		}
		logger.GetLogger("quiver").Errorf("error updating user: %v", err)
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, 0, "success", user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	env := c.Locals("env").(string)
	userID, err := services.GetUserID(c)
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return err
	}

	if err := h.userService.DeleteUser(env, userID); err != nil {
		if err.Error() == "user not found" {
			logger.GetLogger("quiver").Errorf("user not found: %v", err)
			return utils.NotFound(c, "user not found")
		}

		logger.GetLogger("quiver").Errorf("user delete %d failed: %v", userID, err)
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, 0, "success", fiber.Map{"user_id": userID})
}
