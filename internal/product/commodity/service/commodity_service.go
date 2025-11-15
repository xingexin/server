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
// 业务流程：
// 1. 设置创建时间为当前时间
// 2. 设置更新时间为当前时间
// 3. 设置商品状态为true（上架状态）
// 4. 设置库存初始值为0（需后续手动设置库存）
// 5. 调用Repository层创建商品记录
//
// 注意：
// - 商品创建后状态默认为true（上架）
// - 库存初始值强制为0，忽略传入的Stock值
// - 创建后需要手动更新库存或使用库存管理功能
// - 如需使用Redis库存缓存，需调用StockCacheService.InitStockCache初始化
func (c *CommodityService) CreateCommodity(commodity *model.Commodity) error {
	// 设置创建时间和更新时间为当前时间
	commodity.CreatedAt = time.Now()
	commodity.UpdateAt = time.Now()
	// 设置商品状态为true（上架状态，可供购买）
	commodity.Status = true
	// 强制设置库存为0，防止创建时直接设置库存导致Redis和MySQL不一致
	commodity.Stock = 0

	return c.cRepo.CreateCommodity(commodity)
}

// RemoveCommodity 根据ID删除商品
func (c *CommodityService) RemoveCommodity(id int) error {
	return c.cRepo.DeleteCommodity(id)
}

// UpdateCommodity 更新商品信息，保留原有的创建时间和状态，更新更新时间
// 业务流程：
// 1. 根据ID查询原有商品信息
// 2. 保留原商品的创建时间（CreatedAt）
// 3. 保留原商品的状态（Status）
// 4. 设置更新时间为当前时间
// 5. 调用Repository层更新商品记录
//
// 保留字段说明：
// - CreatedAt: 保留原有创建时间，不允许修改
// - Status: 保留原有状态，需通过专门的上下架接口修改
//
// 可更新字段：
// - Name: 商品名称
// - Price: 商品价格
// - Stock: 商品库存（注意：更新MySQL库存不会自动同步到Redis）
//
// 注意：
// - 如果使用了Redis库存缓存，更新Stock后需手动刷新Redis缓存
func (c *CommodityService) UpdateCommodity(commodity *model.Commodity) error {
	// 查询原有商品信息，用于保留不可修改的字段
	com, err := c.cRepo.FindCommodityById(commodity.ID)
	if err != nil {
		return err
	}

	// 保留原有的创建时间，不允许修改
	commodity.CreatedAt = com.CreatedAt
	// 保留原有的状态，状态变更需通过专门的接口
	commodity.Status = com.Status
	// 设置更新时间为当前时间
	commodity.UpdateAt = time.Now()

	return c.cRepo.UpdateCommodity(commodity)
}

// FindCommodityById 根据ID查找商品
func (c *CommodityService) FindCommodityById(id int) (*model.Commodity, error) {
	return c.cRepo.FindCommodityById(id)
}

// FindCommodityByName 根据名称查找商品
func (c *CommodityService) FindCommodityByName(name string) ([]*model.Commodity, error) {
	return c.cRepo.FindCommodityByName(name)

}

// ListCommodity 获取所有商品列表
func (c *CommodityService) ListCommodity() ([]*model.Commodity, error) {
	return c.cRepo.ListCommodity()
}
