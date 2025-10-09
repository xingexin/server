package repository

import (
	"server/internal/product/model"

	"gorm.io/gorm"
)

type CommodityWriter interface {
	CreateCommodity(commodity *model.Commodity) error

	DeleteCommodity(id int) error
	UpdateCommodity(commodity *model.Commodity) error
}

type CommodityReader interface {
	FindCommodityById(id int) (*model.Commodity, error)
	ListCommodity() ([]*model.Commodity, error)
}

type CommodityRepository interface {
	CommodityWriter
	CommodityReader
} //商品操作

type gormCommodityRepository struct {
	gormDB *gorm.DB
}

func NewCommodityRepository(gDB *gorm.DB) CommodityRepository {
	return &gormCommodityRepository{gormDB: gDB}
}

func (cRepo *gormCommodityRepository) CreateCommodity(commodity *model.Commodity) error {
	err := cRepo.gormDB.Create(commodity).Error
	return err
}

func (cRepo *gormCommodityRepository) DeleteCommodity(id int) error {
	err := cRepo.gormDB.Delete(&model.User{}, id).Error
	return err
}

func (cRepo *gormCommodityRepository) UpdateCommodity(commodity *model.Commodity) error {
	err := cRepo.gormDB.Where("id=?", commodity.ID).Updates(commodity).Error
	return err
}

func (cRepo *gormCommodityRepository) FindCommodityById(id int) (*model.Commodity, error) {
	var commodity model.Commodity
	err := cRepo.gormDB.Where("id=?", id).Find(&commodity).Error
	return &commodity, err
}

func (cRepo *gormCommodityRepository) ListCommodity() ([]*model.Commodity, error) {
	cArr := make([]*model.Commodity, 0)
	if err := cRepo.gormDB.Find(&cArr).Error; err != nil {
		return nil, err
	}
	return cArr, nil
}
