package workpool

import (
	"context"
	"runtime"
	"sync/atomic"
)

func StartQueue(wp *WorkPool) {
	wp.isQueTask = 1
	for {
		tmp := wp.queue.Pop()
		if wp.IsClosed() {
			wp.queue.Close()
			break
		}
		if tmp != nil {
			fn := tmp.(TaskHandler)
			if fn != nil {
				wp.task <- fn
			}
		} else {
			break
		}
	}
	atomic.StoreInt32(&wp.isQueTask, 0)
}

func WaitTask(wp *WorkPool) {
	for {
		runtime.Gosched() //用于让出CPU时间片
		if wp.IsDone() {
			if atomic.LoadInt32(&wp.isQueTask) == 0 {
				break
			}
		}
	}
}

func Loop(max int, wp *WorkPool) {
	// 启动任务队列
	go StartQueue(wp)
	// 设定工具人数量
	wp.wg.Add(max)
	// 根据max创建工具人协程，让工具人工作。
	for i := 0; i < max; i++ {
		go func() {
			defer wp.wg.Done()
			for wt := range wp.task {
				if wt == nil || atomic.LoadInt32(&wp.closed) == 1 {
					continue
				}

				closed := make(chan struct{}, 1)
				if wp.timeout > 0 {
					ct, cancel := context.WithTimeout(context.Background(), wp.timeout)
					go func() {
						select {
						case <-ct.Done():
							wp.errChan <- ct.Err()
							atomic.StoreInt32(&wp.closed, 1)
							cancel()
						case <-closed:
						}
					}()
				}

				err := wt()
				close(closed)
				if err != nil {
					select {
					case wp.errChan <- err:
						atomic.StoreInt32(&wp.closed, 1)
					default:
					}
				}
			}
		}()
	}
}
