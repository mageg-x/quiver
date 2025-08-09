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
	// 用户管理 （只允许管理员）
	users := envGroup.Group("/users")
	{
		users.Post("/", handler.NewUserHandler().CreateUser)
		users.Get("/", handler.NewUserHandler().ListUser)
		users.Get("/:user_id", handler.NewUserHandler().GetUser)
		users.Put("/:user_id", handler.NewUserHandler().UpdateUser)
		users.Delete("/:user_id", handler.NewUserHandler().DeleteUser)
	}

	// 权限管理
	permissions := users.Group("/:user_id/permissions")
	{
		permissions.Post("/", handler.NewPermissionHandler().CreatePermission)
		permissions.Get("/", handler.NewPermissionHandler().ListPermission)
		permissions.Get("/:permission_id", handler.NewPermissionHandler().GetPermission)
		permissions.Put("/:permission_id", handler.NewPermissionHandler().UpdatePermission)
		permissions.Delete("/:permission_id", handler.NewPermissionHandler().DeletePermission)
	}

	// accesskey 管理
	accesskeys := users.Group("/:user_id/accesskeys")
	{
		accesskeys.Post("/", handler.NewAccessKeyHandler().CreateAccessKey)
		accesskeys.Get("/", handler.NewAccessKeyHandler().ListAccessKey)
		accesskeys.Get("/:accesskey", handler.NewAccessKeyHandler().GetAccessKey)
		accesskeys.Delete("/:accesskey", handler.NewAccessKeyHandler().DeleteAccessKey)
	}

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
		namespaces.Post("/", namespaceController.CreateNamespace)                     // 创建命名空间
		namespaces.Get("/", namespaceController.ListNamespace)                        // 获取命名空间列表
		namespaces.Get("/:namespace_name", namespaceController.GetNamespace)          // 获取单个命名空间
		namespaces.Delete("/:namespace_name", namespaceController.DeleteNamespace)    // 删除命名空间
		namespaces.Post("/:namespace_name/discard", namespaceController.DiscardDraft) // 丢弃草稿
	}

	// Item(草稿，尚未发布的编辑数据) 管理路由
	items := namespaces.Group("/:namespace_name/items")
	{
		itemController := handler.NewItemHandler()
		items.Post("/", itemController.SetItem)
		items.Get("/", itemController.ListItem)
		items.Get("/:key", itemController.GetItem)
		items.Delete("/:key", itemController.DeleteItem)
	}

	// 版本管理路由
	releases := namespaces.Group("/:namespace_name/releases")
	{
		releaseHandler := handler.NewReleaseHandler()
		releases.Post("/", releaseHandler.PublishRelease)       // 创建发布
		releases.Get("/", releaseHandler.ListReleases)          // 获取发布列表
		releases.Get("/:release_id", releaseHandler.GetRelease) // 获取发布详情
	}

	// 回滚
	rollback := namespaces.Group("/:namespace_name/rollback")
	{
		rollbackHandler := handler.NewReleaseHandler()
		rollback.Post("/:release_id", rollbackHandler.RollbackRelease) // 回滚发布
	}

	app.Use(middleware.NotFoundHandler)
}
