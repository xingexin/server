package router

import (
	"server/internal/middleware"
	"server/internal/product"

	"github.com/gin-gonic/gin"
)

var secret = []byte("gee")

func RegisterRoutes(r *gin.Engine, handler *product.Handler) {
	v1 := r.Group("/v1")
	v1.POST("/login", handler.Login)
	v1.POST("/register", handler.Register)
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleWare(secret))
	auth.POST("/createCommodity", handler.CreateCommodity)
	auth.POST("/updateCommodity", handler.UpdateCommodity)
	auth.GET("/listCommodity", handler.ListCommodity)
	auth.DELETE("/deleteCommodity", handler.DeleteCommodity)
	auth.GET("/getCommodity", handler.FindCommodityByName)
}
