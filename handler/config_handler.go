package handler

import (
	"quiver/services"
	"quiver/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ConfigHandler 配置控制器
type ConfigHandler struct {
	configService *services.ConfigService
}

// NewConfigHandler 创建配置控制器实例
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		configService: services.NewConfigService(),
	}
}

// GetConfigs 获取命名空间配置
func (c *ConfigHandler) GetConfigs(ctx *fiber.Ctx) error {
	env := ctx.Params("env")
	appName := ctx.Params("appName")
	clusterName := ctx.Params("clusterName")
	namespaceName := ctx.Params("namespaceName")

	releaseKey := ctx.Query("releaseKey")
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	size, _ := strconv.Atoi(ctx.Query("size", "100"))

	// 参数验证
	if !utils.ValidateAppName(appName) {
		return utils.BadRequest(ctx, "Invalid app_name format")
	}

	if !utils.ValidateClusterName(clusterName) {
		return utils.BadRequest(ctx, "Invalid cluster_name format")
	}

	if !utils.ValidateNamespaceName(namespaceName) {
		return utils.BadRequest(ctx, "Invalid namespace_name format")
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 500 {
		size = 100
	}

	// 获取配置
	response, err := c.configService.GetNamespaceConfigs(env, appName, clusterName, namespaceName, releaseKey, page, size)
	if err != nil {
		if err.Error() == "no release found" {
			return utils.NotFound(ctx, "No configuration released")
		}
		return utils.InternalError(ctx, err.Error())
	}

	return utils.Success(ctx, 0, "success", response)
}

// ReleaseNamespace 发布命名空间
func (c *ConfigHandler) ReleaseNamespace(ctx *fiber.Ctx) error {
	return nil
}

// GetNotifications 获取配置变更通知（长轮询）
func (c *ConfigHandler) GetNotifications(ctx *fiber.Ctx) error {
	return nil
}
