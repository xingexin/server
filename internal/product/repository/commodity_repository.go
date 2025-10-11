package repository

import (
	"fmt"
	"server/internal/product/model"

	"gorm.io/gorm"
)

// CommodityWriter 定义商品写操作接口
type CommodityWriter interface {
	CreateCommodity(commodity *model.Commodity) error

	DeleteCommodity(id int) error
	UpdateCommodity(commodity *model.Commodity) error
}

// CommodityReader 定义商品读操作接口
type CommodityReader interface {
	FindCommodityById(id int) (*model.Commodity, error)
	ListCommodity() ([]*model.Commodity, error)
	FindCommodityByName(name string) ([]*model.Commodity, error)
}

// CommodityRepository 商品操作的数据访问接口，组合了读写操作
type CommodityRepository interface {
	CommodityWriter
	CommodityReader
} //商品操作

type gormCommodityRepository struct {
	gormDB *gorm.DB
}

// NewCommodityRepository 创建一个新的商品仓储实例
func NewCommodityRepository(gDB *gorm.DB) CommodityRepository {
	return &gormCommodityRepository{gormDB: gDB}
}

// CreateCommodity 在数据库中创建新商品记录
func (cRepo *gormCommodityRepository) CreateCommodity(commodity *model.Commodity) error {
	err := cRepo.gormDB.Create(commodity).Error
	return err
}

// DeleteCommodity 根据ID从数据库中删除商品记录
func (cRepo *gormCommodityRepository) DeleteCommodity(id int) error {
	err := cRepo.gormDB.Delete(&model.Commodity{}, id)
	if err.RowsAffected == 0 {
		return fmt.Errorf("commodity with id=%d not found", id)
	}
	return err.Error
}

// UpdateCommodity 更新数据库中的商品信息
func (cRepo *gormCommodityRepository) UpdateCommodity(commodity *model.Commodity) error {
	err := cRepo.gormDB.Where("id=?", commodity.ID).Updates(commodity).Error
	return err
}

// FindCommodityById 根据ID从数据库中查找商品
func (cRepo *gormCommodityRepository) FindCommodityById(id int) (*model.Commodity, error) {
	var commodity model.Commodity
	err := cRepo.gormDB.Where("id=?", id).Find(&commodity).Error
	return &commodity, err
}

func (cRepo *gormCommodityRepository) FindCommodityByName(name string) ([]*model.Commodity, error) {
	var commodities []*model.Commodity
	err := cRepo.gormDB.Where("name=?", name).Find(&commodities).Error
	if err != nil {
		return nil, err
	}
	return commodities, nil
}

// ListCommodity 从数据库中获取所有商品列表
func (cRepo *gormCommodityRepository) ListCommodity() ([]*model.Commodity, error) {
	cArr := make([]*model.Commodity, 0)
	if err := cRepo.gormDB.Find(&cArr).Error; err != nil {
		return nil, err
	}
	return cArr, nil
}
