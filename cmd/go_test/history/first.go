package history

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// interface 只能是指针类型的实例,new(GoProgrammer) or &GoProgrammer{}
// 空接口可以表达任何类型，通过断言可以将空接口转换为定制类型 v, ok := p.(int)

// os.Exit 不会调用defer函数，不会输出调用栈信息

func returnMultiValues() (int, int) {
	return rand.Intn(10), rand.Intn(20)
}

func Sum(op ...int) int {
	ret := 0
	for _, v := range op {
		ret += v
	}
	return ret
}

func Clear() {
	fmt.Println("Clear resource")
}

type Employee struct {
	Id   string
	Name string
	Age  int
}

func (e *Employee) String() string {
	e.Age = 100
	return fmt.Sprintf("111 Id [%s], Name [%s], Age [%d]", e.Id, e.Name, e.Age)
}

func (e *Employee) String2() string {
	return fmt.Sprintf("222 Id [%s], Name [%s], Age [%d]", e.Id, e.Name, e.Age)
}


type IntConv func(op int) int

func timeSpent(inner IntConv) IntConv {
	return func(n int) int {
		start := time.Now()
		ret := inner(n)
		fmt.Println("time spent: ", time.Since(start).Seconds())
		return ret
	}
}

func slowFunc(op int) int {
	time.Sleep(time.Second * 1)
	return op
}

type Pet struct {
	Name string
}

func (p *Pet) Speak() {
	fmt.Println("...")
}

func (p *Pet) SpeakTo(host string) {
	p.Speak()
	fmt.Println(" ", host)
}

type Dog struct {
	Pet
}

func (d *Dog) Speak() {
	//d.p.Speak()
	fmt.Println("wangwang!")
}

/*func (d *Dog) SpeakTo(host string) {
	d.p.Speak()
	fmt.Println(" ", host)
}*/

// * &的区别
type Rect struct {
	Width  int
	Height int
}

func (r *Rect) size() int {
	return r.Height * r.Width
}

func EmptyInterface(p interface{}) {
	switch v := p.(type) {
	case int:
		fmt.Println("Integer: ", v)
	case string:
		fmt.Println("String: ", v)
	default:
		fmt.Println("Unknow Type")
	}

	/*if i, ok := p.(int); ok {
		fmt.Println("Integer: ", i)
	} else if s, ok := p.(string); ok {
		fmt.Println("String: ", s)
	} else {
		fmt.Println("Unknow Type")
	}
	return*/
}

var LessThanTwoError = errors.New("LessThanTwoError")
var LargerThanHunderdError = errors.New("LargerThanHunderdError")

func GetFib(n int) ([]int, error) {

	fibList := []int{1, 1}
	if n < 2 {
		return nil, LessThanTwoError
	}
	if n > 100 {
		return nil, LargerThanHunderdError
	}

	for i := 2; i < n; i++ {
		fibList = append(fibList, fibList[i-2]+fibList[i-1])
	}
	return fibList, nil
}

/*func init() {
	fmt.Println("init1")
}

func init() {
	fmt.Println("init2")
}*/

func listD(a string, d ...string) string {
	for _, v := range d {
		fmt.Println(a, v)
	}
	return ""
}

func main_first() {
	//EmptyInterface("10")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recovered from", err)
		}
	}()




	/*retCh := AsyncService()
	otherTask()
	t := <-retCh
	fmt.Println(t)
	time.Sleep(time.Second*1)*/

	/*s := "zhang"
	d := []string{"1","2","3"}
	listD(s,d...)*/

	/*var wg sync.WaitGroup
	for i:=0;i<10;i++ {
		wg.Add(1)
		go func (i int){
			fmt.Println(i)
			wg.Done()
		}(i)
		wg.Wait()
	}*/

	//go waitGroup
	/*var mut sync.Mutex
	var wg sync.WaitGroup
	counter := 0
	for i:=0;i <5000;i++ {
		wg.Add(1)
		go func() {
			defer func() {
				mut.Unlock()
			}()
			mut.Lock()
			counter++
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(counter)*/

	//go lock
	/*var mut sync.Mutex
	counter := 0
	for i:=0;i <5000;i++ {
		go func() {
			defer func() {
				mut.Unlock()
			}()
			mut.Lock()
			counter++
		}()
	}
	time.Sleep(1*time.Second)
	fmt.Println(counter)*/

	//go
	/*counter := 0
	for i:=0;i <5000;i++ {
		go func() {
			counter++
		}()
	}
	time.Sleep(1*time.Second)
	fmt.Println(counter)*/

	/*for i:=0;i<10;i++ {
		go func (i int) {
			fmt.Println(i)
		}(i)
		time.Sleep(1*time.Millisecond)
	}*/

	/*m := concurrent_map.CreateConcurrentMap(99)
	m.Set(concurrent_map.StrKey("key"), 10)
	fmt.Println(m.Get(concurrent_map.StrKey("key")))*/

	/*fmt.Println("start")
	panic(errors.New("wrong"))*/
	/*if s, err := GetFib(-1); err != nil {
		if err == LessThanTwoError {
			fmt.Println("less")
		}

		fmt.Println(err.Error())
	} else {
		fmt.Println(s)
	}*/

	/*r := &Rect{Width: 100, Height: 100}
	fmt.Println(r.size())*/

	/*goProg := new(GoProgrammer)
	javaProg := new(JavaProgrammer)

	WriteFirstProgram(goProg)
	WriteFirstProgram(javaProg)*/

	/*ts := timeSpent(slowFunc)
	fmt.Println(ts(10))*/

	/*var p Programmer
	p = new(GoProgrammer)
	ret := p.WriteHelloWorld()
	fmt.Println(ret)*/
}
