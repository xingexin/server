package dto

import "server/internal/product/order/model"

// OrderResponse 订单响应
type OrderResponse struct {
	Order *model.Order `json:"order"`
}
