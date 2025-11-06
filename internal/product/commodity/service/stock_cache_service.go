package service

import (
	"context"
	"server/internal/product/commodity/repository"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// StockCacheService 提供库存缓存相关的业务逻辑服务
type StockCacheService struct {
	cRedisSvc repository.StockCacheRepository
}

// NewStockCacheService 创建一个新的库存缓存服务实例
func NewStockCacheService(cRedisSvc repository.StockCacheRepository) *StockCacheService {
	return &StockCacheService{cRedisSvc: cRedisSvc}
}

// SyncAllStock 同步所有有变化的库存到数据库
func (s *StockCacheService) SyncAllStock(ctx context.Context) error {
	keys, err := s.cRedisSvc.GetAllDeltaKey(ctx)
	if err != nil {
		return err
	}
	validKeys := make([]int, 0)
	for _, key := range keys {
		value, err := s.cRedisSvc.GetDeltaValue(ctx, key)
		if err != nil {
			log.Warning("fail to get delta value " + err.Error())
			continue
		}
		if value != 0 {
			key = strings.TrimPrefix(key, "delta_key_")
			i, err := strconv.Atoi(key)
			if err != nil {
				log.Warning("invalid key " + err.Error())
				continue
			}
			validKeys = append(validKeys, i)
		}
	}
	for _, key := range validKeys {
		err = s.cRedisSvc.SyncStock(ctx, key)
		if err != nil {
			log.Warning("fail to sync " + err.Error())
			continue
		}
	}
	return nil
}
