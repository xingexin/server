package repository

import (
	"fmt"
	"server/internal/product/cart/model"

	"gorm.io/gorm"
)

type cartWriter interface {
	CreateCart(cart *model.Cart) error
	DeleteCart(id int) error
	UpdateCart(cart *model.Cart) error
}

type cartReader interface {
	FindCartById(id int) (*model.Cart, error)
	FindCartByUserId(userId int) ([]*model.Cart, error)
	ListCart() ([]*model.Cart, error)
}

type CartRepository interface {
	cartWriter
	cartReader
}

type gormCartRepository struct {
	gormDB *gorm.DB
}

func NewCartRepository(gDB *gorm.DB) CartRepository {
	return &gormCartRepository{gormDB: gDB}
}

func (cRepo *gormCartRepository) CreateCart(cart *model.Cart) error {
	err := cRepo.gormDB.Create(cart).Error
	return err
}

func (cRepo *gormCartRepository) DeleteCart(id int) error {
	err := cRepo.gormDB.Delete(&model.Cart{}, id)
	if err.RowsAffected == 0 {
		return fmt.Errorf("cart with id=%d not found", id)
	}
	return err.Error
}

func (cRepo *gormCartRepository) UpdateCart(cart *model.Cart) error {
	err := cRepo.gormDB.Where("id=?", cart.Id).Updates(cart).Error
	return err
}

func (cRepo *gormCartRepository) FindCartById(id int) (*model.Cart, error) {
	var cart model.Cart
	err := cRepo.gormDB.First(&cart, id).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (cRepo *gormCartRepository) FindCartByUserId(userId int) ([]*model.Cart, error) {
	var carts []*model.Cart
	err := cRepo.gormDB.Where("user_id = ?", userId).Find(&carts).Error
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (cRepo *gormCartRepository) ListCart() ([]*model.Cart, error) {
	var carts []*model.Cart
	err := cRepo.gormDB.Find(&carts).Error
	if err != nil {
		return nil, err
	}
	return carts, nil
}
