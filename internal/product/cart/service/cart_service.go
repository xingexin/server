package service

import (
	"server/internal/product/cart/model"
	"server/internal/product/cart/repository"
	"time"
)

type CartService struct {
	cartRepo repository.CartRepository
}

func NewCartService(repo repository.CartRepository) *CartService {
	return &CartService{cartRepo: repo}
}

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

func (cs *CartService) RemoveFromCart(cartId int) error {
	return cs.cartRepo.DeleteCart(cartId)
}

func (cs *CartService) UpdateCart(cartId int, quantity int) error {
	cart, err := cs.cartRepo.FindCartById(cartId)
	if err != nil {
		return err
	}
	cart.Quantity = quantity
	cart.UpdatedAt = time.Now()
	return cs.cartRepo.UpdateCart(cart)
}

func (cs *CartService) GetCart(cartId int) (*model.Cart, error) {
	return cs.cartRepo.FindCartById(cartId)
}
