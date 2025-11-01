package dto

import "server/internal/product/cart/model"

// CartResponse 购物车响应
type CartResponse struct {
	Items []*model.Cart `json:"items"`
}
