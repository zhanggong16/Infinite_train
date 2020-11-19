package concurrence

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

//chan string          // 可以发送接收string
//chan<- struct{}      // 只能发送struct{}
//<-chan int           // 只能从chan接收int

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

	//首先把令牌交给第一个worker
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