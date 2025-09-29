package service

import "server/internal/product/repository"

type CommodityService struct {
	cRepo repository.CommodityRepository
}

func NewCommodityService(repository repository.CommodityRepository) *CommodityService {
	return &CommodityService{cRepo: repository}
}

func (c *CommodityService) AddProduct()
