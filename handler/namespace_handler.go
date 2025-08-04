package handler

import (
	"github.com/gofiber/fiber/v2"
	"quiver/logger"
	"quiver/models"
	"quiver/services"
	"quiver/utils"
	"strconv"
)

// NamespaceHandler 命名空间控制器
type NamespaceHandler struct {
	namespaceService *services.NamespaceService
}

// NewNamespaceHandler 创建命名空间控制器实例
func NewNamespaceHandler() *NamespaceHandler {
	return &NamespaceHandler{
		namespaceService: services.NewNamespaceService(),
	}
}

// CreateNamespace 创建命名空间
func (c *NamespaceHandler) CreateNamespace(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言

	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")

	if !utils.ValidateAppName(appName) {
		logger.GetLogger("quiver").Errorf("invalid app_name %s", appName)
		return utils.BadRequest(ctx, "invalid app_name")
	}

	if !utils.ValidateClusterName(clusterName) {
		logger.GetLogger("quiver").Errorf("invalid cluster_name %s", clusterName)
		return utils.BadRequest(ctx, "invalid cluster_name")
	}

	var namespace models.Namespace
	if err := ctx.BodyParser(&namespace); err != nil {
		logger.GetLogger("quiver").Errorf("invalid request body")
		return utils.BadRequest(ctx, "invalid request body")
	}

	// 验证命名空间名称
	if !utils.ValidateNamespaceName(namespace.NamespaceName) {
		logger.GetLogger("quiver").Errorf("invalid namespace_name %s", namespace.NamespaceName)
		return utils.BadRequest(ctx, "Invalid namespace_name")
	}

	namespace.AppName = appName
	namespace.ClusterName = clusterName

	if err := c.namespaceService.CreateNamespace(env, &namespace); err != nil {
		if err.Error() == "app not found" {
			logger.GetLogger("quiver").Errorf("app not found %s", appName)
			return utils.NotFound(ctx, err.Error())
		}

		if err.Error() == "cluster not found" {
			logger.GetLogger("quiver").Errorf("cluster not found %s", clusterName)
			return utils.NotFound(ctx, err.Error())
		}
		if err.Error() == "namespace already exists" {
			logger.GetLogger("quiver").Errorf("namespace already exists %s", namespace.NamespaceName)
			return utils.BadRequest(ctx, err.Error())
		}
		return utils.InternalError(ctx, err.Error())
	}

	response := fiber.Map{
		"env":       env,
		"namespace": namespace,
	}

	return utils.Success(ctx, 0, "success", response)
}

// ListNamespace 获取命名空间列表
func (c *NamespaceHandler) ListNamespace(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	size, _ := strconv.Atoi(ctx.Query("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}

	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")

	if !utils.ValidateAppName(appName) {
		logger.GetLogger("quiver").Errorf("invalid app_name %s", appName)
		return utils.BadRequest(ctx, "invalid app_name")
	}

	if !utils.ValidateClusterName(clusterName) {
		logger.GetLogger("quiver").Errorf("invalid cluster_name %s", clusterName)
		return utils.BadRequest(ctx, "Invalid cluster_name")
	}

	namespaces, total, err := c.namespaceService.ListNamespace(env, appName, clusterName, page, size)
	if err != nil {
		if err.Error() == "app not found" {
			logger.GetLogger("quiver").Errorf("app not found %s", appName)
			return utils.NotFound(ctx, err.Error())
		}
		if err.Error() == "cluster not found" {
			logger.GetLogger("quiver").Errorf("cluster not found %s", clusterName)
			return utils.NotFound(ctx, err.Error())
		}

		logger.GetLogger("quiver").Errorf("list namespace failed %s", err.Error())
		return utils.InternalError(ctx, err.Error())
	}

	dtos := make([]map[string]interface{}, len(namespaces))
	for i, ns := range namespaces {
		dtos[i] = map[string]interface{}{
			"namespace_name": ns.NamespaceName,
			"description":    ns.Description,
			"create_time":    ns.CreateTime,
			"update_time":    ns.UpdateTime,
		}
	}

	response := fiber.Map{
		"env":          env,
		"app_name":     appName,
		"cluster_name": clusterName,
		"total":        total,
		"page":         page,
		"size":         size,
		"namespaces":   dtos,
	}

	return utils.Success(ctx, 0, "success", response)
}

// GetNamespace 获取单个命名空间
func (c *NamespaceHandler) GetNamespace(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	//env := ctx.Params("env")
	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")
	namespaceName := ctx.Params("namespace_name")

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
		return utils.BadRequest(ctx, "Invalid namespace_name")
	}

	namespace, err := c.namespaceService.GetNamespace(env, appName, clusterName, namespaceName)
	if err != nil {
		if err.Error() == "app not found" {
			logger.GetLogger("quiver").Errorf("app not found %s", appName)
			return utils.NotFound(ctx, err.Error())
		}

		if err.Error() == "cluster not found" {
			logger.GetLogger("quiver").Errorf("cluster not found %s", clusterName)
			return utils.NotFound(ctx, err.Error())
		}

		if err.Error() == "namespace not found" {
			logger.GetLogger("quiver").Errorf("namespace not found %s", namespaceName)
			return utils.NotFound(ctx, err.Error())
		}

		logger.GetLogger("quiver").Errorf("get namespace %s failed %s",
			appName+"/"+clusterName+"/"+namespaceName, err.Error())

		return utils.BadRequest(ctx, err.Error())
	}

	response := fiber.Map{
		"env":            env,
		"app_name":       namespace.AppName,
		"cluster_name":   namespace.ClusterName,
		"namespace_name": namespace.NamespaceName,
		"description":    namespace.Description,
		"create_time":    namespace.CreateTime,
		"update_time":    namespace.UpdateTime,
	}

	return utils.Success(ctx, 0, "success", response)
}

// DeleteNamespace 删除命名空间
func (c *NamespaceHandler) DeleteNamespace(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言

	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")
	namespaceName := ctx.Params("namespace_name")

	if !utils.ValidateAppName(appName) {
		logger.GetLogger("quiver").Errorf("invalid app_name %s", appName)
		return utils.BadRequest(ctx, "invalid app_name")
	}

	if !utils.ValidateClusterName(clusterName) {
		logger.GetLogger("quiver").Errorf("invalid cluster_name %s", clusterName)
		return utils.BadRequest(ctx, "Invalid cluster_name")
	}

	if !utils.ValidateNamespaceName(namespaceName) {
		logger.GetLogger("quiver").Errorf("invalid namespace_name %s", namespaceName)
		return utils.BadRequest(ctx, "invalid namespace_name")
	}

	if err := c.namespaceService.DeleteNamespace(env, appName, clusterName, namespaceName); err != nil {
		if err.Error() == "namespace not found" {
			logger.GetLogger("quiver").Errorf("namespace not found %s",
				appName+"/"+clusterName+"/"+namespaceName)
			return utils.NotFound(ctx, err.Error())
		}

		logger.GetLogger("quiver").Errorf("delete namespace %s failed %s",
			appName+"/"+clusterName+"/"+namespaceName, err.Error())
		return utils.InternalError(ctx, err.Error())
	}

	return utils.Success(ctx, 0, "success", fiber.Map{
		"env":            env,
		"app_name":       appName,
		"cluster_name":   clusterName,
		"namespace_name": namespaceName,
	})
}
