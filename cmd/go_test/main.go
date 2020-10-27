package main

import (
	"fmt"
	"math/rand"
	"time"
)

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


type Programmer interface {
	WriteHelloWorld() string
}

type GoProgrammer struct {
	Name	string
}

func (g *GoProgrammer) WriteHelloWorld() string {
	g.Name = "zhg"
	s := fmt.Sprintf("hello world %s", g.Name)
	return s
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

func main() {
	dog := new(Dog)
	dog.SpeakTo("zhg")

	/*ts := timeSpent(slowFunc)
	fmt.Println(ts(10))*/

	/*var p Programmer
	p = new(GoProgrammer)
	ret := p.WriteHelloWorld()
	fmt.Println(ret)*/
}