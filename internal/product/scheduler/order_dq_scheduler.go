package scheduler

import (
	"server/internal/product/order/service"
	"time"

	log "github.com/sirupsen/logrus"
)

type OrderDQScheduler struct {
	oCancelScheduler service.OrderCancelService
	stopChan         chan struct{}
}

func NewOrderDQScheduler(oCancelScheduler service.OrderCancelService) *OrderDQScheduler {
	return &OrderDQScheduler{oCancelScheduler: oCancelScheduler, stopChan: make(chan struct{})}
}

func (s *OrderDQScheduler) Start() error {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	count := 0
	for {
		select {
		case <-ticker.C:
			if err := s.oCancelScheduler.RemoveTimeoutOrderTasks(); err != nil {
				log.Warn("Order DQ processing failed:", err)

			}
			count++
			log.Info("Order DQ processing scheduler finished", " count:", count)
		case <-s.stopChan:
			log.Info("Order DQ processing scheduler stopped")
			return nil
		}
	}
}
