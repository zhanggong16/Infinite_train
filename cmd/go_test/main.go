package main

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
	for _, v :=range op {
		ret += v
	}
	return ret
}


func Clear() {
	fmt.Println("Clear resource")
}

type Employee struct {
	Id string
	Name	string
	Age	int
}

func (e *Employee) String() string {
	e.Age = 100
	return fmt.Sprintf("111 Id [%s], Name [%s], Age [%d]", e.Id, e.Name, e.Age)
}

func (e *Employee) String2() string {
	return fmt.Sprintf("222 Id [%s], Name [%s], Age [%d]", e.Id, e.Name, e.Age)
}

type Code string

type Programmer interface {
	WriteHelloWorld() Code
}

type GoProgrammer struct {
}

func (g *GoProgrammer) WriteHelloWorld() Code {
	return "golang: hello world"
}

type JavaProgrammer struct {
}

func (j *JavaProgrammer) WriteHelloWorld() Code {
	return "java: hello world"
}


func WriteFirstProgram(p Programmer) {
	fmt.Printf("%T, %v\n", p, p.WriteHelloWorld())
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
	Sex	string
}

func (d *Dog) Speak() {
	fmt.Println("Wang")
}

// * &的区别
type Rect struct {
	Width 	int
	Height	int
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


	fibList := []int{1,1}
	if n < 2 {
		return nil, LessThanTwoError
	}
	if  n > 100 {
		return nil, LargerThanHunderdError
	}

	for i:=2; i<n; i++ {
		fibList = append(fibList, fibList[i-2] + fibList[i-1])
	}
	return fibList, nil
}





func main() {
	//EmptyInterface("10")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recovered from", err)
		}
	}()

	fmt.Println("start")
	panic(errors.New("wrong"))
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