package handler

import (
	"server/internal/product/cart/dto"
	"server/internal/product/cart/service"
	userService "server/internal/product/user/service"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CartHandler struct {
	cartService *service.CartService
	userService *userService.UserService
}

func NewCartHandler(cartService *service.CartService, uService *userService.UserService) *CartHandler {
	return &CartHandler{cartService: cartService, userService: uService}
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	var req dto.AddToCartRequest
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
	var req dto.RemoveFromCartRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid URI")
		return
	}
	err = h.cartService.RemoveFromCart(req.ID)
	if err != nil {
		response.BadRequest(c, response.CodeInternalError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "remove success", nil)
	log.Info("cart", req.ID, "remove from cart success")
	return
}

func (h *CartHandler) UpdateCart(c *gin.Context) {
	var req dto.UpdateCartRequest
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
	var req dto.GetCartRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid JSON")
		return
	}
	cart, err := h.cartService.GetCart(req.ID)
	if err != nil {
		response.BadRequest(c, response.CodeInternalError, err.Error())
		return
	}
	res := dto.CartResponse{Cart: cart}
	response.Success(c, res)
	log.Info("cart", req.ID, "get cart success")
	return
}
