package service

import (
	"server/internal/product/model"
	"server/internal/product/repository"
	"time"
)

type CommodityService struct {
	cRepo repository.CommodityRepository
}

func NewCommodityService(repository repository.CommodityRepository) *CommodityService {
	return &CommodityService{cRepo: repository}
}

func (c *CommodityService) CreateCommodity(commodity *model.Commodity) error {
	commodity.CreatedAt = time.Now()
	commodity.UpdateAt = time.Now()
	commodity.Status = true
	commodity.Stock = 0
	return c.cRepo.CreateCommodity(commodity)
}

func (c *CommodityService) RemoveCommodity(id int) error {
	return c.cRepo.DeleteCommodity(id)
}

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

func (c *CommodityService) FindCommodity(id int) (*model.Commodity, error) {
	return c.cRepo.FindCommodityById(id)
}

func (c *CommodityService) ListCommodity() ([]*model.Commodity, error) {
	return c.cRepo.ListCommodity()
}
