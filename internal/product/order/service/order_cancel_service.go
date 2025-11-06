package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	commodityRepository "server/internal/product/commodity/repository"
	"server/internal/product/order/repository"
)

// OrderCancelService 订单取消服务接口，负责超时订单的自动取消和库存归还
type OrderCancelService interface {
	createOrderTask(orderId int, commodityId int, stock int) error // 创建订单取消任务（15分钟后执行）
	RemoveTimeoutOrderTasks() error                                // 处理超时订单，归还库存
}

type cancelService struct {
	cRedisRepo  commodityRepository.StockCacheRepository
	oRepo       repository.OrderRepository
	redisDQRepo repository.OrderDQRepository
	rDB         *redis.Client
}

// NewOrderCancelService 创建一个新的订单取消服务实例
func NewOrderCancelService(redisDQRepo repository.OrderDQRepository, oRepo repository.OrderRepository, cRedisRepo commodityRepository.StockCacheRepository, rDB *redis.Client) OrderCancelService {
	return &cancelService{
		redisDQRepo: redisDQRepo,
		oRepo:       oRepo,
		cRedisRepo:  cRedisRepo,
		rDB:         rDB,
	}
}

// createOrderTask 创建订单超时取消任务，15分钟后自动取消未支付订单
func (s *cancelService) createOrderTask(orderId int, commodityId int, stock int) error {
	// 将 commodityId 和 stock 用逗号分隔存储，格式: "commodityId,stock"
	payload := strconv.Itoa(commodityId) + "," + strconv.Itoa(stock)
	err := s.redisDQRepo.EnqueueDelayTask(context.TODO(), strconv.Itoa(orderId), payload, time.Minute*15)
	if err != nil {
		return err
	}
	return nil
}

// RemoveTimeoutOrderTasks 扫描并处理超时订单，自动归还库存到Redis
func (s *cancelService) RemoveTimeoutOrderTasks() error {
	ids, err := s.redisDQRepo.GetReadyTasks(context.TODO(), 100)

	// 处理预期的错误情况（这些不是真正的错误）
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoTasksDue):
			log.Debug("No tasks are due yet, waiting...")
			time.Sleep(time.Millisecond * 500)
			return nil
		case errors.Is(err, repository.ErrNoTasksInQueue):
			log.Debug("No tasks in queue, waiting...")
			time.Sleep(time.Second * 1)
			return nil
		case errors.Is(err, repository.ErrQueueOperationFailed):
			log.Warnf("Queue operation failed: %v", err)
			time.Sleep(time.Millisecond * 200)
			return err
		default:
			// 未知错误
			return err
		}
	}

	// 处理获取到的过期订单
	for _, id := range ids {
		orderId, err := strconv.Atoi(id)
		if err != nil {
			log.Warn("invalid order id:", err)
			continue
		}

		// 从 Redis 获取 payload: "commodityId,stock"
		payload, err := s.rDB.Get(context.TODO(), "dq:payload:"+id).Result()
		if err != nil {
			log.Warn("fail to get order payload:", err)
			continue
		}

		// 解析 payload
		parts := strings.Split(payload, ",")
		if len(parts) != 2 {
			log.Warnf("invalid payload format for order %d: %s", orderId, payload)
			continue
		}

		commodityId, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Warnf("invalid commodityId in payload: %s", parts[0])
			continue
		}

		stock, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Warnf("invalid stock in payload: %s", parts[1])
			continue
		}

		// 恢复库存
		err = s.cRedisRepo.IncreaseStock(context.TODO(), commodityId, stock)
		if err != nil {
			log.Errorf("Failed to increase stock for order %d: %v", orderId, err)
			return err
		}

		// 从延时队列移除
		err = s.redisDQRepo.RemoveTask(context.TODO(), id)
		if err != nil {
			log.Errorf("Failed to remove task %s: %v", id, err)
			return err
		}

		log.Infof("Successfully cancelled expired order %d and restored stock %d", orderId, stock)
	}
	return nil
}
