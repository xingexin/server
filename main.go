package main

import (
	"database/sql"
	"os"
	"os/signal"
	"server/config"
	cartHandler "server/internal/product/cart/handler"
	commodityHandler "server/internal/product/commodity/handler"
	orderHandler "server/internal/product/order/handler"
	"server/internal/product/scheduler"
	userHandler "server/internal/product/user/handler"
	"syscall"

	"server/internal/router"
	"server/pkg/container"
	"server/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func main() {
	// 创建 dig 容器
	c := container.BuildContainer()

	// 使用 Invoke 来运行应用，dig 会自动解析所有依赖
	err := c.Invoke(func(
		cfg *config.Config,
		rdb *redis.Client,
		gormDB *gorm.DB,
		r *gin.Engine,
		uHandler *userHandler.UserHandler,
		cHandler *commodityHandler.CommodityHandler,
		caHandler *cartHandler.CartHandler,
		oHandler *orderHandler.OrderHandler,
		stockScheduler *scheduler.Scheduler,
		orderDQScheduler *scheduler.OrderDQScheduler,
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
		router.RegisterRoutes(r, uHandler, cHandler, caHandler, oHandler)

		// 启动Scheduler（在独立goroutine中）
		go func() {
			log.Info("Starting stock sync scheduler...")
			if err := stockScheduler.Start(); err != nil {
				log.Error("Scheduler error:", err)
			}
		}()
		go func() {
			log.Info("Starting Order DQ Scheduler...")
			if err := orderDQScheduler.Start(); err != nil {
				log.Error("Order DQ Scheduler error:", err)
			}
		}()
		// 监听退出信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// 在goroutine中启动HTTP服务器
		go func() {
			log.Info("Server is running at http://localhost:8080")
			if err := r.Run("localhost:" + strconv.Itoa(cfg.Server.Port)); err != nil {
				log.Error("Server error:", err)
			}
		}()

		// 等待退出信号
		<-quit
		log.Info("Shutting down gracefully...")

		// 停止Scheduler
		stockScheduler.Stop()

		log.Info("Server stopped")
		return nil
	})

	if err != nil {
		panic(err)
	}
}
