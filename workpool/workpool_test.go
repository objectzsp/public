package workpool

import (
	"fmt"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	wp := New(10)
	for i := 0; i < 10; i++ {
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(1 * time.Second)
			}
			return nil
		})
	}
	wp.Wait()
	fmt.Println("down")
}
