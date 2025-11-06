package scheduler

import (
	"server/internal/product/order/service"
	"time"

	log "github.com/sirupsen/logrus"
)

// OrderDQScheduler 订单延迟队列调度器，定时处理超时订单并归还库存
type OrderDQScheduler struct {
	oCancelScheduler service.OrderCancelService
	stopChan         chan struct{}
}

// NewOrderDQScheduler 创建一个新的订单延迟队列调度器实例
func NewOrderDQScheduler(oCancelScheduler service.OrderCancelService) *OrderDQScheduler {
	return &OrderDQScheduler{oCancelScheduler: oCancelScheduler, stopChan: make(chan struct{})}
}

// Start 启动调度器，每10秒检查一次超时订单
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
