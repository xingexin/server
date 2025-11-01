package dto

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	CommodityId int `json:"commodity_id" binding:"required"`
	Quantity    int `json:"quantity" binding:"required"`
}

// UpdateCartRequest 更新购物车请求
type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}
