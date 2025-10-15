package repository

import (
	"server/internal/product/model"

	"gorm.io/gorm"
)

type orderWriter interface {
	CreateOrder(order *model.Order) error
	UpdateOrder(order *model.Order) error
	DeleteOrder(orderId int) error
}

type orderReader interface {
	FindOrderById(orderId int) (*model.Order, error)
	FindOrdersByUserId(userId int) ([]*model.Order, error)
}

type OrderRepository interface {
	orderWriter
	orderReader
}

type gormOrderRepository struct {
	gormDB *gorm.DB
}

func NewOrderRepository(gDB *gorm.DB) OrderRepository {
	return &gormOrderRepository{gormDB: gDB}
}

func (oRepo *gormOrderRepository) CreateOrder(order *model.Order) error {
	return oRepo.gormDB.Create(order).Error
}

func (oRepo *gormOrderRepository) UpdateOrder(order *model.Order) error {
	return oRepo.gormDB.Where("id=?", order.Id).Updates(order).Error
}

func (oRepo *gormOrderRepository) DeleteOrder(orderId int) error {
	err := oRepo.gormDB.Delete(&model.Order{}, orderId)
	if err.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return err.Error
}

func (oRepo *gormOrderRepository) FindOrderById(orderId int) (*model.Order, error) {
	var order model.Order
	err := oRepo.gormDB.First(&order, orderId).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (oRepo *gormOrderRepository) FindOrdersByUserId(userId int) ([]*model.Order, error) {
	var orders []*model.Order
	err := oRepo.gormDB.Where("user_id = ?", userId).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
