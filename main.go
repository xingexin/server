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
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
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

	sqlDB, err := gormDB.DB()
	if err != nil {
		panic(err)
	}
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

	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"}, // 允许的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.LoggerWithWriter(log.StandardLogger().Out))
	r.Use(gin.Recovery())

	log.Info("Server is running at http://localhost:8080")

	router.RegisterRoutes(r, handler)
	err = r.Run("localhost:" + strconv.Itoa(cfg.Server.Port))
	if err != nil {
		panic(err)
	}

	return
}
