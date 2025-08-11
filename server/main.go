package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"quiver/cache"
	"quiver/config"
	"quiver/database"
	LOG "quiver/logger"
	"quiver/middleware"
	"quiver/routes"
	"quiver/services"
	"quiver/utils"
	"quiver/web"
	"strings"
	"time"
)

func isStaticFile(path string) bool {
	// 常见静态资源扩展名
	staticExts := []string{".js", ".css", ".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg", ".woff", ".woff2", ".ttf", ".pdf"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range staticExts {
		if ext == e {
			return true
		}
	}
	return false
}

func main() {
	log := LOG.GetLogger("quiver")
	log.Info("Quiver starting...")

	// 加载配置
	conf, _ := config.LoadConfig("config/config.yaml")
	log.Infof("get config : %+v", conf)

	// 初始化缓存
	err := cache.Init(cache.Options{MaxMemCost: 1 << 30, DiskDir: os.TempDir() + "/quiver/cache"})
	if err != nil {
		log.Warnf("init cache error: %v", err)
	}
	err = cache.AddWatch([]cache.WatchTable{
		{Env: "dev", Name: "accesskey", Interval: 5 * time.Second, KeyUpdateCB: services.OnKeyUpdate4AccessKey},
		{Env: "dev", Name: "permission", Interval: 5 * time.Second, KeyUpdateCB: services.OnKeyUpdate4Permission},
		{Env: "dev", Name: "namespace_release", Interval: time.Second, KeyUpdateCB: services.OnKeyUpdate4Release},

		{Env: "pro", Name: "accesskey", Interval: 5 * time.Second, KeyUpdateCB: services.OnKeyUpdate4AccessKey},
		{Env: "pro", Name: "permission", Interval: 5 * time.Second, KeyUpdateCB: services.OnKeyUpdate4Permission},
		{Env: "pro", Name: "namespace_release", Interval: time.Second, KeyUpdateCB: services.OnKeyUpdate4Release},
	})

	if err != nil {
		log.Warnf("add watch error: %v", err)
	}

	// 创建 Fiber 应用
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		BodyLimit:    10 * 1024 * 1024, // 10MB
	})

	// 开启压缩（默认是 gzip）
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 压缩等级：fastest, default, best-compression
	}))

	// 中间件
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// 设置路由
	routes.SetupRoutes(app)

	// Api 文档
	app.Static("/docs", "./docs")

	// 健康检查
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// 前端静态文件
	dist, err := fs.Sub(web.WebDistFS, "dist")
	if err != nil {
		log.Fatalf("Failed to get sub filesystem: %v", err)
	}
	httpFS := http.FS(dist)
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   httpFS,
		Browse: false,
	}))

	// 如果上面没匹配到（即 404），返回 index.html
	app.Use(func(c *fiber.Ctx) error {
		// 只对非 API、非静态资源的路径返回 index.html
		path := c.Path()
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/docs") {
			return utils.NotFound(c, "Route not found")
		}
		if isStaticFile(path) {
			return utils.NotFound(c, "Route not found")
		}

		// 否则返回 SPA 入口
		file, err := dist.Open("index.html")
		if err != nil {
			return utils.BadRequest(c, "Failed to open index.html")
		}
		defer file.Close()

		content, _ := io.ReadAll(file)
		c.Set("Content-Type", "text/html; charset=utf-8")
		c.Status(200)
		return c.Send(content)
	})

	// 尝试连接数据库
	database.GetDB("dev")
	database.GetDB("pro")

	// 启动服务器
	host := config.GetServerConfig().Host
	port := config.GetServerConfig().Port
	log.Infof("Server starting on  %s:%d", host, port)

	if err := app.Listen(fmt.Sprintf("%s:%d", host, port)); err != nil {
		log.Fatal("Server failed to start:", err)
	}

}
