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
	userHandler "server/internal/product/user/handler"
	userRepo "server/internal/product/user/repository"
	userService "server/internal/product/user/service"
	"server/pkg/db"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

// BuildContainer 构建依赖注入容器
func BuildContainer() *dig.Container {
	container := dig.New()

	// 提供配置
	container.Provide(config.LoadConfig)

	// 提供数据库连接
	container.Provide(func(cfg *config.Config) (*gorm.DB, error) {
		return db.InitDB(cfg)
	})

	// 提供 Repositories
	container.Provide(userRepo.NewUserRepository)
	container.Provide(commodityRepo.NewCommodityRepository)
	container.Provide(cartRepo.NewCartRepository)
	container.Provide(orderRepo.NewOrderRepository)

	// 提供 Services
	container.Provide(userService.NewUserService)
	container.Provide(commodityService.NewCommodityService)
	container.Provide(cartService.NewCartService)
	container.Provide(orderService.NewOrderService)

	// 提供 Handlers
	container.Provide(userHandler.NewUserHandler)
	container.Provide(commodityHandler.NewCommodityHandler)
	container.Provide(cartHandler.NewCartHandler)
	container.Provide(orderHandler.NewOrderHandler)

	// 提供 Gin Engine
	container.Provide(gin.Default)

	return container
}
