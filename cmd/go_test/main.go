package main

import (
	"Infinite_train/pkg/common/utils/linux"
	"fmt"
	"time"
)

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


func main() {

	s, err := linux.ExecWithTimeout("1111","ping www.baidu.com", 5)
	fmt.Printf("result %s, error %s", s, err.Error())

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

	/*goPro := new(GoProgrammer)
	javaPro := new(JavaProgrammer)
	WriteFirstProgram(goPro)
	WriteFirstProgram(javaPro)*/
}
