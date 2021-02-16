package concurrence

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

//CSP，通信顺序进程，golang引入channel类型实现CSP思想。Don’t communicate by sharing memory, share memory by communicating.
//“communicate by sharing memory”和“share memory by communicating”是两种不同的并发处理模式。“communicate by sharing memory”是传统的并发编程处理方式，就是指，共享的数据需要用锁进行保护，goroutine 需要获取到锁，才能并发访问数据。
//“share memory by communicating”则是类似于 CSP 模型的方式，通过通信的方式，一个 goroutine 可以把数据的“所有权”交给另外一个 goroutine（虽然 Go 中没有“所有权”的概念，但是从逻辑上说，你可以把它理解为是所有权的转移）。

// 数据交流，数据传递，信号通知，任务编排，锁

//chan string          // 可以发送接收string
//chan<- struct{}      // 只能发送struct{}
//<-chan int           // 只能从chan接收int

func TestConChannel1(t *testing.T) {
	var ch = make(chan int, 10)
	for i:=0;i<10;i++ {
		select {
		case ch <- i:
		case v:=<-ch:
			t.Log(v)
		}
	}
}



// 1-4 每秒循环打印
func TestChannelForPrintNumberPreSecond(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	ch4 := make(chan int)

	go func() {
		for {
			fmt.Println("I'm goroutine 1")
			time.Sleep(1 * time.Second)
			ch2 <-1 //I'm done, you turn
			<-ch1
		}

	}()

	go func() {
		for {
			<-ch2
			fmt.Println("I'm goroutine 2")
			time.Sleep(1 * time.Second)
			ch3 <- 1
		}
	}()

	go func() {
		for {
			<-ch3
			fmt.Println("I'm goroutine 3")
			time.Sleep(1 * time.Second)
			ch4 <- 1

		}
	}()

	go func() {
		for {
			<-ch4
			fmt.Println("I'm goroutine 4")
			time.Sleep(1 * time.Second)
			ch1 <-1
		}
	}()
	select {}
}

type Token struct{}

func newWorker(id int, ch chan Token, nextCh chan Token) {
	for {
		token := <-ch         // 取得令牌
		fmt.Println((id + 1)) // id从1开始
		time.Sleep(time.Second)
		nextCh <- token
	}
}
func TestChannelForPrintNumberPreSecond2(t *testing.T) {
	chs := []chan Token{make(chan Token), make(chan Token), make(chan Token), make(chan Token)}

	// 创建4个worker
	for i := 0; i < 4; i++ {
		go newWorker(i, chs[i], chs[(i+1)%4])
	}

	// 首先把令牌交给第一个worker
	chs[0] <- struct{}{}

	select {}
}


// 1-4顺序打印
func TestChannelForSQPrintNumber(t *testing.T) {
	var wg sync.WaitGroup
	for i:=0;i<4;i++{
		wg.Add(1)
		go func(i int) {
			time.Sleep(time.Millisecond*time.Duration(i))
			fmt.Println(i+1)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// 1-4顺序打印
func TestChannelForSQPrintNumber2(t *testing.T) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	count := 1
	for i:=0;i<4;i++{
		wg.Add(1)
		go func(i int) {
			mu.Lock()
			fmt.Println(count)
			count += 1
			mu.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
}