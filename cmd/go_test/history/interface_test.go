package history

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

//func TestInterface(t *testing.T) {
//	var p Programmer
//	p = new(GoProgrammer)
//	res := p.WriteHelloWorld()
//	t.Log(res)
//}