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

// GetAndMoveToProcessing 原子性地从ready队列获取到期任务并移动到processing队列
// 使用Lua脚本保证：同一任务只会被一个调度器获取，防止重复处理
//
// 参数：
// - count: 最多获取多少个任务
// - processingTimeout: 处理超时时间（秒），超时后任务会被恢复到ready队列
//
// 返回：
// - []string: 任务ID列表
func GetAndMoveToProcessing(ctx context.Context, rdb *redis.Client, count int64, processingTimeout int64) ([]string, error) {
	now := time.Now().Unix()
	timeoutScore := now + processingTimeout // 处理超时时间戳

	luaScript := `
	local ready_key = "dq:ready"
	local processing_key = "dq:processing"
	local now = tonumber(ARGV[1])
	local timeout_score = tonumber(ARGV[2])
	local count = tonumber(ARGV[3])

	-- 获取到期的任务（score <= now）
	local tasks = redis.call("ZRANGEBYSCORE", ready_key, "-inf", now, "LIMIT", 0, count)

	if #tasks == 0 then
		return {}
	end

	-- 原子性操作：从ready删除，加入processing
	for i, task in ipairs(tasks) do
		redis.call("ZREM", ready_key, task)
		redis.call("ZADD", processing_key, timeout_score, task)
	end

	return tasks
	`

	result, err := rdb.Eval(ctx, luaScript, []string{}, now, timeoutScore, count).Result()
	if err != nil {
		log.Errorf("Failed to get and move tasks: %v", err)
		return nil, err
	}

	// 转换结果为字符串数组
	tasks := make([]string, 0)
	if arr, ok := result.([]interface{}); ok {
		for _, v := range arr {
			if str, ok := v.(string); ok {
				tasks = append(tasks, str)
			}
		}
	}

	return tasks, nil
}

// RecoverTimedOutTasks 将processing队列中超时的任务移回ready队列
// 应该定期调用（建议每60秒一次）
//
// 工作原理：
// 1. 扫描processing队列中score < now的任务（已超时）
// 2. 将这些任务移回ready队列，延迟retryDelaySeconds秒后重试
// 3. 返回恢复的任务数量
//
// 参数：
// - retryDelaySeconds: 重试延迟时间（秒）
//
// 返回：
// - int: 恢复的任务数量
func RecoverTimedOutTasks(ctx context.Context, rdb *redis.Client, retryDelaySeconds int64) (int, error) {
	now := time.Now().Unix()

	luaScript := `
	local processing_key = "dq:processing"
	local ready_key = "dq:ready"
	local now = tonumber(ARGV[1])
	local retry_delay = tonumber(ARGV[2])

	-- 获取超时的任务（score < now）
	local tasks = redis.call("ZRANGEBYSCORE", processing_key, "-inf", now)

	if #tasks == 0 then
		return 0
	end

	-- 移回ready队列，延迟retry_delay秒后重试
	for i, task in ipairs(tasks) do
		redis.call("ZREM", processing_key, task)
		redis.call("ZADD", ready_key, now + retry_delay, task)
	end

	return #tasks
	`

	result, err := rdb.Eval(ctx, luaScript, []string{}, now, retryDelaySeconds).Result()
	if err != nil {
		log.Errorf("Failed to recover timed out tasks: %v", err)
		return 0, err
	}

	count := int(result.(int64))
	if count > 0 {
		log.Warnf("Recovered %d timed out tasks from processing to ready", count)
	}

	return count, nil
}

// Ack 确认任务成功完成，从processing队列删除并清理任务数据
func Ack(ctx context.Context, rdb *redis.Client, id string) error {
	pipe := rdb.TxPipeline()
	pipe.ZRem(ctx, "dq:processing", id) // 从processing队列删除
	pipe.Del(ctx, "dq:payload:"+id)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("Failed to acknowledge job %s: %v", id, err)
		return err
	}
	return nil
}
