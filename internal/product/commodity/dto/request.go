package dto

// CreateCommodityRequest 创建商品请求
type CreateCommodityRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
	Stock int     `json:"stock" binding:"required"`
}

// UpdateCommodityRequest 更新商品请求
type UpdateCommodityRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
	Stock int     `json:"stock" binding:"required"`
}
