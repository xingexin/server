package dto

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	UserId      int `json:"user_id" binding:"required"`
	CommodityId int `json:"commodity_id" binding:"required"`
	Quantity    int `json:"quantity" binding:"required"`
}

// RemoveFromCartRequest 从购物车移除请求
type RemoveFromCartRequest struct {
	ID int `json:"id" binding:"required"`
}

// UpdateCartRequest 更新购物车请求
type UpdateCartRequest struct {
	CartId   int `json:"cart_id" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
}

// GetCartRequest 获取购物车请求
type GetCartRequest struct {
	ID int `json:"id" binding:"required"`
}
