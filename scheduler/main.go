package main

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"math/rand"
	"sync"
	"time"
)

// https://mp.weixin.qq.com/s/Zr4_GfR2-z8GSvES4AxYDg
// 轻量级 任务调度系统

// Task 表示可执行的任务单元
type Task interface {
	ID() string
	Execute(ctx context.Context) (interface{}, error)
}

// SimpleTask 基础任务实现
type SimpleTask struct {
	id     string
	action func() (interface{}, error)
}

func (s *SimpleTask) ID() string {
	return s.id
}

func (s *SimpleTask) Execute(ctx context.Context) (interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.action()
	}
}

// **************************************Scheduler ******************************//

// Scheduler 调度器实现
type Scheduler struct {
	workerNum   int //工作协程数
	rateLimiter *rate.Limiter
	taskQueue   chan Task
	ResultChan  chan *Result
	errorChan   chan *Error
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// Result 任务执行结果
type Result struct {
	TaskId   string
	Output   interface{}
	Attempts int
}

// Error 任务执行错误
type Error struct {
	TaskId   string
	Err      error
	Attempts int
}

// NewScheduler 创建调度器
func NewScheduler(workerNum int, queueSize int) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		workerNum:  workerNum,
		taskQueue:  make(chan Task, queueSize),
		ResultChan: make(chan *Result, queueSize),
		errorChan:  make(chan *Error, queueSize),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// SetRateLimit 设置每秒最大任务数
func (s *Scheduler) SetRateLimit(perSecond int) {
	s.rateLimiter = rate.NewLimiter(rate.Limit(perSecond), perSecond)
}

// Submit 提交任务
func (s *Scheduler) Submit(task Task) {
	s.taskQueue <- task
}

// Start 启动调度器
func (s *Scheduler) Start() {
	for i := 0; i < s.workerNum; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

func (s *Scheduler) worker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case task := <-s.taskQueue:
			//速率限制
			if s.rateLimiter != nil {
				if err := s.rateLimiter.Wait(s.ctx); err != nil {
					s.errorChan <- &Error{
						TaskId:   task.ID(),
						Err:      fmt.Errorf("rate limit wait failed: %w", err),
						Attempts: 0,
					}
					continue
				}
			}
			//执行任务
			output, err := task.Execute(s.ctx)
			if err != nil {
				s.errorChan <- &Error{
					TaskId: task.ID(),
					Err:    err,
				}
			} else {
				s.ResultChan <- &Result{
					TaskId: task.ID(),
					Output: output,
				}
			}
		}
	}
}

// AdjustWorkers 动态调整worker数量：根据系统负载动态增减worker数量
func (s *Scheduler) AdjustWorkers(newNum int) {
	if newNum > s.workerNum {
		// 增加worker
		for i := s.workerNum; i < newNum; i++ {
			s.wg.Add(1)
			go s.worker()
		}
	} else if newNum < s.workerNum {
		// 减少worker (通过context取消)
		s.workerNum = newNum
	}
}

//func (s *Scheduler) SubmitAffinity(task Task, affinityKey string) {
//	// 根据affinityKey选择特定worker
//	workerIdx := hash(affinityKey) % s.workerNum
//	s.affinityQueues[workerIdx] <- task
//}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cancel()
	s.wg.Wait()
	close(s.taskQueue)
	close(s.ResultChan)
	close(s.errorChan)
}

// ************************************** 任务重试机制 ******************************//

// WithRetry 待重试的任务包装器
type WithRetry struct {
	task    Task
	max     int //最大重试时间
	backOff time.Duration
}

func (w *WithRetry) ID() string {
	return w.task.ID()
}

func (w *WithRetry) Execute(ctx context.Context) (interface{}, error) {
	var lastErr error
	for i := 0; i < w.max; i++ {
		if i > 0 {
			select {
			case <-time.After(w.backOff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		output, err := w.task.Execute(ctx)
		if err == nil {
			return output, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("after %d attempts: %w", w.max, lastErr)
}

// ************************************** 任务超时控制 ******************************//

// WithTimeout 带超时控制的任务包装器
type WithTimeout struct {
	task    Task
	timeout time.Duration
}

func (t *WithTimeout) ID() string {
	return t.task.ID()
}

func (t *WithTimeout) Execute(ctx context.Context) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()
	return t.task.Execute(ctx)
}

// ************************************** 优先级队列 ******************************//

// PriorityTask 带优先级的任务
type PriorityTask struct {
	Task     Task
	Priority int // 优先级，数字越大优先级越高
}

// PriorityQueue 优先队列实现
type PriorityQueue []*PriorityTask

func (pq *PriorityQueue) Len() int { return len(*pq) }

func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].Priority > (*pq)[j].Priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*PriorityTask))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// ************************************** 优先级队列 ******************************//

// BatchTask 批量任务处理：合并小任务为批量任务提高吞吐量
type BatchTask struct {
	tasks []Task
}

func (b *BatchTask) Execute(ctx context.Context) (interface{}, error) {
	results := make([]interface{}, len(b.tasks))
	for i, task := range b.tasks {
		out, err := task.Execute(ctx)
		if err != nil {
			return nil, err
		}
		results[i] = out
	}
	return results, nil
}

func main() {
	// 创建调度器：3个工作协程，队列大小100
	scheduler := NewScheduler(3, 100)

	// 设置每秒最多处理5个任务
	scheduler.SetRateLimit(5)

	// 启动调度器
	scheduler.Start()

	// 收集结果
	go func() {
		for result := range scheduler.ResultChan {
			fmt.Printf("Task %s completed: %v\n", result.TaskId, result.Output)
		}
	}()

	// 处理错误
	go func() {
		for err := range scheduler.errorChan {
			fmt.Printf("Task %s failed: %v (attempts: %d)\n",
				err.TaskId, err.Err, err.Attempts)
		}
	}()

	// 提交任务
	for i := 0; i < 20; i++ {
		taskID := fmt.Sprintf("task-%d", i)
		task := &SimpleTask{
			id: taskID,
			action: func() (interface{}, error) {
				// 模拟任务执行
				time.Sleep(time.Millisecond * 100)
				if rand.Intn(10) == 0 { // 10%失败率
					return nil, fmt.Errorf("random error")
				}
				return fmt.Sprintf("result of %s", taskID), nil
			},
		}

		// 包装为带重试的任务
		retryTask := &WithRetry{
			task:    task,
			max:     3,
			backOff: time.Millisecond * 200,
		}

		scheduler.Submit(retryTask)
	}

	// 等待所有任务完成
	time.Sleep(time.Second * 5)
	scheduler.Stop()
}
