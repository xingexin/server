package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	myRedis "server/pkg/redis"
)

// OrderDQRepository 订单延迟队列的数据访问接口
type OrderDQRepository interface {
	EnqueueDelayTask(ctx context.Context, id, payload string, execTime time.Duration) error // 将延迟任务加入队列
	GetReadyTasks(ctx context.Context, count int64) ([]string, error)                       // 获取到期的任务
	RemoveTask(ctx context.Context, id string) error                                        // 从队列中移除任务
}

type redisOrderDQRepository struct {
	redisDB *redis.Client
}

// NewOrderDQRepository 创建一个新的订单延迟队列仓储实例
func NewOrderDQRepository(rDB *redis.Client) OrderDQRepository {
	return &redisOrderDQRepository{redisDB: rDB}
}

// EnqueueDelayTask 将订单延迟任务加入Redis队列
func (oRedisRepo *redisOrderDQRepository) EnqueueDelayTask(ctx context.Context, id, payload string, execTime time.Duration) error {
	unixTime := int64(execTime.Seconds()) + oRedisRepo.redisDB.Time(ctx).Val().Unix()
	return myRedis.EnqueueDelayTask(ctx, oRedisRepo.redisDB, id, payload, unixTime)
}

// GetReadyTasks 获取到期的延迟任务（已超时的订单）
func (oRedisRepo *redisOrderDQRepository) GetReadyTasks(ctx context.Context, count int64) ([]string, error) {
	head, err := oRedisRepo.redisDB.ZRangeWithScores(ctx, "dq:ready", 0, 0).Result()

	if err != nil {
		log.Warn("Failed to get ready tasks: %v", err)
		return nil, ErrQueueOperationFailed
	}
	if len(head) == 0 {
		log.Debug("No ready tasks found")
		return nil, ErrNoTasksInQueue
	}
	dueTime := int64(head[0].Score)
	redisNow := oRedisRepo.redisDB.Time(ctx).Val().Unix()
	if dueTime > redisNow {
		log.Debug("No tasks are due yet")
		return nil, ErrNoTasksDue
	}
	ids, err := oRedisRepo.redisDB.ZRangeByScore(ctx, "dq:ready", &redis.ZRangeBy{
		Min:    "-inf",
		Max:    fmt.Sprintf("%d", redisNow),
		Offset: 0,
		Count:  count,
	}).Result()
	if err != nil {
		log.Warn("Failed to get ready tasks by score: %v", err)
		return nil, ErrQueueOperationFailed
	}
	return ids, nil
}

// RemoveTask 从延迟队列中移除已处理的任务
func (oRedisRepo *redisOrderDQRepository) RemoveTask(ctx context.Context, id string) error {
	return myRedis.Ack(ctx, oRedisRepo.redisDB, id)
}
