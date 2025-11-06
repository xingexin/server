package scheduler

import (
	"context"
	"server/internal/product/commodity/service"
	"time"

	log "github.com/sirupsen/logrus"
)

// Scheduler 库存同步调度器，定时将Redis库存变化同步到MySQL
type Scheduler struct {
	cStockSvc *service.StockCacheService
	stopChan  chan struct{}
}

// NewScheduler 创建一个新的库存同步调度器实例
func NewScheduler(cStockSvc *service.StockCacheService) *Scheduler {
	return &Scheduler{
		cStockSvc: cStockSvc,
		stopChan:  make(chan struct{}),
	}
}

// Start 启动调度器，每10秒同步一次库存
func (s *Scheduler) Start() error {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	log.Info("Stock sync scheduler started")
	count := 0
	for {
		select {
		case <-ticker.C:
			if err := s.cStockSvc.SyncAllStock(context.Background()); err != nil {
				log.Error("Stock sync failed:", err)
			}
			count++
			log.Info("Stock sync scheduler finished", " count:", count)
		case <-s.stopChan:
			log.Info("Stock sync scheduler stopped")
			return nil
		}
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	close(s.stopChan)
}
