package main

import (
	"fmt"
	"quiver/config"
	"quiver/database"
	LOG "quiver/logger"
	"quiver/middleware"
	"quiver/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	log := LOG.GetLogger("quiver")
	log.Info("Quiver starting...")

	// 加载配置
	conf, _ := config.LoadConfig("config/config.yaml")
	log.Infof("get config : %+v", conf)

	// 创建 Fiber 应用
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		BodyLimit:    10 * 1024 * 1024, // 10MB
	})

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

	// Api 文档
	app.Static("/docs", "./docs")

	// 健康检查
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// 设置路由
	routes.SetupRoutes(app)

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
