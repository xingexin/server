package dto

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	CommodityId int    `json:"commodity_id" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
	TotalPrice  string `json:"total_price" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

// UpdateOrderRequest 更新订单请求
type UpdateOrderRequest struct {
	Status  string `json:"status"`
	Address string `json:"address"`
}
