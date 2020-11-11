package concurrence

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Synchronization_primitives
// 共享资源，并发读写资源，会出现数据竞争，需要mutex来保护
// 任务编排，需要goroutine按照一定规律执行。使用waitGroup来实现
// 消息传递，不同goroutine之间的线程安全的数据交流，channel实现

// 在真实的实践中，我们使用互斥锁的时候，很少在一个方法中单独申请锁，而在另外一个方法中单独释放锁，一般都会在同一个方法中获取锁和释放锁。
// Unlock 方法可以被任意的 goroutine 调用释放锁，即使是没持有这个互斥锁的 goroutine，也可以进行这个操作。这是因为，Mutex 本身并没有包含持有这把锁的 goroutine 的信息，所以，Unlock 也不会对此进行检查。Mutex 的这个设计一直保持至今。

// 保证 Lock/Unlock 成对出现，尽可能采用 defer mutex.Unlock 的方式，把它们成对、紧凑地写在一起。


const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexStarving
	mutexWaiterShift = iota
)

type Counter struct {
	mu 		sync.Mutex
	Count 	uint64
}

// 解析state状态码
func (c *Counter) mutexCount() int {
	// 获取state字段的值
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&c.mu)))
	v = v >> mutexWaiterShift //得到等待者的数值
	v = v + (v & mutexLocked) //再加上锁持有者的数量，0或者1
	return int(v)
}

// 锁是否被持有
func (c *Counter) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&c.mu)))
	return state&mutexLocked == mutexLocked
}

// 是否有等待者被唤醒
func (c *Counter) IsWoken() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&c.mu)))
	return state&mutexWoken == mutexWoken
}

// 锁是否处于饥饿状态
func (c *Counter) IsStarving() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&c.mu)))
	return state&mutexStarving == mutexStarving
}


// 业务逻辑代码
func (c *Counter) incr() {
	c.mu.Lock()
	c.Count++
	c.mu.Unlock()
}

func (c *Counter) count() uint64 {
	defer c.mu.Unlock()
	c.mu.Lock()
	return c.Count
}

func main() {

	var counter Counter
	var wg sync.WaitGroup
	wg.Add(10)

	for i:=0;i<10;i++ {
		//wg.Add(1)
		go func() {
			defer wg.Done()
			for j :=0;j<100000;j++ {
				counter.incr()
			}
		}()
	}
	wg.Wait()
	fmt.Println(counter.count())
	fmt.Printf("waitings: %d, isLocked: %t, woken: %t, starving: %t\n", counter.mutexCount(), counter.IsLocked(), counter.IsWoken(), counter.IsStarving())

	/*var c Counter
	for i := 0; i < 1000; i++ { // 启动1000个goroutine
		go func() {
			c.mu.Lock()
			time.Sleep(time.Second)
			c.mu.Unlock()
		}()
	}
	fmt.Printf("waitings: %d, isLocked: %t, woken: %t, starving: %t\n", c.mutexCount(), c.IsLocked(), c.IsWoken(), c.IsStarving())*/
}

// 安全的slice

type SliceQueue struct {
	data []interface{}
	mu   sync.Mutex
}

func NewSliceQueue(n int) (q *SliceQueue) {
	return &SliceQueue{data: make([]interface{}, 0, n)}
}

// Enqueue 把值放在队尾
func (q *SliceQueue) Enqueue(v interface{}) {
	q.mu.Lock()
	q.data = append(q.data, v)
	q.mu.Unlock()
}

// Dequeue 移去队头并返回
func (q *SliceQueue) Dequeue() interface{} {
	q.mu.Lock()
	if len(q.data) == 0 {
		q.mu.Unlock()
		return nil
	}
	v := q.data[0]
	q.data = q.data[1:]
	q.mu.Unlock()
	return v
}