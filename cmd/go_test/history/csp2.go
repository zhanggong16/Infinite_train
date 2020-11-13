package history

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Singleton struct {
}

var singleInstance *Singleton
var once sync.Once

func GetSingletonObj() *Singleton {
	once.Do(func() {
		fmt.Println("Create Obj")
		singleInstance = new(Singleton)
	})
	return singleInstance
}

// 只运行一次就返回
func runTask(id int) string {
	time.Sleep(time.Millisecond * 10)
	return fmt.Sprintf("The result is from %d", id)
}

func FirstResponse() string {
	numOfRunner := 10
	//ch := make(chan string) 存在协程内存泄露
	ch := make(chan string, numOfRunner)
	for i := 0; i < numOfRunner; i++ {
		go func(i int) {
			ret := runTask(i)
			ch <- ret
		}(i)
	}
	return <-ch
}

func AllResponse() string {
	numOfRunner := 10
	ch := make(chan string, numOfRunner)
	for i := 0; i < numOfRunner; i++ {
		go func(i int) {
			ret := runTask(i)
			ch <- ret
		}(i)
	}


	finalRet := ""

	for j := 0; j < numOfRunner; j++ {
		finalRet += <-ch + "\n"
	}
	return finalRet
}

func t(p interface{}) {
	rt := reflect.TypeOf(p)
	fmt.Printf("%T", rt)
}

// 取消任务
/*func cancle_2(cancelChan chan struct{}) {
	close(cancelChan)
}

func isCancelled(cancelChan chan struct{}) bool {
	select {
	case <-cancelChan:
		return true
	default:
		return false
	}
}*/

func main_csp2() {
	/*cancelChan := make(chan struct{}, 0)
	for i:=0;i<5;i++ {
		go func(i int, cancelCh chan struct{}) {
			for {
				if isCancelled(cancelChan) {
					break
				}
				time.Sleep(time.Millisecond*5)
			}
			fmt.Println(i, "Canceled")
		}(i, cancelChan)
	}
	//cancel()
	cancle_2(cancelChan)
	time.Sleep(time.Second*1)*/
	/*fmt.Println("before:", runtime.NumGoroutine())
	fmt.Println(AllResponse())
	time.Sleep(time.Second * 1)
	fmt.Println("after:", runtime.NumGoroutine())*/

	/*var wg sync.WaitGroup
	for i:=0;i<10;i++{
		wg.Add(1)
		go func() {
			defer wg.Done()
			obj := GetSingletonObj()
			fmt.Println(unsafe.Pointer(obj))
		}()
	}
	wg.Wait()*/


}