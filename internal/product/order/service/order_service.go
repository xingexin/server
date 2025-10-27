package service

import (
	"server/internal/product/order/model"
	"server/internal/product/order/repository"
	"time"
)

type OrderService struct {
	oRepo repository.OrderRepository
}

func NewOrderService(oRepo repository.OrderRepository) *OrderService {
	return &OrderService{oRepo: oRepo}
}

func (os *OrderService) CreateOrder(userId int, commodityId int, quantity int, totalPrice string, address string) error {
	order := &model.Order{UserId: userId, CommodityId: commodityId, Quantity: quantity, TotalPrice: totalPrice, Status: "pending", Address: address, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	return os.oRepo.CreateOrder(order)
}

func (os *OrderService) UpdateOrderStatus(id int, status string) error {
	order := &model.Order{Status: status, Id: id}
	return os.oRepo.UpdateOrder(order)
}

func (os *OrderService) UpdateOrderAddress(id int, address string) error {
	order := &model.Order{Address: address, Id: id}
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
