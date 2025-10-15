package handler

import (
	"server/internal/product/service"
	"server/pkg/response"

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

type reqUpdate struct {
	CartId   int `json:"cart_id"`
	Quantity int `json:"quantity"`
}

func NewCartHandler(cartService *service.CartService, userService *service.UserService) *CartHandler {
	return &CartHandler{cartService: cartService, userService: userService}
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	req := &reqAdd{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	err = h.cartService.AddToCart(req.UserId, req.CommodityId, req.Quantity)
	if err != nil {
		response.BadRequest(c, response.CodeInternalError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "add success", nil)
	log.Info("user", req.UserId, "add to cart success")
}

func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	cardId := 0
	err := c.ShouldBindUri(&cardId)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid URI")
		return
	}
	err = h.cartService.RemoveFromCart(cardId)
	if err != nil {
		response.BadRequest(c, response.CodeInternalError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "remove success", nil)
	log.Info("cart", cardId, "remove from cart success")
	return
}

func (h *CartHandler) UpdateCart(c *gin.Context) {
	req := &reqUpdate{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	err = h.cartService.UpdateCart(req.CartId, req.Quantity)
	if err != nil {
		response.BadRequest(c, response.CodeInternalError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "update success", nil)
	log.Info("cart", req.CartId, "update cart success")
	return
}

func (h *CartHandler) GetCart(c *gin.Context) {
	cartId := 0
	err := c.ShouldBindJSON(&cartId)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid JSON")
		return
	}
	cart, err := h.cartService.GetCart(cartId)
	if err != nil {
		response.BadRequest(c, response.CodeInternalError, err.Error())
		return
	}
	response.Success(c, gin.H{"cart": cart})
	log.Info("cart", cartId, "get cart success")
	return
}
