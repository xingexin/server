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

func getStockCacheKey(commodityId int) string {
	return "stock_key_" + strconv.Itoa(commodityId)
}
func getDeltaCacheKey(commodityId int) string {
	return "delta_key_" + strconv.Itoa(commodityId)
}

type StockCacheRepository interface {
	InitStockCache(ctx context.Context, commodityId int, stock int) error
	DecreaseStock(ctx context.Context, commodityId int, quantity int) (bool, error)
	IncreaseStock(ctx context.Context, commodityId int, quantity int) error
	SyncStock(ctx context.Context, commodityId int) error
	GetAllDeltaKey(ctx context.Context) ([]string, error)
	GetDeltaValue(ctx context.Context, key string) (int, error)
}

type RedisCommodityRepository struct {
	cRedisRepo *redis.Client
	cRepo      *gorm.DB
}

func NewRedisCommodityRepository(cRedisRepo *redis.Client, cRepo *gorm.DB) StockCacheRepository {
	return &RedisCommodityRepository{cRedisRepo: cRedisRepo, cRepo: cRepo}
}

func (rRepo *RedisCommodityRepository) InitStockCache(ctx context.Context, commodityId int, stock int) error {
	if err := rRepo.cRedisRepo.Set(ctx, getStockCacheKey(commodityId), stock, time.Hour*24).Err(); err != nil {
		log.Error("Failed to initialize stock cache:", err)
		return err
	}
	log.Debug("Initialized stock cache for commodity ID", commodityId, "with stock", stock)
	return nil
}

func (rRepo *RedisCommodityRepository) DecreaseStock(ctx context.Context, commodityId int, quantity int) (bool, error) {
	if quantity <= 0 {
		log.Warning("Invalid quantity:", quantity)
		return false, fmt.Errorf("invalid quantity %d", quantity)
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

	if current_stock < quantity then
		return -2
	end
	redis.call("DECRBY", stock_key, quantity)
	redis.call("INCRBY", delta_key, quantity)
	redis.call("EXPIRE", delta_key, 86400)
	return current_stock - quantity
`
	result, err := rRepo.cRedisRepo.Eval(ctx, luaScript, []string{stockKey, deltaKey}, quantity).Result()
	if err != nil {
		log.Error("Failed to decrease stock:", err)
		return false, err
	}
	if result.(int64) == -1 {
		log.Warning("Stock cache not initialized for commodity ID", commodityId)
		return false, fmt.Errorf("stock cache not initialized")
	}
	if result.(int64) == -2 {
		log.Warning("Insufficient stock for commodity ID", commodityId)
		return false, fmt.Errorf("insufficient stock")
	}
	log.Debug("Decreased stock for commodity ID", commodityId, "by", quantity)
	return true, nil
}

func (rRepo *RedisCommodityRepository) IncreaseStock(ctx context.Context, commodityId int, quantity int) error {
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

func (rRepo *RedisCommodityRepository) SyncStock(ctx context.Context, commodityId int) error {
	deltaKey := getDeltaCacheKey(commodityId)
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
	delta, err := rRepo.cRedisRepo.Eval(ctx, luaScript, []string{deltaKey}).Result()
	if err != nil {
		log.Error("Failed to execute sync stock Lua script:", err)
		return err
	}
	if delta.(int64) == -1 {
		log.Debug("No delta to sync for commodity ID", commodityId)
		return nil
	}
	if delta.(int64) == 0 {
		log.Debug("Delta is zero, no need to sync for commodity ID", commodityId)
		return nil
	}
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

func (rRepo *RedisCommodityRepository) GetAllDeltaKey(ctx context.Context) ([]string, error) {
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

func (rRepo *RedisCommodityRepository) GetDeltaValue(ctx context.Context, key string) (int, error) {
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
