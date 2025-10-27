package router

import (
	"server/internal/middleware"
	cartHandler "server/internal/product/cart/handler"
	commodityHandler "server/internal/product/commodity/handler"
	orderHandler "server/internal/product/order/handler"
	userHandler "server/internal/product/user/handler"

	"github.com/gin-gonic/gin"
)

var secret = []byte("gee")

func RegisterRoutes(r *gin.Engine, uHandler *userHandler.UserHandler, cHandler *commodityHandler.CommodityHandler, caHandler *cartHandler.CartHandler, oHandler *orderHandler.OrderHandler) {
	v1 := r.Group("/v1")
	v1.POST("/login", uHandler.Login)
	v1.POST("/register", uHandler.Register)
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleWare(secret))
	auth.POST("/createCommodity", cHandler.CreateCommodity)
	auth.POST("/updateCommodity", cHandler.UpdateCommodity)
	auth.GET("/listCommodity", cHandler.ListCommodity)
	auth.DELETE("/deleteCommodity", cHandler.DeleteCommodity)
	auth.GET("/getCommodity", cHandler.FindCommodityByName)

	auth.POST("/addToCart", caHandler.AddToCart)
	auth.DELETE("/removeFromCart", caHandler.RemoveFromCart)
	auth.PUT("/updateCart", caHandler.UpdateCart)
	auth.GET("/getCart", caHandler.GetCart)

	auth.POST("/createOrder", oHandler.CreateOrder)
	auth.PUT("/updateOrder", oHandler.UpdateOrderStatus)
	auth.DELETE("/deleteOrder", oHandler.DeleteOrder)
	auth.GET("/getOrder", oHandler.GetOrder)
}
