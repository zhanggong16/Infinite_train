package ut

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSq(t *testing.T) {
	inputs := [...]int{1, 2, 3}
	expected := [...]int{1, 4, 9}
	for i := 0; i < len(inputs); i++ {
		ret := sq(i)
		if ret == expected[i] {
			t.Errorf("input is %d, expect is %d, the actual %d", inputs[i], expected[i], ret)
		}
	}
}

func TestSqAssert(t *testing.T) {
	inputs := [...]int{1, 2, 3}
	expected := [...]int{1, 4, 9}
	for i := 0; i < len(inputs); i++ {
		ret := sq(i)
		assert.Equal(t, expected[i], ret)
		/*if ret == expected[i] {
			t.Errorf("input is %d, expect is %d, the actual %d", inputs[i], expected[i], ret)
		}*/
	}
}

// error，测试失败，该测试继续
// fatal，测试失败，测试终止

// go test -v -cover

func TestError(t *testing.T) {
	fmt.Println("start")
	t.Error("error")
	fmt.Println("end")
}

func TestFatal(t *testing.T) {
	fmt.Println("start")
	t.Fatal("error")
	fmt.Println("end")
}

// benchmark
// go test -bench=. -benchmem
// windows	go test -bench="."

func TestConcatStringByAdd(t *testing.T) {
	assert := assert.New(t)
	elems := []string{"1", "2", "3", "4", "5"}
	ret := ""
	for _, v := range elems {
		ret += v
	}
	assert.Equal("12345", ret)
}

func TestConcatStringByBytesBuffer(t *testing.T) {
	assert := assert.New(t)
	elems := []string{"1", "2", "3", "4", "5"}
	var buf bytes.Buffer
	for _, v := range elems {
		buf.WriteString(v)
	}
	assert.Equal("12345", buf.String())
}

func BenchmarkConcatStringByAdd(b *testing.B) {
	elems := []string{"1", "2", "3", "4", "5"}
	b.ResetTimer()
	for i:=0;i<b.N;i++ {
		ret := ""
		for _, v := range elems {
			ret += v
		}
	}
	b.StopTimer()
}

func BenchmarkConcatStringByBytesBuffer(b *testing.B) {
	elems := []string{"1", "2", "3", "4", "5"}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		var buf bytes.Buffer
		for _, v := range elems {
			buf.WriteString(v)
		}
	}

	b.StopTimer()
}