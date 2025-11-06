package service

import (
	"context"
	commodityRepository "server/internal/product/commodity/repository"
	"server/internal/product/order/model"
	"server/internal/product/order/repository"
	"time"
)

// OrderService 提供订单相关的业务逻辑服务
type OrderService struct {
	oRepo              repository.OrderRepository
	cRedisRepo         commodityRepository.StockCacheRepository
	commodityRepo      commodityRepository.CommodityRepository
	orderCancelService OrderCancelService
}

// NewOrderService 创建一个新的订单服务实例
func NewOrderService(oRepo repository.OrderRepository, cRedisRepo commodityRepository.StockCacheRepository, commodityRepo commodityRepository.CommodityRepository, orderCancelService OrderCancelService) *OrderService {
	return &OrderService{
		oRepo:              oRepo,
		cRedisRepo:         cRedisRepo,
		commodityRepo:      commodityRepo,
		orderCancelService: orderCancelService,
	}
}

// CreateOrder 创建订单，先扣减Redis库存，成功后创建订单并加入延迟取消队列
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
	err = os.oRepo.CreateOrder(order)
	if err != nil {
		return err
	}
	err = os.orderCancelService.createOrderTask(order.Id, commodityId, quantity)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrderStatus 更新订单状态
func (os *OrderService) UpdateOrderStatus(id int, status string) error {
	order := &model.Order{Status: status, Id: id}
	return os.oRepo.UpdateOrder(order)
}

// UpdateOrderAddress 更新订单地址
func (os *OrderService) UpdateOrderAddress(id int, address string) error {
	order := &model.Order{Address: address, Id: id}
	return os.oRepo.UpdateOrder(order)
}

// DeleteOrder 删除订单
func (os *OrderService) DeleteOrder(orderId int) error {
	return os.oRepo.DeleteOrder(orderId)
}

// GetOrderById 根据订单ID获取订单
func (os *OrderService) GetOrderById(orderId int) (*model.Order, error) {
	return os.oRepo.FindOrderById(orderId)
}

// GetOrdersByUserId 根据用户ID获取该用户的所有订单
func (os *OrderService) GetOrdersByUserId(userId int) ([]*model.Order, error) {
	return os.oRepo.FindOrdersByUserId(userId)
}
