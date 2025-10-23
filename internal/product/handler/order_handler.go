package handler

import (
	"server/internal/product/service"
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

type reqCreateOrder struct {
	UserId      int    `json:"user_id"`
	CommodityId int    `json:"commodity_id"`
	Quantity    int    `json:"quantity"`
	TotalPrice  string `json:"total_price"`
	Address     string `json:"address"`
}

type reqUpdateOrder struct {
	Id      int    `json:"id"`
	Status  string `json:"status"`
	Address string `json:"address"`
}

type reqDeleteOrder struct {
	Id int `json:"id"`
}

type reqGetOrder struct {
	Id int `json:"id"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	req := &reqCreateOrder{}
	err := c.ShouldBindJSON(req)
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
	req := reqUpdateOrder{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	if req.Status != "" {
		err = h.oSvc.UpdateOrderStatus(req.Id, req.Status)
		if err != nil {
			response.InternalServerError(c, response.CodeInternalError, "server busy")
			return
		}
	}
	if req.Address != "" {
		err = h.oSvc.UpdateOrderAddress(req.Id, req.Address)
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
	req := reqDeleteOrder{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid JSON")
		return
	}
	err = h.oSvc.DeleteOrder(req.Id)
	if err != nil {
		response.InternalServerError(c, response.CodeInternalError, "server busy")
		return
	}
	response.SuccessWithMessage(c, "delete success", nil)
	log.Info("order delete success:", req.Id)
	return
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	req := reqGetOrder{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	order, err := h.oSvc.GetOrderById(req.Id)
	if err != nil {
		response.InternalServerError(c, response.CodeInternalError, "server busy")
		return
	}
	response.Success(c, gin.H{"order": order})
	log.Info("order get success:", req.Id)
	return
}
