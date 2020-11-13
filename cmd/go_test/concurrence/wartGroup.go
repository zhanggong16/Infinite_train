package concurrence

import (
	"fmt"
	"sync"
	"time"
)

// waitGroup，现在有一个 goroutine A 在检查点（checkpoint）等待一组 goroutine 全部完成，如果在执行任务的这些 goroutine 还没全部完成，那么 goroutine A 就会阻塞在检查点，直到所有 goroutine 都完成后才能继续执行。

// Add，用来设置 WaitGroup 的计数值；放在go外面，与Done成对出现。
// Done，用来将 WaitGroup 的计数值减 1，其实就是调用了 Add(-1)；
// Wait，调用这个方法的 goroutine 会一直阻塞，直到 WaitGroup 的计数值变为 0


type Counter2 struct {
	mu 		sync.Mutex
	counter	uint64
}

func (c *Counter2) insr() {
	c.mu.Lock()
	c.counter++
	c.mu.Unlock()
}

func (c *Counter2) Count() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counter
}

func main() {
	c := new(Counter2)
	c.counter = 0

	var wg sync.WaitGroup
	wg.Add(10)
	for i:=0;i<10;i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Second)
			c.insr()
		}()
	}
	wg.Wait()
	fmt.Println(c.Count())
}