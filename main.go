package main

import (
	"database/sql"
	"server/config"
	"server/internal/product/handler"
	"server/internal/router"
	"server/pkg/container"
	"server/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func main() {
	// 创建 dig 容器
	c := container.BuildContainer()

	// 使用 Invoke 来运行应用，dig 会自动解析所有依赖
	err := c.Invoke(func(
		cfg *config.Config,
		gormDB *gorm.DB,
		r *gin.Engine,
		userHandler *handler.UserHandler,
		commodityHandler *handler.CommodityHandler,
		cartHandler *handler.CartHandler,
		orderHandler *handler.OrderHandler,
	) error {
		// 初始化日志
		logger.InitLogger(cfg.Logger.Level)

		// 获取 SQL DB 用于资源清理
		sqlDB, err := gormDB.DB()
		if err != nil {
			return err
		}
		defer func(sqlDB *sql.DB) {
			err = sqlDB.Close()
			if err != nil {
				log.Error("Failed to close database:", err)
			}
		}(sqlDB)

		r.Use(gin.LoggerWithWriter(log.StandardLogger().Out))
		r.Use(gin.Recovery())

		// 注册路由
		router.RegisterRoutes(r, userHandler, commodityHandler, cartHandler, orderHandler)

		log.Info("Server is running at http://localhost:8080")

		// 启动服务器
		return r.Run("localhost:" + strconv.Itoa(cfg.Server.Port))
	})

	if err != nil {
		panic(err)
	}
}
