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

	auth.POST("/commodity", cHandler.CreateCommodity)
	auth.PUT("/commodity/:id", cHandler.UpdateCommodity)
	auth.GET("/commodity", cHandler.ListCommodity)
	auth.DELETE("/commodity/:id", cHandler.DeleteCommodity)
	auth.GET("/commodity/search", cHandler.FindCommodityByName)

	auth.POST("/cart", caHandler.AddToCart)
	auth.DELETE("/cart/:id", caHandler.RemoveFromCart)
	auth.PUT("/cart/:id", caHandler.UpdateCart)
	auth.GET("/cart", caHandler.GetCart)

	auth.POST("/order", oHandler.CreateOrder)
	auth.PUT("/order/:id", oHandler.UpdateOrderStatus)
	auth.DELETE("/order/:id", oHandler.DeleteOrder)
	auth.GET("/order/:id", oHandler.GetOrder)
}
