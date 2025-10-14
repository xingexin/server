package container

import (
	"server/config"
	"server/internal/product/handler"
	"server/internal/product/repository"
	"server/internal/product/service"
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
	container.Provide(repository.NewUserRepository)
	container.Provide(repository.NewCommodityRepository)
	container.Provide(repository.NewCartRepository)

	// 提供 Services
	container.Provide(service.NewUserService)
	container.Provide(service.NewCommodityService)
	container.Provide(service.NewCartService)

	// 提供 Handlers
	container.Provide(handler.NewUserHandler)
	container.Provide(handler.NewCommodityHandler)
	container.Provide(handler.NewCartHandler)

	// 提供 Gin Engine
	container.Provide(gin.Default)

	return container
}
