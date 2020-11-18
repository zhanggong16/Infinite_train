package history

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func service() string {
	time.Sleep(time.Millisecond * 500)
	fmt.Println("zhanggong")
	return "1 Done"
}

func otherTask() {
	fmt.Println("2 working on something else")
	time.Sleep(time.Millisecond * 100)
	fmt.Println("2 Task is done")
}

func AsyncService() chan string {
	//retCh := make(chan string)
	retCh := make(chan string, 1)
	go func() {
		ret := service()
		fmt.Println("1 returned result")
		retCh <- ret
		fmt.Println("1 service exit")
	}()
	return retCh
}

func dataProducer(ch chan int, wg *sync.WaitGroup) {
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
		wg.Done()
	}()
}

func dataReceiver(ch chan int, wg *sync.WaitGroup) {
	go func() {
		for {
			if data, ok := <-ch; ok {
				fmt.Println(data)
			} else {
				break
			}
		}
		wg.Done()
	}()
}

/*func isCancelled(cancelChan chan struct{}) bool {
	select {
	case <-cancelChan:
		return true
	default:
		return false
	}
}*/

func cancle_1(cancelChan chan struct{}) {
	cancelChan <- struct{}{}
}

func cancle_2(cancelChan chan struct{}) {
	close(cancelChan)
}

func isCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func main_csp() {
	/*for i:=0;i<10;i++ {
		go func(i int){
			fmt.Println(i)
		}(i)
	}*/
	/*var mut sync.Mutex
	count := 0
	for i:=0;i<5000;i++ {
		go func() {
			defer func() {
				mut.Unlock()
			}()
			mut.Lock()
			count++
		}()
	}
	time.Sleep(time.Second*1)
	fmt.Println(count)*/

	/*var mut sync.Mutex
	var wg sync.WaitGroup
	count := 0
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				mut.Unlock()
			}()
			mut.Lock()
			count++
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(count)*/

	//cancelChan := make(chan struct{}, 0)
	ctx, cancel := context.WithCancel(context.Background())
	for i:=0;i<5;i++ {
		go func(i int, ctx context.Context) {
			for {
				if isCancelled(ctx) {
					break
				}
				time.Sleep(time.Millisecond*5)
			}
			fmt.Println(i, "Canceled")
		}(i, ctx)
	}
	cancel()
	//cancle_2(cancelChan)
	time.Sleep(time.Second*1)

	/*var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(1)
	dataProducer(ch, &wg)
	wg.Add(1)
	dataReceiver(ch, &wg)
	wg.Wait()*/

	/*select {
	case ret := <-AsyncService():
		fmt.Println(ret)
	case <-time.After(time.Millisecond*100):
		fmt.Println("time out")
	}
	*/

	/*retCh := AsyncService()
	otherTask()
	t := <-retCh
	fmt.Println(t)
	time.Sleep(time.Second*1)*/
}
