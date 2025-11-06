package repository

import (
	"server/internal/product/order/model"

	"gorm.io/gorm"
)

// orderWriter 定义订单写操作接口
type orderWriter interface {
	CreateOrder(order *model.Order) error
	UpdateOrder(order *model.Order) error
	DeleteOrder(orderId int) error
}

// orderReader 定义订单读操作接口
type orderReader interface {
	FindOrderById(orderId int) (*model.Order, error)
	FindOrdersByUserId(userId int) ([]*model.Order, error)
}

// OrderRepository 订单操作的数据访问接口，组合了读写操作
type OrderRepository interface {
	orderWriter
	orderReader
}

type gormOrderRepository struct {
	gormDB *gorm.DB
}

// NewOrderRepository 创建一个新的订单仓储实例
func NewOrderRepository(gDB *gorm.DB) OrderRepository {
	return &gormOrderRepository{gormDB: gDB}
}

// CreateOrder 在数据库中创建新订单记录
func (oRepo *gormOrderRepository) CreateOrder(order *model.Order) error {
	return oRepo.gormDB.Create(order).Error
}

// UpdateOrder 更新数据库中的订单信息
func (oRepo *gormOrderRepository) UpdateOrder(order *model.Order) error {
	return oRepo.gormDB.Where("id=?", order.Id).Updates(order).Error
}

// DeleteOrder 根据ID从数据库中删除订单记录
func (oRepo *gormOrderRepository) DeleteOrder(orderId int) error {
	err := oRepo.gormDB.Delete(&model.Order{}, orderId)
	if err.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return err.Error
}

// FindOrderById 根据ID从数据库中查找订单
func (oRepo *gormOrderRepository) FindOrderById(orderId int) (*model.Order, error) {
	var order model.Order
	err := oRepo.gormDB.First(&order, orderId).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindOrdersByUserId 根据用户ID从数据库中查找该用户的所有订单
func (oRepo *gormOrderRepository) FindOrdersByUserId(userId int) ([]*model.Order, error) {
	var orders []*model.Order
	err := oRepo.gormDB.Where("user_id = ?", userId).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
