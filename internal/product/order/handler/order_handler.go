package handler

import (
	"server/internal/product/order/dto"
	"server/internal/product/order/service"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type OrderHandler struct {
	oSvc *service.OrderService
}

func NewOrderHandler(oSvc *service.OrderService) *OrderHandler {
	return &OrderHandler{oSvc: oSvc}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	err = h.oSvc.CreateOrder(req.UserId, req.CommodityId, req.Quantity, req.TotalPrice, req.Address)
	if err != nil {
		response.InternalServerError(c, response.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
	log.Info("order create success:", req.UserId, req.CommodityId)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	var req dto.UpdateOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	if req.Status != "" {
		err = h.oSvc.UpdateOrderStatus(req.ID, req.Status)
		if err != nil {
			response.InternalServerError(c, response.CodeInternalError, "server busy")
			return
		}
	}
	if req.Address != "" {
		err = h.oSvc.UpdateOrderAddress(req.ID, req.Address)
		if err != nil {
			response.InternalServerError(c, response.CodeInternalError, "server busy")
			return
		}
	}
	response.Success(c, nil)
	log.Info("order update success:", req.Status, req.Address)
	return
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	var req dto.DeleteOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid JSON")
		return
	}
	err = h.oSvc.DeleteOrder(req.ID)
	if err != nil {
		response.InternalServerError(c, response.CodeInternalError, "server busy")
		return
	}
	response.SuccessWithMessage(c, "delete success", nil)
	log.Info("order delete success:", req.ID)
	return
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	var req dto.GetOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	order, err := h.oSvc.GetOrderById(req.ID)
	if err != nil {
		response.InternalServerError(c, response.CodeInternalError, "server busy")
		return
	}
	res := dto.OrderResponse{Order: order}
	response.Success(c, res)
	log.Info("order get success:", req.ID)
	return
}
