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
	return c.cRepo.UpdateCommodity(commodity)
}

func (c *CommodityService) FindCommodityByName(name string) (*model.Commodity, error) {
	return c.cRepo.FindCommodityByName(name)

}

func (c *CommodityService) ListCommodity() ([]*model.Commodity, error) {
	return c.cRepo.ListCommodity()
}
