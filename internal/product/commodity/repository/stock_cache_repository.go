package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

func getStockCacheKey(commodityId int) string {
	return "commodity_stock_" + strconv.Itoa(commodityId)
}
func getDeltaCacheKey(commodityId int) string {
	return "commodity_delta_" + strconv.Itoa(commodityId)
}

type StockCacheRepository interface {
	initStockCache(ctx context.Context, commodityId int, stock int) error
	decreaseStock(ctx context.Context, commodityId int, quantity int) (bool, error)
	increaseStock(ctx context.Context, commodityId int, quantity int) error
	syncStock(ctx context.Context, commodityId int, stock int) error
}

type RedisCommodityRepository struct {
	cRedisRepo *redis.Client
}

func NewRedisCommodityRepository(cRedisRepo *redis.Client) StockCacheRepository {
	return &RedisCommodityRepository{cRedisRepo: cRedisRepo}
}

func (rRepo *RedisCommodityRepository) initStockCache(ctx context.Context, commodityId int, stock int) error {
	if err := rRepo.cRedisRepo.Set(ctx, getStockCacheKey(commodityId), stock, time.Hour*24).Err(); err != nil {
		log.Error("Failed to initialize stock cache:", err)
		return err
	}
	log.Debug("Initialized stock cache for commodity ID", commodityId, "with stock", stock)
	return nil
}

func (rRepo *RedisCommodityRepository) decreaseStock(ctx context.Context, commodityId int, quantity int) (bool, error) {
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
