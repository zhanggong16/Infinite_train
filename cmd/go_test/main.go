package main

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

func main() {
	t(1)

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
