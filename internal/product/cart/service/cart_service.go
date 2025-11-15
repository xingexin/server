package service

import (
	"server/internal/product/cart/model"
	"server/internal/product/cart/repository"
	"time"
)

// CartService 提供购物车相关的业务逻辑服务
type CartService struct {
	cartRepo repository.CartRepository
}

// NewCartService 创建一个新的购物车服务实例
func NewCartService(repo repository.CartRepository) *CartService {
	return &CartService{cartRepo: repo}
}

// AddToCart 添加商品到购物车
// 业务流程：
// 1. 构建购物车项模型对象
// 2. 设置用户ID、商品ID、数量
// 3. 设置创建时间和更新时间为当前时间
// 4. 调用Repository层创建购物车记录
//
// 注意：
// - 此方法直接创建新记录，不会检查是否已存在相同商品
// - 如需累加数量，应先查询是否存在，存在则调用UpdateCart
// - 添加时不会验证商品是否存在或库存是否充足
// - 数量校验应在业务层或Handler层完成
func (cs *CartService) AddToCart(userId int, commodityId int, quantity int) error {
	// 构建购物车项模型对象
	cart := &model.Cart{
		UserId:      userId,      // 用户ID
		CommodityId: commodityId, // 商品ID
		Quantity:    quantity,    // 商品数量
		CreatedAt:   time.Now(),  // 创建时间
		UpdatedAt:   time.Now(),  // 更新时间
	}

	return cs.cartRepo.CreateCart(cart)
}

// RemoveFromCart 从购物车中移除商品
func (cs *CartService) RemoveFromCart(cartId int) error {
	return cs.cartRepo.DeleteCart(cartId)
}

// UpdateCart 更新购物车中商品的数量
// 业务流程：
// 1. 根据购物车ID查询购物车项
// 2. 更新商品数量
// 3. 设置更新时间为当前时间
// 4. 调用Repository层更新购物车记录
//
// 注意：
// - 数量可以增加或减少
// - 如果数量设为0，建议使用RemoveFromCart方法删除记录
// - 不会验证库存是否充足，下单时才验证
func (cs *CartService) UpdateCart(cartId int, quantity int) error {
	// 查询购物车项
	cart, err := cs.cartRepo.FindCartById(cartId)
	if err != nil {
		return err
	}

	// 更新数量和更新时间
	cart.Quantity = quantity
	cart.UpdatedAt = time.Now()

	return cs.cartRepo.UpdateCart(cart)
}

// GetCart 获取用户的购物车
func (cs *CartService) GetCart(userId int) ([]*model.Cart, error) {
	return cs.cartRepo.FindCartByUserId(userId)
}
