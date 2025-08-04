package handler

import (
	"github.com/gofiber/fiber/v2"
	"quiver/logger"
	"quiver/services"
	"quiver/utils"
)

// ItemHandler 命名空间控制器
type ItemHandler struct {
	itemService *services.ItemService
}

// NewItemHandler 创建命名空间控制器实例
func NewItemHandler() *ItemHandler {
	return &ItemHandler{
		itemService: services.NewItemService(),
	}
}

// SetItem 设置配置项
func (c *ItemHandler) SetItem(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")
	namespaceName := ctx.Params("namespace_name")

	// 参数验证
	if !utils.ValidateAppName(appName) {
		logger.GetLogger("quiver").Errorf("invalid app_name %s", appName)
		return utils.BadRequest(ctx, "invalid app_name")
	}

	if !utils.ValidateClusterName(clusterName) {
		logger.GetLogger("quiver").Errorf("invalid cluster_name %s", clusterName)
		return utils.BadRequest(ctx, "invalid cluster_name")
	}

	if !utils.ValidateNamespaceName(namespaceName) {
		logger.GetLogger("quiver").Errorf("invalid namespace_name %s", namespaceName)
		return utils.BadRequest(ctx, "invalid namespace_name")
	}

	var request struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		logger.GetLogger("quiver").Errorf("invalid request body")
		return utils.BadRequest(ctx, "invalid request body")
	}

	if !utils.ValidateItemKey(request.Key) {
		logger.GetLogger("quiver").Errorf("invalid item key %s", request.Key)
		return utils.BadRequest(ctx, "invalid item key")
	}

	err := c.itemService.SetItem(env, appName, clusterName, namespaceName, request.Key, request.Value)
	if err != nil {
		if err.Error() == "namespace not found" {
			logger.GetLogger("quiver").Errorf("namespace not found")
			return utils.NotFound(ctx, err.Error())
		}

		logger.GetLogger("quiver").Errorf("set item item failed %s", err.Error())
		return utils.InternalError(ctx, err.Error())
	}

	response := fiber.Map{
		"env":            env,
		"app_name":       appName,
		"cluster_name":   clusterName,
		"namespace_name": namespaceName,
		"key":            request.Key,
		"value":          request.Value,
	}

	return utils.Success(ctx, 0, "success", response)
}

// GetItem 获取单个配置项
func (c *ItemHandler) GetItem(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言

	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")
	namespaceName := ctx.Params("namespace_name")
	key := ctx.Params("key")

	// 参数验证
	if !utils.ValidateAppName(appName) {
		logger.GetLogger("quiver").Errorf("invalid app_name %s", appName)
		return utils.BadRequest(ctx, "invalid app_name")
	}

	if !utils.ValidateClusterName(clusterName) {
		logger.GetLogger("quiver").Errorf("invalid cluster_name %s", clusterName)
		return utils.BadRequest(ctx, "invalid cluster_name")
	}

	if !utils.ValidateNamespaceName(namespaceName) {
		logger.GetLogger("quiver").Errorf("invalid namespace_name %s", namespaceName)
		return utils.BadRequest(ctx, "invalid namespace_name")
	}

	if !utils.ValidateItemKey(key) {
		logger.GetLogger("quiver").Errorf("invalid item key %s", key)
		return utils.BadRequest(ctx, "invalid item key")
	}

	value, err := c.itemService.GetItem(env, appName, clusterName, namespaceName, key)
	if err != nil {
		if err.Error() == "item not found" {
			return utils.NotFound(ctx, err.Error())
		}
		return utils.InternalError(ctx, err.Error())
	}

	response := fiber.Map{
		"env":            env,
		"app_name":       appName,
		"cluster_name":   clusterName,
		"namespace_name": namespaceName,
		"key":            key,
		"value":          value,
	}

	return utils.Success(ctx, 0, "success", response)
}

// DeleteItem 删除配置项
func (c *ItemHandler) DeleteItem(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言

	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")
	namespaceName := ctx.Params("namespace_name")
	key := ctx.Params("key")

	// 参数验证
	if !utils.ValidateAppName(appName) {
		logger.GetLogger("quiver").Errorf("invalid app_name %s", appName)
		return utils.BadRequest(ctx, "invalid app_name")
	}

	if !utils.ValidateClusterName(clusterName) {
		logger.GetLogger("quiver").Errorf("invalid cluster_name %s", clusterName)
		return utils.BadRequest(ctx, "invalid cluster_name")
	}

	if !utils.ValidateNamespaceName(namespaceName) {
		logger.GetLogger("quiver").Errorf("invalid namespace_name %s", namespaceName)
		return utils.BadRequest(ctx, "invalid namespace_name")
	}

	if !utils.ValidateItemKey(key) {
		logger.GetLogger("quiver").Errorf("invalid item key %s", key)
		return utils.BadRequest(ctx, "invalid item key")
	}

	err := c.itemService.DeleteItem(env, appName, clusterName, namespaceName, key)
	if err != nil {
		if err.Error() == "item not found" {
			return utils.NotFound(ctx, err.Error())
		}
		return utils.InternalError(ctx, err.Error())
	}

	response := fiber.Map{
		"env":            env,
		"app_name":       appName,
		"cluster_name":   clusterName,
		"namespace_name": namespaceName,
		"key":            key,
	}

	return utils.Success(ctx, 0, "success", response)
}
