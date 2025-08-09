package handler

import (
	"github.com/gofiber/fiber/v2"
	"quiver/logger"
	"quiver/services"
	"quiver/utils"
	"strconv"
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

	var request struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		logger.GetLogger("quiver").Errorf("invalid request body")
		return utils.BadRequest(ctx, "invalid request body")
	}

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, &clusterName, &namespaceName, &request.Key)
	if !valid {
		return err
	}

	if len(request.Value) == 0 {
		logger.GetLogger("quiver").Errorf("invalid item value")
		return utils.BadRequest(ctx, "invalid item value")
	}

	err = c.itemService.SetItem(env, appName, clusterName, namespaceName, request.Key, request.Value)
	if err != nil {
		if err.Error() == "app not found" {
			logger.GetLogger("quiver").Errorf("app not found")
			return utils.NotFound(ctx, err.Error())
		}

		if err.Error() == "cluster not found" {
			logger.GetLogger("quiver").Errorf("cluster not found %s", clusterName)
			return utils.NotFound(ctx, err.Error())
		}

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

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, &clusterName, &namespaceName, &key)
	if !valid {
		return err
	}

	item, err := c.itemService.GetItem(env, appName, clusterName, namespaceName, key)
	if err != nil {
		if err.Error() == "app not found" {
			logger.GetLogger("quiver").Errorf("app not found")
			return utils.NotFound(ctx, err.Error())
		}

		if err.Error() == "cluster not found" {
			logger.GetLogger("quiver").Errorf("cluster not found %s", clusterName)
			return utils.NotFound(ctx, err.Error())
		}

		if err.Error() == "namespace not found" {
			logger.GetLogger("quiver").Errorf("namespace not found")
			return utils.NotFound(ctx, err.Error())
		}

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
		"value":          item.V,
		"create_time":    item.CreateTime,
		"update_time":    item.UpdateTime,
	}

	return utils.Success(ctx, 0, "success", response)
}

func (c *ItemHandler) ListItem(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")
	namespaceName := ctx.Params("namespace_name")
	search := ctx.Query("search")
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

	items, total, err := c.itemService.ListItem(env, appName, clusterName, namespaceName, search, page, size)
	if err != nil {
		logger.GetLogger("quiver").Errorf("items list failed %s", err.Error())
		return utils.InternalError(ctx, err.Error())
	}

	response := fiber.Map{
		"env":   env,
		"total": total,
		"page":  page,
		"size":  size,
		"items": items,
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

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, &clusterName, &namespaceName, &key)
	if !valid {
		return err
	}

	err = c.itemService.DeleteItem(env, appName, clusterName, namespaceName, key)
	if err != nil {
		if err.Error() == "app not found" {
			return utils.NotFound(ctx, err.Error())
		}

		if err.Error() == "cluster not found" {
			return utils.NotFound(ctx, err.Error())
		}

		if err.Error() == "namespace not found" {
			return utils.NotFound(ctx, err.Error())
		}

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
