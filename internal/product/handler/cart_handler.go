package handler

import (
	"net/http"
	"server/internal/product/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CartHandler struct {
	cartService *service.CartService
	userService *service.UserService
}

type reqAdd struct {
	UserId      int `json:"user_id"`
	CommodityId int `json:"commodity_id"`
	Quantity    int `json:"quantity"`
}

func NewCartHandler(cartService *service.CartService, userService *service.UserService) *CartHandler {
	return &CartHandler{cartService: cartService, userService: userService}
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	req := &reqAdd{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	err = h.cartService.AddToCart(req.UserId, req.CommodityId, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "add success"})
	log.Info("user", req.UserId, "add to cart success")
}
