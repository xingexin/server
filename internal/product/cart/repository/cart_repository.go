package repository

import (
	"fmt"
	"server/internal/product/cart/model"

	"gorm.io/gorm"
)

// cartWriter 定义购物车写操作接口
type cartWriter interface {
	CreateCart(cart *model.Cart) error
	DeleteCart(id int) error
	UpdateCart(cart *model.Cart) error
}

// cartReader 定义购物车读操作接口
type cartReader interface {
	FindCartById(id int) (*model.Cart, error)
	FindCartByUserId(userId int) ([]*model.Cart, error)
	ListCart() ([]*model.Cart, error)
}

// CartRepository 购物车操作的数据访问接口，组合了读写操作
type CartRepository interface {
	cartWriter
	cartReader
}

type gormCartRepository struct {
	gormDB *gorm.DB
}

// NewCartRepository 创建一个新的购物车仓储实例
func NewCartRepository(gDB *gorm.DB) CartRepository {
	return &gormCartRepository{gormDB: gDB}
}

// CreateCart 在数据库中创建新购物车条目
func (cRepo *gormCartRepository) CreateCart(cart *model.Cart) error {
	err := cRepo.gormDB.Create(cart).Error
	return err
}

// DeleteCart 根据ID从数据库中删除购物车条目
func (cRepo *gormCartRepository) DeleteCart(id int) error {
	err := cRepo.gormDB.Delete(&model.Cart{}, id)
	if err.RowsAffected == 0 {
		return fmt.Errorf("cart with id=%d not found", id)
	}
	return err.Error
}

// UpdateCart 更新数据库中的购物车信息
func (cRepo *gormCartRepository) UpdateCart(cart *model.Cart) error {
	err := cRepo.gormDB.Where("id=?", cart.Id).Updates(cart).Error
	return err
}

// FindCartById 根据ID从数据库中查找购物车条目
func (cRepo *gormCartRepository) FindCartById(id int) (*model.Cart, error) {
	var cart model.Cart
	err := cRepo.gormDB.First(&cart, id).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// FindCartByUserId 根据用户ID从数据库中查找该用户的所有购物车条目
func (cRepo *gormCartRepository) FindCartByUserId(userId int) ([]*model.Cart, error) {
	var carts []*model.Cart
	err := cRepo.gormDB.Where("user_id = ?", userId).Find(&carts).Error
	if err != nil {
		return nil, err
	}
	return carts, nil
}

// ListCart 从数据库中获取所有购物车条目
func (cRepo *gormCartRepository) ListCart() ([]*model.Cart, error) {
	var carts []*model.Cart
	err := cRepo.gormDB.Find(&carts).Error
	if err != nil {
		return nil, err
	}
	return carts, nil
}
