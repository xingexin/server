package scheduler

import (
	"server/internal/product/order/service"
	"time"

	log "github.com/sirupsen/logrus"
)

// OrderDQScheduler 订单延迟队列调度器，定时处理超时订单并归还库存
// 工作原理：
// 1. 每10秒扫描一次Redis延迟队列（ZSet）
// 2. 获取所有到期的订单任务（score小于当前时间戳）
// 3. 对每个到期订单：归还库存到Redis、从延迟队列移除任务
// 4. 记录处理日志
//
// 设计思想：
// - 订单创建时加入延迟队列，15分钟后到期
// - 调度器定时扫描到期任务，自动取消未支付订单
// - 归还库存只操作Redis，不影响已创建的订单记录
type OrderDQScheduler struct {
	oCancelScheduler service.OrderCancelService // 订单取消服务
	stopChan         chan struct{}               // 停止信号channel
}

// NewOrderDQScheduler 创建一个新的订单延迟队列调度器实例
func NewOrderDQScheduler(oCancelScheduler service.OrderCancelService) *OrderDQScheduler {
	return &OrderDQScheduler{oCancelScheduler: oCancelScheduler, stopChan: make(chan struct{})}
}

// Start 启动调度器，每10秒检查一次超时订单
// 业务流程：
// 1. 创建10秒定时器
// 2. 每10秒触发一次超时订单处理
// 3. 调用OrderCancelService.RemoveTimeoutOrderTasks处理到期任务
// 4. 记录处理次数和错误日志
// 5. 监听stopChan信号，收到信号后优雅退出
//
// 注意：
// - 此方法会阻塞，应在goroutine中运行
// - 处理失败只记录日志，不会停止调度器
// - 每次处理最多100个到期任务
func (s *OrderDQScheduler) Start() error {
	// 创建10秒定时器
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	count := 0
	for {
		select {
		case <-ticker.C:
			// 定时器触发，处理超时订单
			if err := s.oCancelScheduler.RemoveTimeoutOrderTasks(); err != nil {
				log.Warn("Order DQ processing failed:", err)
			}
			count++
			log.Info("Order DQ processing scheduler finished", " count:", count)
		case <-s.stopChan:
			// 收到停止信号，优雅退出
			log.Info("Order DQ processing scheduler stopped")
			return nil
		}
	}
}
