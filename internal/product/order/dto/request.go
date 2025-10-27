package dto

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	UserId      int    `json:"user_id" binding:"required"`
	CommodityId int    `json:"commodity_id" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
	TotalPrice  string `json:"total_price" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

// UpdateOrderRequest 更新订单请求
type UpdateOrderRequest struct {
	ID      int    `json:"id" binding:"required"`
	Status  string `json:"status"`
	Address string `json:"address"`
}

// DeleteOrderRequest 删除订单请求
type DeleteOrderRequest struct {
	ID int `json:"id" binding:"required"`
}

// GetOrderRequest 获取订单请求
type GetOrderRequest struct {
	ID int `json:"id" binding:"required"`
}
