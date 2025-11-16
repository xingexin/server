package scheduler

import (
	"context"
	myRedis "server/pkg/redis"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

// RecoveryScheduler 超时任务恢复调度器
// 作用：定期扫描processing队列，将超时任务移回ready队列重试
//
// 工作原理：
// 1. 每60秒扫描一次processing队列
// 2. 找出超时的任务（处理时间 > 5分钟）
// 3. 将这些任务移回ready队列，60秒后重试
//
// 使用场景：
// - 调度器进程崩溃，任务留在processing队列
// - 任务处理时间过长（如网络慢、Redis超时）
// - 归还库存失败，任务卡在processing队列
//
// 三队列模型：
// - dq:ready: 待处理队列
// - dq:processing: 处理中队列（本调度器负责恢复）
// - dq:payload:{id}: 任务数据
type RecoveryScheduler struct {
	redisDB  *redis.Client
	stopChan chan struct{}
}

// NewRecoveryScheduler 创建一个新的恢复调度器实例
func NewRecoveryScheduler(redisDB *redis.Client) *RecoveryScheduler {
	return &RecoveryScheduler{
		redisDB:  redisDB,
		stopChan: make(chan struct{}),
	}
}

// Start 启动调度器，每60秒恢复一次超时任务
// 业务流程：
// 1. 创建60秒定时器
// 2. 每60秒触发一次超时任务恢复
// 3. 调用RecoverTimedOutTasks扫描processing队列
// 4. 将超时任务移回ready队列，延迟60秒后重试
// 5. 监听stopChan信号，收到信号后优雅退出
//
// 注意：
// - 此方法会阻塞，应在goroutine中运行
// - 恢复失败只记录日志，不会停止调度器
func (s *RecoveryScheduler) Start() error {
	// 创建60秒定时器
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()

	log.Info("Recovery scheduler started")
	count := 0
	for {
		select {
		case <-ticker.C:
			// 定时器触发，恢复超时任务
			// 参数60表示恢复的任务延迟60秒后重试
			recovered, err := myRedis.RecoverTimedOutTasks(context.Background(), s.redisDB, 60)
			if err != nil {
				log.Error("Recovery failed:", err)
			} else if recovered > 0 {
				log.Infof("Recovered %d timed out tasks", recovered)
			}
			count++
			log.Debug("Recovery scheduler finished, count:", count)
		case <-s.stopChan:
			// 收到停止信号，优雅退出
			log.Info("Recovery scheduler stopped")
			return nil
		}
	}
}

// Stop 停止调度器
func (s *RecoveryScheduler) Stop() {
	close(s.stopChan)
}
