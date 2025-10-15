package service

import (
	"server/internal/product/model"
	"server/internal/product/repository"
	"time"
)

type OrderService struct {
	oRepo repository.OrderRepository
}

func NewOrderService(oRepo repository.OrderRepository) *OrderService {
	return &OrderService{oRepo: oRepo}
}

func (os *OrderService) CreateOrder(userId int, commodityId int, quantity int, totalPrice float64, address string) error {
	order := &model.Order{UserId: userId, CommodityId: commodityId, Quantity: quantity, TotalPrice: totalPrice, Status: "pending", Address: address, CreatedAt: time.Now(), UpdateAt: time.Now()}
	return os.oRepo.CreateOrder(order)
}

func (os *OrderService) UpdateOrderStatus(status string) error {
	order := &model.Order{Status: status}
	return os.oRepo.UpdateOrder(order)
}

func (os *OrderService) UpdateOrderAddress(address string) error {
	order := &model.Order{Address: address}
	return os.oRepo.UpdateOrder(order)
}

func (os *OrderService) DeleteOrder(orderId int) error {
	return os.oRepo.DeleteOrder(orderId)
}

func (os *OrderService) GetOrderById(orderId int) (*model.Order, error) {
	return os.oRepo.FindOrderById(orderId)
}

func (os *OrderService) GetOrdersByUserId(userId int) ([]*model.Order, error) {
	return os.oRepo.FindOrdersByUserId(userId)
}
