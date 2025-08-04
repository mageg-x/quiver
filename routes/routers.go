package routes

import (
	"quiver/handler"
	"quiver/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes 设置所有路由
func SetupRoutes(app *fiber.App) {
	// API v1 路由组
	api := app.Group("/api/v1")

	// 应用限流中间件
	api.Use(middleware.RateLimitMiddleware())

	// 环境路由组
	envGroup := api.Group("/envs/:env")

	//在 envGroup 上应用 EnvMiddleware
	envGroup.Use(middleware.EnvMiddleware())

	// 应用管理路由
	apps := envGroup.Group("/apps")
	{
		AppHandler := handler.NewAppHandler()
		apps.Post("/", AppHandler.CreateApp)            // 创建应用
		apps.Get("/", AppHandler.ListApp)               // 获取应用列表
		apps.Get("/:app_name", AppHandler.GetApp)       // 获取单个应用
		apps.Put("/:app_name", AppHandler.UpdateApp)    // 更新应用
		apps.Delete("/:app_name", AppHandler.DeleteApp) // 删除应用
	}

	// 集群管理路由
	clusters := apps.Group("/:app_name/clusters")
	{
		ClusterHandler := handler.NewClusterHandler()
		clusters.Post("/", ClusterHandler.CreateCluster)                // 创建集群
		clusters.Get("/", ClusterHandler.ListCluster)                   // 获取集群列表
		clusters.Get("/:cluster_name", ClusterHandler.GetCluster)       // 获取单个集群
		clusters.Delete("/:cluster_name", ClusterHandler.DeleteCluster) // 删除集群
	}

	// 命名空间管理路由
	namespaces := clusters.Group("/:cluster_name/namespaces")
	{
		namespaceController := handler.NewNamespaceHandler()
		namespaces.Post("/", namespaceController.CreateNamespace)                  // 创建命名空间
		namespaces.Get("/", namespaceController.ListNamespace)                     // 获取命名空间列表
		namespaces.Get("/:namespace_name", namespaceController.GetNamespace)       // 获取单个命名空间
		namespaces.Delete("/:namespace_name", namespaceController.DeleteNamespace) // 删除命名空间
	}

	// Item 管理路由
	items := namespaces.Group("/:namespace_name/items")
	{
		itemController := handler.NewItemHandler()
		items.Post("/", itemController.SetItem)
		items.Get("/:key", itemController.GetItem)
		items.Delete("/:key", itemController.DeleteItem)
	}
	// 404处理
	app.Use(middleware.NotFoundHandler)
}
