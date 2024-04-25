package queue

import (
	"runtime"
	"sync"
	"sync/atomic"

	"gopkg.in/eapache/queue.v1"
)

type Queue struct {
	sync.Mutex
	cond   *sync.Cond
	buffer *queue.Queue // 队列内容
	count  int32        // 队列长度
	close  bool         // 是否关闭
}

// New 创建Queue
func New() *Queue {
	q := &Queue{
		buffer: queue.New(),
	}
	q.cond = sync.NewCond(&q.Mutex)
	return q
}

// Push 放入队列 （非阻塞模式）
func (q *Queue) Push(v interface{}) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if !q.close {
		q.buffer.Add(v)
		atomic.AddInt32(&q.count, 1)
		q.cond.Signal() // 通知下去队列有新值
	}
}

// Pop 取出队列 （阻塞模式）
func (q *Queue) Pop() (v interface{}) {
	c := q.cond

	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	for q.Len() == 0 && !q.close {
		c.Wait()
	}

	if q.close {
		return
	}

	if q.Len() > 0 {
		buffer := q.buffer
		v = buffer.Peek()
		buffer.Remove()
		atomic.AddInt32(&q.count, -1)
	}
	return
}

// TryPop 试着取出队列（非阻塞模式）返回ok == false 表示空
func (q *Queue) TryPop() (v interface{}, ok bool) {
	buffer := q.buffer

	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if q.Len() > 0 {
		v = buffer.Peek()
		buffer.Remove()
		atomic.AddInt32(&q.count, -1)
		ok = true
	} else if q.close {
		ok = true
	}

	return
}

// 获取队列长度
func (q *Queue) Len() int {
	return (int)(atomic.LoadInt32(&q.count))
}

func (q *Queue) Wait() {
	for {
		if q.close || q.Len() == 0 {
			break
		}
		runtime.Gosched() // 出让时间片
	}
}

// 关闭队列
func (q *Queue) Close() {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if !q.close {
		q.close = true
		atomic.StoreInt32(&q.count, 0)
		q.cond.Broadcast() // 广播
	}
}

// 判断是否已经关闭
func (q *Queue) IsClose() bool {
	return q.close
}
