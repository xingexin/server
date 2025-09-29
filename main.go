package main

import (
	"database/sql"
	"server/internal/product"
	"server/internal/product/repository"
	"server/internal/product/service"
	"server/internal/router"
	"server/pkg/db"

	"github.com/gin-gonic/gin"
)

func main() {
	gormDB, err := db.InitDB()
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
	repo := repository.NewUserRepository(gormDB)
	srvc := service.NewUserService(repo)
	handler := product.NewHandler(srvc)
	r := gin.Default()
	router.RegisterRoutes(r, handler)
	err = r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
	return
}
