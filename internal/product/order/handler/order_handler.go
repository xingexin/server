package handler

import (
	"server/internal/product/order/dto"
	"server/internal/product/order/service"
	"server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// OrderHandler 处理订单相关的HTTP请求
type OrderHandler struct {
	oSvc *service.OrderService
}

// NewOrderHandler 创建一个新的订单处理器实例
func NewOrderHandler(oSvc *service.OrderService) *OrderHandler {
	return &OrderHandler{oSvc: oSvc}
}

// CreateOrder 处理创建订单请求
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// 从JWT获取认证后的userID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, response.CodeUnauthorized, "user not authenticated")
		return
	}
	uid := userID.(int)

	var req dto.CreateOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	err = h.oSvc.CreateOrder(uid, req.CommodityId, req.Quantity, req.TotalPrice, req.Address)
	if err != nil {
		response.InternalServerError(c, response.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
	log.Info("order create success, userID:", uid, "commodityID:", req.CommodityId)
}

// UpdateOrderStatus 处理更新订单请求（状态或地址）
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid id parameter")
		return
	}

	var req dto.UpdateOrderRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	if req.Status != "" {
		err = h.oSvc.UpdateOrderStatus(id, req.Status)
		if err != nil {
			response.InternalServerError(c, response.CodeInternalError, "server busy")
			return
		}
	}
	if req.Address != "" {
		err = h.oSvc.UpdateOrderAddress(id, req.Address)
		if err != nil {
			response.InternalServerError(c, response.CodeInternalError, "server busy")
			return
		}
	}
	response.Success(c, nil)
	log.Info("order update success:", req.Status, req.Address)
	return
}

// DeleteOrder 处理删除订单请求
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid id parameter")
		return
	}
	err = h.oSvc.DeleteOrder(id)
	if err != nil {
		response.InternalServerError(c, response.CodeInternalError, "server busy")
		return
	}
	response.SuccessWithMessage(c, "delete success", nil)
	log.Info("order delete success:", id)
	return
}

// GetOrder 处理获取订单详情请求
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid id parameter")
		return
	}
	order, err := h.oSvc.GetOrderById(id)
	if err != nil {
		response.InternalServerError(c, response.CodeInternalError, "server busy")
		return
	}
	res := dto.OrderResponse{Order: order}
	response.Success(c, res)
	log.Info("order get success:", id)
	return
}
