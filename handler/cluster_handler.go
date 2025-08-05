package handler

import (
	"github.com/gofiber/fiber/v2"
	"quiver/logger"
	"quiver/models"
	"quiver/services"
	"quiver/utils"
	"strconv"
)

// ClusterHandler 集群控制器
type ClusterHandler struct {
	clusterService *services.ClusterService
}

// NewClusterHandler 创建集群控制器实例
func NewClusterHandler() *ClusterHandler {
	return &ClusterHandler{
		clusterService: services.NewClusterService(),
	}
}

// CreateCluster 创建集群
func (c *ClusterHandler) CreateCluster(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")

	var cluster models.Cluster
	if err := ctx.BodyParser(&cluster); err != nil {
		logger.GetLogger("quiver").Error("invalid request body")
		return utils.BadRequest(ctx, "invalid request body")
	}

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, &cluster.ClusterName, nil, nil)
	if !valid {
		return err
	}

	cluster.AppName = appName

	if err := c.clusterService.CreateCluster(env, &cluster); err != nil {
		if err.Error() == "app not found" {
			logger.GetLogger("quiver").Errorf("app not found %s", appName)
			return utils.NotFound(ctx, err.Error())
		}
		if err.Error() == "cluster already exists in this app" {
			logger.GetLogger("quiver").Errorf("cluster already exists in this app %s", cluster.ClusterName)
			return utils.BadRequest(ctx, err.Error())
		}
		return utils.InternalError(ctx, err.Error())
	}

	response := fiber.Map{
		"env":     env,
		"cluster": cluster,
	}

	return utils.Success(ctx, 0, "success", response)
}

// ListCluster 获取集群列表
func (c *ClusterHandler) ListCluster(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, nil, nil, nil)
	if !valid {
		return err
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	size, _ := strconv.Atoi(ctx.Query("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}

	if !utils.ValidateAppName(appName) {
		logger.GetLogger("quiver").Errorf("invalid app_name %s", appName)
		return utils.BadRequest(ctx, "invalid app_name format")
	}

	clusters, total, err := c.clusterService.ListCluster(env, appName, page, size)
	if err != nil {
		if err.Error() == "app not found" {
			logger.GetLogger("quiver").Errorf("app not found %s", appName)
			return utils.NotFound(ctx, err.Error())
		}

		logger.GetLogger("quiver").Errorf("app_name %s not found", appName)
		return utils.InternalError(ctx, err.Error())
	}

	dtos := make([]map[string]interface{}, len(clusters))
	for i, cluster := range clusters {
		dtos[i] = map[string]interface{}{
			"cluster_name": cluster.ClusterName,
			"description":  cluster.Description,
			"create_time":  cluster.CreateTime,
			"update_time":  cluster.UpdateTime,
		}
	}

	response := fiber.Map{
		"env":      env,
		"app_name": appName,
		"clusters": dtos,
		"total":    total,
		"page":     page,
		"size":     size,
	}

	return utils.Success(ctx, 0, "success", response)
}

// GetCluster 获取单个集群
func (c *ClusterHandler) GetCluster(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, &clusterName, nil, nil)
	if !valid {
		return err
	}

	cluster, err := c.clusterService.GetCluster(env, appName, clusterName)
	if err != nil {
		logger.GetLogger("quiver").Errorf("cluster not found %s", clusterName)
		return utils.NotFound(ctx, "Cluster not found")
	}

	response := fiber.Map{
		"env":          env,
		"app_name":     cluster.AppName,
		"cluster_name": cluster.ClusterName,
		"description":  cluster.Description,
		"create_time":  cluster.CreateTime,
		"update_time":  cluster.UpdateTime,
	}

	return utils.Success(ctx, 0, "success", response)
}

// DeleteCluster 删除集群
func (c *ClusterHandler) DeleteCluster(ctx *fiber.Ctx) error {
	env := ctx.Locals("env").(string) // 类型断言
	appName := ctx.Params("app_name")
	clusterName := ctx.Params("cluster_name")

	// 验证输入
	valid, err := services.CheckACNKFormat(ctx, &env, &appName, &clusterName, nil, nil)
	if !valid {
		return err
	}

	if err := c.clusterService.DeleteCluster(env, appName, clusterName); err != nil {
		if err.Error() == "cluster not found" {
			logger.GetLogger("quiver").Errorf("cluster not found %s", clusterName)
			return utils.NotFound(ctx, err.Error())
		}
		logger.GetLogger("quiver").Errorf("cluster delete failed %s", err.Error())
		return utils.InternalError(ctx, err.Error())
	}

	return utils.Success(ctx, 0, "success", fiber.Map{
		"env":          env,
		"app_name":     appName,
		"cluster_name": clusterName,
	})
}
