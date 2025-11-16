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

// main 主函数，应用程序入口
// 应用架构：
// 1. 使用dig进行依赖注入管理
// 2. 四个goroutine并发运行：
//    - HTTP服务器（处理API请求）
//    - 库存同步调度器（每10秒同步Redis库存到MySQL）
//    - 订单延迟队列调度器（每10秒处理超时订单）
//    - 超时任务恢复调度器（每60秒恢复processing队列中的超时任务）
// 3. 支持优雅关闭（监听SIGINT和SIGTERM信号）
func main() {
	// 创建dig容器，集中管理所有依赖的生命周期
	// BuildContainer会按顺序提供：配置、数据库、Redis、服务、Handler等
	c := container.BuildContainer()

	// 使用Invoke运行应用，dig会自动解析并注入所有依赖
	// 参数列表中的所有对象都由dig容器自动提供
	err := c.Invoke(func(
		cfg *config.Config,                        // 配置对象
		rdb *redis.Client,                         // Redis客户端
		gormDB *gorm.DB,                           // GORM数据库连接
		r *gin.Engine,                             // Gin Web引擎
		uHandler *userHandler.UserHandler,         // 用户Handler
		cHandler *commodityHandler.CommodityHandler, // 商品Handler
		caHandler *cartHandler.CartHandler,        // 购物车Handler
		oHandler *orderHandler.OrderHandler,       // 订单Handler
		stockScheduler *scheduler.Scheduler,       // 库存同步调度器
		orderDQScheduler *scheduler.OrderDQScheduler, // 订单延迟队列调度器
		recoveryScheduler *scheduler.RecoveryScheduler, // 超时任务恢复调度器
	) error {
		// 1. 初始化日志系统（根据配置文件设置日志级别）
		logger.InitLogger(cfg.Logger.Level)

		// 2. 获取底层SQL DB连接，用于资源清理
		sqlDB, err := gormDB.DB()
		if err != nil {
			return err
		}
		// 注册数据库连接的清理函数（程序退出时自动调用）
		defer func(sqlDB *sql.DB) {
			err = sqlDB.Close()
			if err != nil {
				log.Error("Failed to close database:", err)
			}
		}(sqlDB)

		// 3. 配置Gin中间件
		r.Use(gin.LoggerWithWriter(log.StandardLogger().Out)) // 使用logrus作为日志输出
		r.Use(gin.Recovery())                                 // panic恢复中间件

		// 4. 注册所有HTTP路由（包括公开路由和需要认证的路由）
		router.RegisterRoutes(r, uHandler, cHandler, caHandler, oHandler)

		// 5. 启动库存同步调度器（在独立goroutine中运行）
		// 作用：每10秒将Redis中的库存变化批量同步到MySQL
		go func() {
			log.Info("Starting stock sync scheduler...")
			if err := stockScheduler.Start(); err != nil {
				log.Error("Scheduler error:", err)
			}
		}()

		// 6. 启动订单延迟队列调度器（在独立goroutine中运行）
		// 作用：每10秒扫描超时订单（15分钟未支付），自动取消并归还库存
		go func() {
			log.Info("Starting Order DQ Scheduler...")
			if err := orderDQScheduler.Start(); err != nil {
				log.Error("Order DQ Scheduler error:", err)
			}
		}()

		// 7. 启动超时任务恢复调度器（在独立goroutine中运行）
		// 作用：每60秒扫描processing队列，将超时任务（处理超过5分钟）移回ready队列重试
		// 场景：调度器崩溃、处理失败、网络超时等导致任务卡在processing队列
		go func() {
			log.Info("Starting Recovery Scheduler...")
			if err := recoveryScheduler.Start(); err != nil {
				log.Error("Recovery Scheduler error:", err)
			}
		}()

		// 8. 设置系统信号监听（用于优雅关闭）
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 监听Ctrl+C和kill信号

		// 8. 启动HTTP服务器（在独立goroutine中运行）
		go func() {
			log.Info("Server is running at http://localhost:8080")
			if err := r.Run("localhost:" + strconv.Itoa(cfg.Server.Port)); err != nil {
				log.Error("Server error:", err)
			}
		}()

		// 9. 阻塞等待退出信号
		<-quit
		log.Info("Shutting down gracefully...")

		// 10. 优雅关闭：停止调度器
		stockScheduler.Stop()
		recoveryScheduler.Stop()
		// TODO: 也应该停止orderDQScheduler（需要添加Stop方法）

		log.Info("Server stopped")
		return nil
	})

	// 如果依赖注入失败，程序panic
	if err != nil {
		panic(err)
	}
}
