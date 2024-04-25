package workpool

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/objectzsp/public/queue"
)

type TaskHandler func() error

type WorkPool struct {
	closed    int32            // 标记工作池是否已关闭。
	task      chan TaskHandler // 通过通道传输任务对象。
	queue     *queue.Queue     // 任务队列。
	isQueTask int32            // 标记是否队列取出任务。
	timeout   time.Duration    //
	wg        sync.WaitGroup   //
	errChan   chan error       // 通过通道返回错误信息。
}

// 注册工作池，并设置最大并发数。
func New(max int) *WorkPool {
	// 线程数必须大于1，否则线程池就没有任何作用。
	if max < 1 {
		max = 1
	}

	wp := &WorkPool{
		errChan: make(chan error, 1),
		task:    make(chan TaskHandler, 2*max),
		queue:   queue.New(),
	}
	go Loop(max, wp)
	return wp
}

// 添加任务到工作池，并立即返回。
func (wp *WorkPool) Do(th TaskHandler) {
	if wp.IsClosed() {
		// 如果工作池已关闭，立即返回。
		return
	}
	wp.queue.Push(th)
}

// 添加任务到工作池，并等待执行完成之后再返回。
func (wp *WorkPool) DoWait(th TaskHandler) {
	if wp.IsClosed() {
		// 如果工作池已关闭，立即返回。
		return
	}

	doneChan := make(chan struct{})
	wp.queue.Push(TaskHandler(func() error {
		defer close(doneChan)
		return th()
	}))
	<-doneChan
}

// 等待工作线程执行结束
func (wp *WorkPool) Wait() error {
	wp.queue.Wait()
	wp.queue.Close()
	WaitTask(wp)
	close(wp.task)
	wp.wg.Wait()
	select {
	case err := <-wp.errChan:
		return err
	default:
		return nil
	}
}

// 判断是否完成 (非阻塞)
func (wp *WorkPool) IsDone() bool {
	if wp == nil || wp.task == nil {
		return true
	}
	// 返回队列和任务通道是否为nil
	return wp.queue.Len() == 0 && len(wp.task) == 0
}

// 是否已经关闭
func (wp *WorkPool) IsClosed() bool {
	// if atomic.LoadInt32(&wp.closed) == 1 {
	// 	return true
	// }
	// return false
	return atomic.LoadInt32(&wp.closed) == 1
}

// 设置超时时间
func (wp *WorkPool) SetTimeOut(timeout time.Duration) {
	wp.timeout = timeout
}
