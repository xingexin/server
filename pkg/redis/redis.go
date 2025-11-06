package redis

import (
	"context"
	"server/config"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

// InitRedis 初始化 Redis 连接池
func InitRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return nil, err
	}
	return rdb, nil
}

// EnqueueDelayTask 将延迟任务加入 Redis 延迟队列（使用 ZSet 实现）
func EnqueueDelayTask(ctx context.Context, rdb *redis.Client, id, payload string, time int64) error {
	pipe := rdb.TxPipeline()
	pipe.ZAdd(ctx, "dq:ready", redis.Z{Score: float64(time), Member: id})
	pipe.Set(ctx, "dq:payload:"+id, payload, 0)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("Failed to enqueue job %s: %v", id, err)
		return err
	}
	return nil
}

// Ack 确认并删除延迟队列中的任务
func Ack(ctx context.Context, rdb *redis.Client, id string) error {
	pipe := rdb.TxPipeline()
	pipe.ZRem(ctx, "dq:ready", id)
	pipe.Del(ctx, "dq:payload:"+id)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("Failed to acknowledge job %s: %v", id, err)
		return err
	}
	return nil
}
