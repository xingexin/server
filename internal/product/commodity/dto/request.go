package dto

// CreateCommodityRequest 创建商品请求
type CreateCommodityRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
	Stock int     `json:"stock" binding:"required"`
}

// UpdateCommodityRequest 更新商品请求
type UpdateCommodityRequest struct {
	ID    int     `json:"id" binding:"required"`
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
	Stock int     `json:"stock" binding:"required"`
}

// DeleteCommodityRequest 删除商品请求
type DeleteCommodityRequest struct {
	ID int `json:"id" binding:"required"`
}

// FindCommodityByNameRequest 根据名称查找商品请求
type FindCommodityByNameRequest struct {
	Name string `json:"name" binding:"required"`
}
