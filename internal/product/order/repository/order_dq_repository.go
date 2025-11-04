package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	myRedis "server/pkg/redis"
)

type OrderDQRepository interface {
	EnqueueDelayTask(ctx context.Context, id, payload string, execTime time.Duration) error
	GetReadyTasks(ctx context.Context, count int64) ([]string, error)
	RemoveTask(ctx context.Context, id string) error
}

type redisOrderDQRepository struct {
	redisDB *redis.Client
}

func NewOrderDQRepository(rDB *redis.Client) OrderDQRepository {
	return &redisOrderDQRepository{redisDB: rDB}
}

func (oRedisRepo *redisOrderDQRepository) EnqueueDelayTask(ctx context.Context, id, payload string, execTime time.Duration) error {
	unixTime := int64(execTime.Seconds()) + oRedisRepo.redisDB.Time(ctx).Val().Unix()
	return myRedis.EnqueueDelayTask(ctx, oRedisRepo.redisDB, id, payload, unixTime)
}

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

func (oRedisRepo *redisOrderDQRepository) RemoveTask(ctx context.Context, id string) error {
	return myRedis.Ack(ctx, oRedisRepo.redisDB, id)
}
