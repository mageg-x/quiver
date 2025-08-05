package handler

import (
	"github.com/gofiber/fiber/v2"
	"quiver/logger"
	"quiver/models"
	"quiver/services"
	"quiver/utils"
	"strconv"
)

// AppHandler 应用控制器
type AppHandler struct {
	appService *services.AppService
}

// NewAppHandler 创建应用控制器实例
func NewAppHandler() *AppHandler {
	return &AppHandler{
		appService: services.NewAppService(),
	}
}

// CreateApp 创建应用
func (c *AppHandler) CreateApp(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言

	var app models.App
	if err := ctx.BodyParser(&app); err != nil {
		logger.GetLogger("quiver").Error("invalid request body")
		return utils.BadRequest(ctx, "invalid request body")
	}

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &app.AppName, nil, nil, nil)
	if !valid {
		return err
	}

	if len(app.Description) > 500 {
		logger.GetLogger("quiver").Error("description is too long")
		return utils.BadRequest(ctx, "description is too long")
	}

	// 创建应用
	if err := c.appService.CreateApp(env, &app); err != nil {
		if err.Error() == "app_name already exists" {
			logger.GetLogger("quiver").Errorf("app_name already exists")
			return utils.Conflict(ctx, err.Error())
		}
		logger.GetLogger("quiver").Errorf("app create failed %s", err.Error())
		return utils.InternalError(ctx, err.Error())
	}
	logger.GetLogger("quiver").Infof("app created: %+v", app)

	return utils.Success(ctx, 0, "success", app)
}

// ListApp 获取应用列表
func (c *AppHandler) ListApp(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	size, _ := strconv.Atoi(ctx.Query("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}

	apps, total, err := c.appService.ListApp(env, page, size)
	if err != nil {
		logger.GetLogger("quiver").Errorf("app list failed %s", err.Error())
		return utils.InternalError(ctx, err.Error())
	}

	response := fiber.Map{
		"env":   env,
		"total": total,
		"page":  page,
		"size":  size,
		"apps":  apps,
	}

	return utils.Success(ctx, 0, "success", response)
}

// GetApp 获取单个应用
func (c *AppHandler) GetApp(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, nil, nil, nil)
	if !valid {
		return err
	}

	app, err := c.appService.GetApp(env, appName)
	if err != nil {
		logger.GetLogger("quiver").Errorf("app_name %s not found", appName)
		return utils.NotFound(ctx, "app_name not found")
	}

	response := fiber.Map{
		"env":         env,
		"app_name":    app.AppName,
		"description": app.Description,
		"create_time": app.CreateTime,
		"update_time": app.UpdateTime,
	}

	return utils.Success(ctx, 0, "success", response)
}

// UpdateApp 更新应用
func (c *AppHandler) UpdateApp(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, nil, nil, nil)
	if !valid {
		return err
	}

	var updates map[string]interface{}
	if err := ctx.BodyParser(&updates); err != nil {
		logger.GetLogger("quiver").Error("invalid request body")
		return utils.BadRequest(ctx, "invalid request body")
	}

	app, err := c.appService.UpdateApp(env, appName, updates)
	if err != nil {
		logger.GetLogger("quiver").Error("app update failed %s", err.Error())
		if err.Error() == "app_name not exist" {
			return utils.NotFound(ctx, "app_name not found")
		}
		return utils.InternalError(ctx, err.Error())
	}

	response := fiber.Map{
		"env":         env,
		"app_name":    app.AppName,
		"description": app.Description,
		"create_time": app.CreateTime,
		"update_time": app.UpdateTime,
	}
	return utils.Success(ctx, 0, "success", response)
}

// DeleteApp 删除应用
func (c *AppHandler) DeleteApp(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, nil, nil, nil)
	if !valid {
		return err
	}

	if err := c.appService.DeleteApp(env, appName); err != nil {
		if err.Error() == "app_name not exist" {
			return utils.NotFound(ctx, "app_name not found")
		}
		logger.GetLogger("quiver").Errorf("app delete failed %s", err.Error())
		return utils.InternalError(ctx, err.Error())
	}

	return utils.Success(ctx, 0, "success", fiber.Map{
		"env":      env,
		"app_name": appName,
	})
}
