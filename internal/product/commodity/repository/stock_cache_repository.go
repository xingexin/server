package repository

import (
	"context"
	"fmt"
	"server/internal/product/commodity/model"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// getStockCacheKey 生成库存缓存的Redis key
// 格式：stock_key_{商品ID}
// 例如：stock_key_123 表示商品ID为123的库存
// 用途：存储商品的实时库存数量
func getStockCacheKey(commodityId int) string {
	return "stock_key_" + strconv.Itoa(commodityId)
}

// getDeltaCacheKey 生成库存增量的Redis key
// 格式：delta_key_{商品ID}
// 例如：delta_key_123 表示商品ID为123的库存变化量
// 用途：存储Redis与MySQL之间的库存差值，用于异步同步
func getDeltaCacheKey(commodityId int) string {
	return "delta_key_" + strconv.Itoa(commodityId)
}

type StockCacheRepository interface {
	InitStockCache(ctx context.Context, commodityId int, stock int) error          // 初始化商品库存缓存
	DecreaseStock(ctx context.Context, commodityId int, quantity int) (int, error) // 扣减库存（原子操作）
	IncreaseStock(ctx context.Context, commodityId int, quantity int) error        // 增加库存（用于归还）
	SyncStock(ctx context.Context, commodityId int) error                          // 同步库存增量到数据库
	GetAllDeltaKey(ctx context.Context) ([]string, error)                          // 获取所有有变化的库存key
	GetDeltaValue(ctx context.Context, key string) (int, error)                    // 获取库存增量值
}

type redisCommodityRepository struct {
	cRedisRepo *redis.Client
	cRepo      *gorm.DB
}

// NewRedisCommodityRepository 创建一个新的Redis库存缓存仓储实例
func NewRedisCommodityRepository(cRedisRepo *redis.Client, cRepo *gorm.DB) StockCacheRepository {
	return &redisCommodityRepository{cRedisRepo: cRedisRepo, cRepo: cRepo}
}

// InitStockCache 初始化商品库存缓存到Redis
// 创建两个key：stock_key_商品ID 和 delta_key_商品ID
// stock_key_商品ID：存储商品的实时库存数量
// delta_key_商品ID：存储Redis与MySQL之间的库存差值，初始值为0
// 创建delta_key目的是为了不影响主动修改库存的业务场景
// 例如：管理员手动调整库存时，不应影响订单扣减的库存同步逻辑
func (rRepo *redisCommodityRepository) InitStockCache(ctx context.Context, commodityId int, stock int) error {
	if err := rRepo.cRedisRepo.Set(ctx, getStockCacheKey(commodityId), stock, time.Hour*24).Err(); err != nil {
		log.Error("Failed to initialize stock cache:", err)
		return err
	}
	log.Debug("Initialized stock cache for commodity ID", commodityId, "with stock", stock)
	return nil
}

// DecreaseStock 使用Lua脚本原子性地扣减库存，防止超卖
func (rRepo *redisCommodityRepository) DecreaseStock(ctx context.Context, commodityId int, quantity int) (int, error) {
	if rRepo == nil || rRepo.cRedisRepo == nil {
		panic("redis repository is nil")
	}
	if quantity <= 0 {
		log.Warning("Invalid quantity:", quantity)
		return -1, fmt.Errorf("invalid quantity %d", quantity)
	}
	stockKey := getStockCacheKey(commodityId)
	deltaKey := getDeltaCacheKey(commodityId)

	luaScript := `
	local stock_key = KEYS[1]
	local delta_key = KEYS[2]
	local quantity = tonumber(ARGV[1])

	local current_stock = tonumber(redis.call("GET", stock_key))
	if not current_stock then
		return -2
	end

	if current_stock < quantity then
		return -3
	end
	redis.call("DECRBY", stock_key, quantity)
	redis.call("INCRBY", delta_key, quantity)
	redis.call("EXPIRE", delta_key, 86400)
	return current_stock - quantity
`
	// 执行Lua脚本进行原子性库存扣减
	result, err := rRepo.cRedisRepo.Eval(ctx, luaScript, []string{stockKey, deltaKey}, quantity).Result()

	// 错误码说明：
	// - code=非-2 -3值: 扣减成功
	// - code=1: Redis执行失败（网络错误、Redis宕机等）
	// - code=2: 缓存未初始化（触发缓存穿透处理，需从MySQL加载）
	// - code=3: 库存不足（防止超卖，订单创建失败）

	if err != nil {
		// code=1: Redis执行失败
		// 处理方式：返回错误，订单创建失败
		log.Error("Failed to decrease stock:", err)
		return 1, err
	}

	if result.(int64) == -2 {
		// code=-2: Redis缓存未初始化（缓存穿透）
		// 处理方式：调用方需从MySQL查询库存并调用InitStockCache初始化，然后重新扣减
		log.Warning("Stock cache not initialized for commodity ID", commodityId)
		return 2, fmt.Errorf("stock cache not initialized")
	}

	if result.(int64) == -3 {
		// code=-3: 库存不足
		// 处理方式：直接返回错误，拒绝创建订单，防止超卖
		log.Warning("Insufficient stock for commodity ID", commodityId)
		return 3, fmt.Errorf("insufficient stock")
	}

	// code=0: 扣减成功
	log.Debug("Decreased stock for commodity ID", commodityId, "by", quantity)
	return 0, nil
}

// IncreaseStock 使用Lua脚本原子性地增加库存（用于订单取消）
func (rRepo *redisCommodityRepository) IncreaseStock(ctx context.Context, commodityId int, quantity int) error {
	if quantity <= 0 {
		log.Warning("Invalid quantity:", quantity)
		return fmt.Errorf("invalid quantity %d", quantity)
	}
	stockKey := getStockCacheKey(commodityId)
	deltaKey := getDeltaCacheKey(commodityId)

	luaScript := `
	local stock_key = KEYS[1]
	local delta_key = KEYS[2]
	local quantity = tonumber(ARGV[1])
	local current_stock = redis.call("GET", stock_key)
	if not current_stock then
		return -1
	end

	redis.call("INCRBY", stock_key, quantity)
	redis.call("DECRBY", delta_key, quantity)
	return current_stock + quantity
`
	result, err := rRepo.cRedisRepo.Eval(ctx, luaScript, []string{stockKey, deltaKey}, quantity).Result()
	if err != nil {
		log.Error("Failed to increase stock:", err)
		return err
	}
	if result.(int64) == -1 {
		log.Warning("Stock cache not initialized for commodity ID", commodityId)
		return fmt.Errorf("stock cache not initialized")
	}
	log.Debug("Increased stock for commodity ID", commodityId, "by", quantity)
	return nil
}

// SyncStock 将Redis中的库存增量同步到MySQL数据库
func (rRepo *redisCommodityRepository) SyncStock(ctx context.Context, commodityId int) error {
	deltaKey := getDeltaCacheKey(commodityId)
	// 获取并重置库存增量
	luaScript := `
	local delta_key = KEYS[1]
	local delta = redis.call("GET", delta_key)
	if not delta then
		return -1
	end
	redis.call("SET", delta_key, 0)
	redis.call("EXPIRE", delta_key, 86400)
	return tonumber(delta)
`
	// 执行Lua脚本：原子性获取delta并重置为0
	delta, err := rRepo.cRedisRepo.Eval(ctx, luaScript, []string{deltaKey}).Result()

	// Delta值说明：
	// - delta=-1: delta_key不存在（商品从未被扣减过库存）
	// - delta=0:  没有库存变化，无需同步
	// - delta>0:  有库存扣减，需要同步到MySQL（UPDATE stock = stock - delta）

	if err != nil {
		// Redis执行失败
		log.Error("Failed to execute sync stock Lua script:", err)
		return err
	}

	if delta.(int64) == -1 {
		// delta_key不存在：商品从未有订单扣减库存
		// 处理方式：无需同步，直接返回
		log.Debug("No delta to sync for commodity ID", commodityId)
		return nil
	}

	if delta.(int64) == 0 {
		// delta为0：上次同步后没有新的库存变化
		// 处理方式：无需同步，避免无效的数据库UPDATE操作
		log.Debug("Delta is zero, no need to sync for commodity ID", commodityId)
		return nil
	}

	// 同步到数据库
	result := rRepo.cRepo.Model(&model.Commodity{}).Where("id = ?", commodityId).Update("stock", gorm.Expr("stock - ?", delta))

	if result.RowsAffected == 0 {
		rRepo.cRedisRepo.IncrBy(ctx, deltaKey, delta.(int64))

		log.Warning("Commodity not found in database for ID", commodityId)
		return fmt.Errorf("commodity with id=%d not found", commodityId)
	}
	if result.Error != nil {
		rRepo.cRedisRepo.IncrBy(ctx, deltaKey, delta.(int64))
		log.Error("Failed to sync stock to database:", result.Error)
		return result.Error
	}

	return nil
}

// GetAllDeltaKey 扫描并获取所有库存增量的key
func (rRepo *redisCommodityRepository) GetAllDeltaKey(ctx context.Context) ([]string, error) {
	res := make([]string, 0)
	iter := rRepo.cRedisRepo.Scan(ctx, 0, "delta_key_*", 100).Iterator() //分100条每页
	for iter.Next(ctx) {
		res = append(res, iter.Val())
	}
	if err := iter.Err(); err != nil {
		log.Error("Failed to scan delta keys:", err)
		return nil, err
	}
	return res, nil
}

// GetDeltaValue 获取指定key的库存增量值
func (rRepo *redisCommodityRepository) GetDeltaValue(ctx context.Context, key string) (int, error) {
	res, err := rRepo.cRedisRepo.Get(ctx, key).Result()
	if err != nil {
		log.Error("Failed to get delta value:", err)
		return 0, err
	}
	delta, err := strconv.Atoi(res)
	if err != nil {
		log.Error("Failed to convert delta value to int:", err)
		return 0, err
	}
	return delta, nil
}
