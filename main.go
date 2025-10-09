package main

import (
	"database/sql"
	"server/config"
	"server/internal/product"
	"server/internal/product/repository"
	"server/internal/product/service"
	"server/internal/router"
	"server/pkg/db"
	"server/pkg/logger"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	logger.InitLogger(cfg.Logger.Level)

	gormDB, err := db.InitDB(cfg)
	if err != nil {
		panic(err)
	}

	sqlDB, _ := gormDB.DB()
	defer func(sqlDB *sql.DB) {
		err = sqlDB.Close()
		if err != nil {
			panic(err)
		}
	}(sqlDB)

	uRepo := repository.NewUserRepository(gormDB)
	cRepo := repository.NewCommodityRepository(gormDB)

	cSvc := service.NewCommodityService(cRepo)
	uSvc := service.NewUserService(uRepo)

	handler := product.NewHandler(uSvc, cSvc)

	r := gin.Default()

	r.Use(gin.LoggerWithWriter(log.StandardLogger().Out))
	r.Use(gin.Recovery())

	log.Info("Server is running at http://localhost:8080")

	router.RegisterRoutes(r, handler)
	err = r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}

	return
}
