package service

import (
	"server/internal/product/commodity/model"
	"server/internal/product/commodity/repository"
	"time"
)

// CommodityService 提供商品相关的业务逻辑服务
type CommodityService struct {
	cRepo repository.CommodityRepository
}

// NewCommodityService 创建一个新的商品服务实例
func NewCommodityService(repository repository.CommodityRepository) *CommodityService {
	return &CommodityService{cRepo: repository}
}

// CreateCommodity 创建新商品，设置创建时间、更新时间、状态和库存初始值
func (c *CommodityService) CreateCommodity(commodity *model.Commodity) error {
	commodity.CreatedAt = time.Now()
	commodity.UpdateAt = time.Now()
	commodity.Status = true
	commodity.Stock = 0
	return c.cRepo.CreateCommodity(commodity)
}

// RemoveCommodity 根据ID删除商品
func (c *CommodityService) RemoveCommodity(id int) error {
	return c.cRepo.DeleteCommodity(id)
}

// UpdateCommodity 更新商品信息，保留原有的创建时间和状态，更新更新时间
func (c *CommodityService) UpdateCommodity(commodity *model.Commodity) error {
	com, err := c.cRepo.FindCommodityById(commodity.ID)
	if err != nil {
		return err
	}
	commodity.CreatedAt = com.CreatedAt
	commodity.Status = com.Status
	commodity.UpdateAt = time.Now()
	return c.cRepo.UpdateCommodity(commodity)
}

// FindCommodityById 根据ID查找商品
func (c *CommodityService) FindCommodityById(id int) (*model.Commodity, error) {
	return c.cRepo.FindCommodityById(id)
}

func (c *CommodityService) FindCommodityByName(name string) ([]*model.Commodity, error) {
	return c.cRepo.FindCommodityByName(name)

}

// ListCommodity 获取所有商品列表
func (c *CommodityService) ListCommodity() ([]*model.Commodity, error) {
	return c.cRepo.ListCommodity()
}
