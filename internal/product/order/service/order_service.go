package service

import (
	"context"
	commodityRepository "server/internal/product/commodity/repository"
	"server/internal/product/order/model"
	"server/internal/product/order/repository"
	"time"
)

// OrderService 提供订单相关的业务逻辑服务
type OrderService struct {
	oRepo              repository.OrderRepository
	cRedisRepo         commodityRepository.StockCacheRepository
	commodityRepo      commodityRepository.CommodityRepository
	orderCancelService OrderCancelService
}

// NewOrderService 创建一个新的订单服务实例
func NewOrderService(oRepo repository.OrderRepository, cRedisRepo commodityRepository.StockCacheRepository, commodityRepo commodityRepository.CommodityRepository, orderCancelService OrderCancelService) *OrderService {
	return &OrderService{
		oRepo:              oRepo,
		cRedisRepo:         cRedisRepo,
		commodityRepo:      commodityRepo,
		orderCancelService: orderCancelService,
	}
}

// CreateOrder 创建订单，先扣减Redis库存，成功后创建订单并加入延迟取消队列
// 业务流程：
// 1. 构建订单对象（状态为pending，设置创建时间和更新时间）
// 2. 尝试从Redis扣减库存（使用Lua脚本保证原子性）
// 3. 如果Redis缓存未初始化（code=2），从MySQL查询库存并初始化Redis缓存
// 4. 扣减成功后创建订单记录
// 5. 将订单加入延迟取消队列（15分钟后自动取消）
//
// 库存扣减返回码说明：
// - code=0: 扣减成功
// - code=1: 扣减失败（网络错误等）
// - code=2: Redis缓存未初始化
// - code=3: 库存不足
//
// ⚠️ 已知问题（需修复）：
// 1. 当code=2时，只初始化了缓存但没有重新扣减库存
// 2. 当code=3时（库存不足），没有return，会继续创建订单
// 3. 这两个问题会导致：超卖、库存数据不一致
//
// 正确的处理逻辑应该是：
// - code=2: 初始化缓存后重新调用DecreaseStock
// - code=3: 直接return错误，不创建订单
func (os *OrderService) CreateOrder(userId int, commodityId int, quantity int, totalPrice string, address string) error {
	// 构建订单对象，初始状态为pending（待支付）
	order := &model.Order{
		UserId:      userId,
		CommodityId: commodityId,
		Quantity:    quantity,
		TotalPrice:  totalPrice,
		Status:      "pending",  // 订单状态：pending待支付、paid已支付、cancelled已取消、completed已完成
		Address:     address,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 尝试从Redis扣减库存（使用Lua脚本保证原子性，防止超卖）
	code, err := os.cRedisRepo.DecreaseStock(context.TODO(), commodityId, quantity)
	switch code {
	case 1: // 扣减失败（网络错误等）
		return err
	case 2: // Redis缓存未初始化
		// 从MySQL查询商品库存
		commodity, err := os.commodityRepo.FindCommodityById(commodityId)
		if err != nil {
			return err
		}
		// 初始化Redis库存缓存
		err = os.cRedisRepo.InitStockCache(context.TODO(), commodityId, commodity.Stock)
		if err != nil {
			return err
		}
		// ⚠️ BUG: 这里应该重新调用DecreaseStock扣减库存，否则库存未扣减就创建了订单
	case 3: // 库存不足
		// ⚠️ BUG: 这里应该return错误，否则库存不足也会创建订单
	}

	// 创建订单记录
	err = os.oRepo.CreateOrder(order)
	if err != nil {
		return err
	}

	// 将订单加入延迟取消队列（15分钟后如果还是pending状态，会自动取消并归还库存）
	err = os.orderCancelService.createOrderTask(order.Id, commodityId, quantity)
	if err != nil {
		return err
	}

	return nil
}

// UpdateOrderStatus 更新订单状态
func (os *OrderService) UpdateOrderStatus(id int, status string) error {
	order := &model.Order{Status: status, Id: id}
	return os.oRepo.UpdateOrder(order)
}

// UpdateOrderAddress 更新订单地址
func (os *OrderService) UpdateOrderAddress(id int, address string) error {
	order := &model.Order{Address: address, Id: id}
	return os.oRepo.UpdateOrder(order)
}

// DeleteOrder 删除订单
func (os *OrderService) DeleteOrder(orderId int) error {
	return os.oRepo.DeleteOrder(orderId)
}

// GetOrderById 根据订单ID获取订单
func (os *OrderService) GetOrderById(orderId int) (*model.Order, error) {
	return os.oRepo.FindOrderById(orderId)
}

// GetOrdersByUserId 根据用户ID获取该用户的所有订单
func (os *OrderService) GetOrdersByUserId(userId int) ([]*model.Order, error) {
	return os.oRepo.FindOrdersByUserId(userId)
}
