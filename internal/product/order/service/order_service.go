package service

import (
	"context"
	commodityRepository "server/internal/product/commodity/repository"
	"server/internal/product/order/model"
	"server/internal/product/order/repository"
	"time"
)

type OrderService struct {
	oRepo         repository.OrderRepository
	cRedisRepo    commodityRepository.StockCacheRepository
	commodityRepo commodityRepository.CommodityRepository
}

func NewOrderService(oRepo repository.OrderRepository, cRedisRepo commodityRepository.StockCacheRepository, commodityRepo commodityRepository.CommodityRepository) *OrderService {
	return &OrderService{oRepo: oRepo, cRedisRepo: cRedisRepo, commodityRepo: commodityRepo}
}

func (os *OrderService) CreateOrder(userId int, commodityId int, quantity int, totalPrice string, address string) error {
	order := &model.Order{UserId: userId, CommodityId: commodityId, Quantity: quantity, TotalPrice: totalPrice, Status: "pending", Address: address, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	code, err := os.cRedisRepo.DecreaseStock(context.TODO(), commodityId, quantity)
	switch code {
	case 1:
		return err
	case 2:
		commodity, err := os.commodityRepo.FindCommodityById(commodityId)
		if err != nil {
			return err
		}
		err = os.cRedisRepo.InitStockCache(context.TODO(), commodityId, commodity.Stock)
		if err != nil {
			return err
		}
	}
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
