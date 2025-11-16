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
// 业务流程：
// 1. 构建payload字符串（格式："commodityId,stock"）
// 2. 将任务加入Redis延迟队列（使用ZSet实现，score为执行时间戳）
// 3. 15分钟后，订单取消调度器会扫描到期任务并处理
//
// 延迟队列实现：
// - 使用Redis ZSet存储任务，key为"dq:ready"
// - score为任务执行时间的Unix时间戳（当前时间 + 15分钟）
// - member为订单ID
// - payload单独存储在"dq:payload:{orderId}"中
//
// 参数说明：
// - orderId: 订单ID，用作延迟队列的任务ID
// - commodityId: 商品ID，用于归还库存时定位商品
// - stock: 需要归还的库存数量
func (s *cancelService) createOrderTask(orderId int, commodityId int, stock int) error {
	// 构建payload字符串，格式："commodityId,stock"
	// 例如："123,5" 表示商品ID=123，数量=5
	payload := strconv.Itoa(commodityId) + "," + strconv.Itoa(stock)

	// 将任务加入延迟队列，15分钟后执行
	err := s.redisDQRepo.EnqueueDelayTask(context.TODO(), strconv.Itoa(orderId), payload, time.Minute*15)
	if err != nil {
		return err
	}

	return nil
}

// RemoveTimeoutOrderTasks 扫描并处理超时订单，自动归还库存到Redis
// 业务流程：
//  1. 从延迟队列中获取到期的任务（最多100个）
//  2. 遍历每个到期任务：
//     a. 设置幂等性键（防止重复处理）
//     b. 解析任务ID（订单ID）
//     c. 从Redis获取payload（包含商品ID和库存数量）
//     d. 调用IncreaseStock归还库存到Redis
//     e. 从延迟队列中移除任务
//  3. 记录处理日志
//
// 错误处理：
// - ErrNoTasksDue: 没有到期任务，等待500ms后返回
// - ErrNoTasksInQueue: 队列为空，等待1秒后返回
// - ErrQueueOperationFailed: 队列操作失败，等待200ms后返回错误
//
// 幂等性保护：
// - 使用Redis SetNX设置幂等性键（格式：order_cancel_idempotent:{orderId}）
// - 如果幂等性键已存在，说明订单已处理过，只删除任务不归还库存
// - 幂等性键24小时后自动过期
// - 解决了"归还库存成功但删除任务失败"导致的重复归还问题
func (s *cancelService) RemoveTimeoutOrderTasks() error {
	// 从延迟队列获取到期任务（最多100个）
	ids, err := s.redisDQRepo.GetReadyTasks(context.TODO(), 100)

	// 处理预期的错误情况（这些不是真正的错误，只是队列状态）
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoTasksDue):
			// 没有到期任务，暂停500ms避免空转
			log.Debug("No tasks are due yet, waiting...")
			time.Sleep(time.Millisecond * 500)
			return nil
		case errors.Is(err, repository.ErrNoTasksInQueue):
			// 队列为空，暂停1秒避免空转
			log.Debug("No tasks in queue, waiting...")
			time.Sleep(time.Second * 1)
			return nil
		case errors.Is(err, repository.ErrQueueOperationFailed):
			// 队列操作失败（Redis错误），暂停200ms后重试
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
		// 解析订单ID
		orderId, err := strconv.Atoi(id)
		if err != nil {
			log.Warn("invalid order id:", err)
			continue // 跳过无效ID，继续处理下一个
		}

		// 幂等性保护：尝试设置幂等性键
		// 格式：order_cancel_idempotent:{orderId}
		// SetNX：只在键不存在时设置，返回true表示设置成功（首次处理）
		idempotentKey := "order_cancel_idempotent:" + id
		ctx := context.TODO()
		success, err := s.rDB.SetNX(ctx, idempotentKey, "1", time.Hour*24).Result()
		if err != nil {
			log.Errorf("Failed to set idempotent key for order %d: %v", orderId, err)
			continue // Redis错误，跳过此任务，下次重试
		}

		if !success {
			// 幂等性键已存在，说明这个订单已经处理过了
			log.Warnf("Order %d already cancelled (idempotent key exists), only removing task", orderId)

			// 只删除任务，不归还库存（因为库存已经归还过了）
			err = s.redisDQRepo.RemoveTask(ctx, id)
			if err != nil {
				log.Errorf("Failed to remove duplicate task %s: %v", id, err)
				// 不return，继续处理下一个任务
			} else {
				log.Infof("Successfully removed duplicate task for order %d", orderId)
			}
			continue
		}

		// 从Redis获取payload（格式："commodityId,stock"）
		payload, err := s.rDB.Get(ctx, "dq:payload:"+id).Result()
		if err != nil {
			log.Warnf("Failed to get payload for order %d: %v", orderId, err)
			// payload不存在，删除幂等性键，下次重试
			s.rDB.Del(ctx, idempotentKey)
			continue
		}

		// 解析payload，提取商品ID和库存数量
		parts := strings.Split(payload, ",")
		if len(parts) != 2 {
			log.Warnf("Invalid payload format for order %d: %s", orderId, payload)
			// payload格式错误，删除幂等性键和任务
			s.rDB.Del(ctx, idempotentKey)
			_ = s.redisDQRepo.RemoveTask(ctx, id)
			continue
		}

		// 将payload解析成commodityId和stock
		commodityId, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Warnf("Invalid commodityId in payload: %s", parts[0])
			s.rDB.Del(ctx, idempotentKey)
			_ = s.redisDQRepo.RemoveTask(ctx, id)
			continue
		}
		stock, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Warnf("Invalid stock in payload: %s", parts[1])
			s.rDB.Del(ctx, idempotentKey)
			_ = s.redisDQRepo.RemoveTask(ctx, id)
			continue
		}

		// 归还库存到Redis（使用Lua脚本保证原子性）
		err = s.cRedisRepo.IncreaseStock(ctx, commodityId, stock)
		if err != nil {
			log.Errorf("Failed to increase stock for order %d: %v", orderId, err)
			// 归还失败，删除幂等性键，任务保留在队列中，下次重试
			s.rDB.Del(ctx, idempotentKey)
			return err
		}

		//  库存归还成功，幂等性键保留
		// 从延迟队列中移除任务
		err = s.redisDQRepo.RemoveTask(ctx, id)
		if err != nil {
			// ⚠️ 任务删除失败，但库存已归还
			// 由于幂等性键已存在，下次处理时会检测到并跳过库存归还，只删除任务
			// 不会重复归还库存！
			log.Errorf("Failed to remove task %s (stock already restored, safe to retry): %v", id, err)
			// 不return，继续处理下一个任务
			continue
		}

		log.Infof("Successfully cancelled expired order %d and restored stock %d for commodity %d", orderId, stock, commodityId)
	}
	return nil
}
