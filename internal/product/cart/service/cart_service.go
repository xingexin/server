package service

import (
	"server/internal/product/cart/model"
	"server/internal/product/cart/repository"
	"time"
)

// CartService 提供购物车相关的业务逻辑服务
type CartService struct {
	cartRepo repository.CartRepository
}

// NewCartService 创建一个新的购物车服务实例
func NewCartService(repo repository.CartRepository) *CartService {
	return &CartService{cartRepo: repo}
}

// AddToCart 添加商品到购物车
func (cs *CartService) AddToCart(userId int, commodityId int, quantity int) error {

	cart := &model.Cart{
		UserId:      userId,
		CommodityId: commodityId,
		Quantity:    quantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return cs.cartRepo.CreateCart(cart)
}

// RemoveFromCart 从购物车中移除商品
func (cs *CartService) RemoveFromCart(cartId int) error {
	return cs.cartRepo.DeleteCart(cartId)
}

// UpdateCart 更新购物车中商品的数量
func (cs *CartService) UpdateCart(cartId int, quantity int) error {
	cart, err := cs.cartRepo.FindCartById(cartId)
	if err != nil {
		return err
	}
	cart.Quantity = quantity
	cart.UpdatedAt = time.Now()
	return cs.cartRepo.UpdateCart(cart)
}

// GetCart 获取用户的购物车
func (cs *CartService) GetCart(userId int) ([]*model.Cart, error) {
	return cs.cartRepo.FindCartByUserId(userId)
}
