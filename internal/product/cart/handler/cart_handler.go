package handler

import (
	"server/internal/product/cart/dto"
	"server/internal/product/cart/service"
	userService "server/internal/product/user/service"
	"server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// CartHandler 处理购物车相关的HTTP请求
type CartHandler struct {
	cartService *service.CartService
	userService *userService.UserService
}

// NewCartHandler 创建一个新的购物车处理器实例
func NewCartHandler(cartService *service.CartService, uService *userService.UserService) *CartHandler {
	return &CartHandler{cartService: cartService, userService: uService}
}

// AddToCart 处理添加商品到购物车请求
func (h *CartHandler) AddToCart(c *gin.Context) {
	// 从JWT获取认证后的userID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, response.CodeUnauthorized, "user not authenticated")
		return
	}
	uid := userID.(int)

	var req dto.AddToCartRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	err = h.cartService.AddToCart(uid, req.CommodityId, req.Quantity)
	if err != nil {
		response.BadRequest(c, response.CodeInternalError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "add success", nil)
	log.Info("user", uid, "add to cart success")
}

// RemoveFromCart 处理从购物车移除商品请求
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid id parameter")
		return
	}
	err = h.cartService.RemoveFromCart(id)
	if err != nil {
		response.BadRequest(c, response.CodeInternalError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "remove success", nil)
	log.Info("cart", id, "remove from cart success")
	return
}

// UpdateCart 处理更新购物车商品数量请求
func (h *CartHandler) UpdateCart(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid id parameter")
		return
	}

	var req dto.UpdateCartRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	err = h.cartService.UpdateCart(id, req.Quantity)
	if err != nil {
		response.BadRequest(c, response.CodeInternalError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "update success", nil)
	log.Info("cart", id, "update cart success")
	return
}

// GetCart 处理获取购物车请求
func (h *CartHandler) GetCart(c *gin.Context) {
	// 从JWT获取认证后的userID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, response.CodeUnauthorized, "user not authenticated")
		return
	}
	uid := userID.(int)

	items, err := h.cartService.GetCart(uid)
	if err != nil {
		response.BadRequest(c, response.CodeInternalError, err.Error())
		return
	}
	res := dto.CartResponse{Items: items}
	response.Success(c, res)
	log.Info("user", uid, "get cart success")
	return
}
