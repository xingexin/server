package container

import (
	"server/config"
	cartHandler "server/internal/product/cart/handler"
	cartRepo "server/internal/product/cart/repository"
	cartService "server/internal/product/cart/service"
	commodityHandler "server/internal/product/commodity/handler"
	commodityRepo "server/internal/product/commodity/repository"
	commodityService "server/internal/product/commodity/service"
	orderHandler "server/internal/product/order/handler"
	orderRepo "server/internal/product/order/repository"
	orderService "server/internal/product/order/service"
	"server/internal/product/scheduler"
	userHandler "server/internal/product/user/handler"
	userRepo "server/internal/product/user/repository"
	userService "server/internal/product/user/service"
	"server/pkg/db"
	myRedis "server/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

// BuildContainer 构建依赖注入容器
func BuildContainer() *dig.Container {
	container := dig.New()

	// 提供配置
	if err := container.Provide(config.LoadConfig); err != nil {
		log.Fatalf("Failed to provide config: %v", err)
	}

	// 提供数据库连接
	if err := container.Provide(func(cfg *config.Config) (*gorm.DB, error) {
		return db.InitDB(cfg)
	}); err != nil {
		log.Fatalf("Failed to provide database connection: %v", err)
	}

	// 提供 Redis 连接
	if err := container.Provide(func(cfg *config.Config) (*redis.Client, error) {
		return myRedis.InitRedis(cfg)
	}); err != nil {
		log.Fatalf("Failed to provide Redis connection: %v", err)
	}

	// 提供 Repositories
	if err := container.Provide(orderRepo.NewOrderDQRepository); err != nil {
		log.Fatalf("Failed to provide OrderDQRepository: %v", err)
	}
	if err := container.Provide(commodityRepo.NewRedisCommodityRepository); err != nil {
		log.Fatalf("Failed to provide redisCommodityRepository: %v", err)
	}
	if err := container.Provide(userRepo.NewUserRepository); err != nil {
		log.Fatalf("Failed to provide UserRepository: %v", err)
	}
	if err := container.Provide(commodityRepo.NewCommodityRepository); err != nil {
		log.Fatalf("Failed to provide CommodityRepository: %v", err)
	}
	if err := container.Provide(cartRepo.NewCartRepository); err != nil {
		log.Fatalf("Failed to provide CartRepository: %v", err)
	}
	if err := container.Provide(orderRepo.NewOrderRepository); err != nil {
		log.Fatalf("Failed to provide OrderRepository: %v", err)
	}

	// 提供 Services
	if err := container.Provide(orderService.NewOrderCancelService); err != nil {
		log.Fatalf("Failed to provide OrderCancelService: %v", err)
	}
	if err := container.Provide(userService.NewUserService); err != nil {
		log.Fatalf("Failed to provide UserService: %v", err)
	}
	if err := container.Provide(commodityService.NewCommodityService); err != nil {
		log.Fatalf("Failed to provide CommodityService: %v", err)
	}
	if err := container.Provide(cartService.NewCartService); err != nil {
		log.Fatalf("Failed to provide CartService: %v", err)
	}
	if err := container.Provide(orderService.NewOrderService); err != nil {
		log.Fatalf("Failed to provide OrderService: %v", err)
	}
	if err := container.Provide(commodityService.NewStockCacheService); err != nil {
		log.Fatalf("Failed to provide StockCacheService: %v", err)
	}

	// 提供 Scheduler
	if err := container.Provide(scheduler.NewOrderDQScheduler); err != nil {
		log.Fatalf("Failed to provide OrderDQScheduler: %v", err)
	}
	if err := container.Provide(scheduler.NewScheduler); err != nil {
		log.Fatalf("Failed to provide Scheduler: %v", err)
	}

	// 提供 Handlers
	if err := container.Provide(userHandler.NewUserHandler); err != nil {
		log.Fatalf("Failed to provide UserHandler: %v", err)
	}
	if err := container.Provide(commodityHandler.NewCommodityHandler); err != nil {
		log.Fatalf("Failed to provide CommodityHandler: %v", err)
	}
	if err := container.Provide(cartHandler.NewCartHandler); err != nil {
		log.Fatalf("Failed to provide CartHandler: %v", err)
	}
	if err := container.Provide(orderHandler.NewOrderHandler); err != nil {
		log.Fatalf("Failed to provide OrderHandler: %v", err)
	}

	// 提供 Gin Engine
	if err := container.Provide(gin.Default); err != nil {
		log.Fatalf("Failed to provide Gin Engine: %v", err)
	}

	return container
}
