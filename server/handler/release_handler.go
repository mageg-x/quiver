package handler

import (
	"github.com/gofiber/fiber/v2"
	"quiver/logger"
	"quiver/services"
	"quiver/utils"
	"strconv"
)

// ReleaseHandler 命名空间控制器
type ReleaseHandler struct {
	releaseService *services.ReleaseService
}

// NewReleaseHandler 创建命名空间控制器实例
func NewReleaseHandler() *ReleaseHandler {
	return &ReleaseHandler{
		releaseService: services.NewReleaseService(),
	}
}

func (c *ReleaseHandler) PublishRelease(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")
	namespaceName := ctx.Params("namespace_name")

	valid, err := services.CheckACNKFormat(ctx, &env, &appName, &clusterName, &namespaceName, nil)
	if !valid {
		return err
	}

	var body map[string]string
	if err := ctx.BodyParser(&body); err != nil {
		logger.GetLogger("quiver").Error("invalid request body")
		return utils.BadRequest(ctx, "invalid request body")
	}

	operator, comment, releaseName := body["operator"], body["comment"], body["release_name"]

	// 验证输入
	if len(operator) == 0 || len(releaseName) == 0 {
		return utils.BadRequest(ctx, "operator  and release_name is required")
	}

	release, err := c.releaseService.PublishRelease(env, appName, clusterName, namespaceName, releaseName, operator, comment)
	if err != nil {
		return utils.BadRequest(ctx, err.Error())
	}

	response := fiber.Map{
		"env":            env,
		"app_name":       appName,
		"cluster_name":   clusterName,
		"namespace_name": namespaceName,
		"release_name":   release.ReleaseName,
		"release_id":     release.ReleaseID,
		"release_time":   release.ReleaseTime,
	}

	return utils.Success(ctx, 0, "success", response)
}

func (c *ReleaseHandler) ListReleases(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")
	namespaceName := ctx.Params("namespace_name")

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, &clusterName, &namespaceName, nil)
	if !valid {
		return err
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	size, _ := strconv.Atoi(ctx.Query("size", "100"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 100
	}

	releases, total, err := c.releaseService.ListRelease(env, appName, clusterName, namespaceName, page, size)
	if err != nil {
		logger.GetLogger("quiver").Errorf("releases list failed %s", err.Error())
		return utils.InternalError(ctx, err.Error())
	}

	response := fiber.Map{
		"env":      env,
		"total":    total,
		"page":     page,
		"size":     size,
		"releases": releases,
	}

	return utils.Success(ctx, 0, "success", response)
}

func (c *ReleaseHandler) GetRelease(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")
	namespaceName := ctx.Params("namespace_name")
	release_id := ctx.Params("release_id")

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, &clusterName, &namespaceName, nil)
	if !valid {
		return err
	}
	release, err := c.releaseService.GetRelease(env, appName, clusterName, namespaceName, release_id)
	if err != nil {
		return utils.BadRequest(ctx, err.Error())
	}

	response := release

	return utils.Success(ctx, 0, "success", response)
}

func (c *ReleaseHandler) RollbackRelease(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")
	namespaceName := ctx.Params("namespace_name")
	releaseId := ctx.Params("release_id")

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, &clusterName, &namespaceName, nil)
	if !valid {
		return err
	}

	var body map[string]string
	if err := ctx.BodyParser(&body); err != nil {
		logger.GetLogger("quiver").Error("invalid request body")
		return utils.BadRequest(ctx, "invalid request body")
	}

	operator, comment := body["operator"], body["comment"]

	// 验证输入
	if len(operator) == 0 {
		return utils.BadRequest(ctx, "operator is required")
	}

	release, err := c.releaseService.RollbackRelease(env, appName, clusterName, namespaceName, releaseId, operator, comment)
	if err != nil {
		return utils.BadRequest(ctx, err.Error())
	}

	response := fiber.Map{
		"env":            env,
		"app_name":       appName,
		"cluster_name":   clusterName,
		"namespace_name": namespaceName,
		"release_name":   release.ReleaseName,
		"release_id":     release.ReleaseID,
		"release_time":   release.ReleaseTime,
		"operator":       operator,
		"comment":        comment,
	}

	return utils.Success(ctx, 0, "success", response)
}
