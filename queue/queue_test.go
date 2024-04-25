package queue

import (
	"fmt"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	que := New()
	for i := 1; i <= 10; i++ {
		que.Push(i)
	}

	go func() {
		for {
			fmt.Println(que.Pop().(int))
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			fmt.Println(que.Pop().(int))
			time.Sleep(1 * time.Second)
		}
	}()

	for i := 11; i <= 20; i++ {
		que.Push(i)
	}

	go func() {
		for {
			fmt.Println(que.Pop().(int))
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			fmt.Println(que.Pop().(int))
			time.Sleep(1 * time.Second)
		}
	}()

	que.Wait()
	fmt.Println("Down")
}

func TestClose(t *testing.T) {
	que := New()
	for i := 0; i < 10; i++ {
		que.Push(i)
	}
	go func() {
		for {
			v := que.Pop()
			if v != nil {
				fmt.Println(v.(int))
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			v := que.Pop()
			if v != nil {
				fmt.Println(v.(int))
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// que.Close()
	que.Wait()
	fmt.Println("Down")
}
