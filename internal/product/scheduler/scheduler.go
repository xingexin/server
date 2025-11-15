package scheduler

import (
	"context"
	"server/internal/product/commodity/service"
	"time"

	log "github.com/sirupsen/logrus"
)

// Scheduler 库存同步调度器，定时将Redis库存变化同步到MySQL
// 工作原理：
// 1. 每10秒扫描一次所有delta_key（库存变化记录）
// 2. 将有变化的库存批量同步到MySQL
// 3. 同步成功后清零delta值
//
// 设计思想：
// - 订单扣减库存时只操作Redis（快速响应）
// - 调度器异步批量同步到MySQL（减轻数据库压力）
// - 最终一致性：Redis为实时数据，MySQL定期同步
type Scheduler struct {
	cStockSvc *service.StockCacheService // 库存缓存服务
	stopChan  chan struct{}               // 停止信号channel
}

// NewScheduler 创建一个新的库存同步调度器实例
func NewScheduler(cStockSvc *service.StockCacheService) *Scheduler {
	return &Scheduler{
		cStockSvc: cStockSvc,
		stopChan:  make(chan struct{}),
	}
}

// Start 启动调度器，每10秒同步一次库存
// 业务流程：
// 1. 创建10秒定时器
// 2. 每10秒触发一次库存同步
// 3. 调用StockCacheService.SyncAllStock批量同步所有有变化的库存
// 4. 记录同步次数和错误日志
// 5. 监听stopChan信号，收到信号后优雅退出
//
// 注意：
// - 此方法会阻塞，应在goroutine中运行
// - 同步失败只记录日志，不会停止调度器
// - 使用select可以响应停止信号，实现优雅关闭
func (s *Scheduler) Start() error {
	// 创建10秒定时器
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	log.Info("Stock sync scheduler started")
	count := 0
	for {
		select {
		case <-ticker.C:
			// 定时器触发，执行库存同步
			if err := s.cStockSvc.SyncAllStock(context.Background()); err != nil {
				log.Error("Stock sync failed:", err)
			}
			count++
			log.Info("Stock sync scheduler finished", " count:", count)
		case <-s.stopChan:
			// 收到停止信号，优雅退出
			log.Info("Stock sync scheduler stopped")
			return nil
		}
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	close(s.stopChan)
}
