package main

import (
	"fmt"
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

func main() {
	goPro := new(GoProgrammer)
	javaPro := new(JavaProgrammer)
	WriteFirstProgram(goPro)
	WriteFirstProgram(javaPro)
}
