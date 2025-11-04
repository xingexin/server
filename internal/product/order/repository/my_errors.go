package repository

import "errors"

var (
	// ErrNoTasksInQueue 队列中没有任务
	ErrNoTasksInQueue = errors.New("no tasks in delay queue")

	// ErrNoTasksDue 队列中没有到期的任务
	ErrNoTasksDue = errors.New("no tasks are due yet")

	// ErrQueueOperationFailed Redis 队列操作失败
	ErrQueueOperationFailed = errors.New("delay queue operation failed")
)
