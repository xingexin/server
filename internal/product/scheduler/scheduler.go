package scheduler

import (
	"context"
	"server/internal/product/commodity/service"
	"time"

	log "github.com/sirupsen/logrus"
)

type Scheduler struct {
	cStockSvc *service.StockCacheService
	stopChan  chan struct{}
}

func NewScheduler(cStockSvc *service.StockCacheService) *Scheduler {
	return &Scheduler{
		cStockSvc: cStockSvc,
		stopChan:  make(chan struct{}),
	}
}

func (s *Scheduler) Start() error {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	log.Info("Stock sync scheduler started")

	for {
		select {
		case <-ticker.C:
			if err := s.cStockSvc.SyncAllStock(context.Background()); err != nil {
				log.Error("Stock sync failed:", err)
			}
		case <-s.stopChan:
			log.Info("Stock sync scheduler stopped")
			return nil
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
}
