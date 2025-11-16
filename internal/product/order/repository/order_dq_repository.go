package repository

import (
	"context"

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

// GetReadyTasks 获取到期的延迟任务（已超时的订单）并原子性移动到processing队列
// 使用三队列模型防止任务重复处理：
// 1. 原子性从ready队列获取到期任务
// 2. 立即移动到processing队列（防止其他调度器重复获取）
// 3. 处理超时设置为300秒，超时后任务会被RecoveryScheduler恢复
func (oRedisRepo *redisOrderDQRepository) GetReadyTasks(ctx context.Context, count int64) ([]string, error) {
	// 使用GetAndMoveToProcessing原子性获取任务并移动到processing队列
	// 处理超时设置为300秒（5分钟），超时后任务会被RecoveryScheduler恢复到ready队列
	ids, err := myRedis.GetAndMoveToProcessing(ctx, oRedisRepo.redisDB, count, 300)
	if err != nil {
		log.Warn("Failed to get and move tasks: %v", err)
		return nil, ErrQueueOperationFailed
	}

	// 如果没有到期任务，返回特定错误
	if len(ids) == 0 {
		log.Debug("No ready tasks found")
		return nil, ErrNoTasksInQueue
	}

	return ids, nil
}

// RemoveTask 从延迟队列中移除已处理的任务
func (oRedisRepo *redisOrderDQRepository) RemoveTask(ctx context.Context, id string) error {
	return myRedis.Ack(ctx, oRedisRepo.redisDB, id)
}
