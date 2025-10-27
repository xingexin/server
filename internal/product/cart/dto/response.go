package dto

import "server/internal/product/cart/model"

// CartResponse 购物车响应
type CartResponse struct {
	Cart *model.Cart `json:"cart"`
}
