package router

import (
	"server/internal/middleware"
	"server/internal/product/handler"

	"github.com/gin-gonic/gin"
)

var secret = []byte("gee")

func RegisterRoutes(r *gin.Engine, userHandler *handler.UserHandler, commodityHandler *handler.CommodityHandler, cartHandler *handler.CartHandler) {
	v1 := r.Group("/v1")
	v1.POST("/login", userHandler.Login)
	v1.POST("/register", userHandler.Register)
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleWare(secret))
	auth.POST("/createCommodity", commodityHandler.CreateCommodity)
	auth.POST("/updateCommodity", commodityHandler.UpdateCommodity)
	auth.GET("/listCommodity", commodityHandler.ListCommodity)
	auth.DELETE("/deleteCommodity", commodityHandler.DeleteCommodity)
	auth.GET("/getCommodity", commodityHandler.FindCommodityByName)
	auth.POST("/addToCart", cartHandler.AddToCart)
}
