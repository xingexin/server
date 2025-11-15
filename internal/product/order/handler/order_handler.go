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
// 业务流程：
// 1. 从JWT中间件获取已认证的用户ID
// 2. 解析请求体中的订单信息（商品ID、数量、总价、地址）
// 3. 调用Service层创建订单（包含：扣减Redis库存、创建订单记录、加入延迟取消队列）
// 4. 返回创建结果
//
// 注意：
// - 订单创建后会自动加入15分钟延迟队列，超时未支付将自动取消并归还库存
// - 库存扣减失败（库存不足或Redis未初始化）会直接返回错误，不会创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// 从JWT中间件注入的上下文中获取认证后的userID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, response.CodeUnauthorized, "user not authenticated")
		return
	}
	uid := userID.(int)

	// 解析请求体，绑定到CreateOrderRequest结构体
	var req dto.CreateOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}

	// 调用Service层创建订单
	// 内部流程：扣减Redis库存 -> 创建订单记录 -> 加入延迟取消队列
	err = h.oSvc.CreateOrder(uid, req.CommodityId, req.Quantity, req.TotalPrice, req.Address)
	if err != nil {
		response.InternalServerError(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, nil)
	log.Info("order create success, userID:", uid, "commodityID:", req.CommodityId)
}

// UpdateOrderStatus 处理更新订单请求（状态或地址）
// 业务流程：
// 1. 从URL路径中提取订单ID
// 2. 解析请求体，支持更新状态和地址
// 3. 根据请求内容选择性更新（状态和地址可以单独或同时更新）
//
// 支持的订单状态：
// - pending: 待支付
// - paid: 已支付
// - cancelled: 已取消
// - completed: 已完成
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	// 从URL路径参数中获取订单ID（如：/order/123）
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid id parameter")
		return
	}

	// 解析请求体，获取要更新的字段
	var req dto.UpdateOrderRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}

	// 如果提供了状态字段，则更新订单状态
	if req.Status != "" {
		err = h.oSvc.UpdateOrderStatus(id, req.Status)
		if err != nil {
			response.InternalServerError(c, response.CodeInternalError, "server busy")
			return
		}
	}

	// 如果提供了地址字段，则更新收货地址
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
// 业务流程：
// 1. 从URL路径中提取订单ID
// 2. 调用Service层删除订单
//
// 注意：
// - 此方法仅删除订单记录，不会自动归还库存
// - 建议在删除前先确认订单状态，避免误删已支付订单
// - 如需取消订单并归还库存，应使用订单取消功能而非直接删除
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	// 从URL路径参数中获取订单ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid id parameter")
		return
	}

	// 调用Service层删除订单
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
// 业务流程：
// 1. 从URL路径中提取订单ID
// 2. 调用Service层查询订单详情
// 3. 将订单模型转换为响应DTO并返回
//
// 返回内容包括：
// - 订单ID、用户ID、商品ID
// - 商品数量、总价、收货地址
// - 订单状态、创建时间、更新时间
func (h *OrderHandler) GetOrder(c *gin.Context) {
	// 从URL路径参数中获取订单ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid id parameter")
		return
	}

	// 调用Service层查询订单详情
	order, err := h.oSvc.GetOrderById(id)
	if err != nil {
		response.InternalServerError(c, response.CodeInternalError, "server busy")
		return
	}

	// 封装响应数据
	res := dto.OrderResponse{Order: order}
	response.Success(c, res)
	log.Info("order get success:", id)
	return
}
