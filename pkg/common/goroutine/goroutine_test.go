package goroutine

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

type boolType bool

func runTask(args interface{}) string {
	data := args.(map[string]interface{})
	id := data["id"].(int)
	result := fmt.Sprintf("The result is from %d", id)
	time.Sleep(time.Second * time.Duration(id))
	result += " finish"
	return result
}

func TestAllResponse(t *testing.T) {
	t.Log("before:", runtime.NumGoroutine())
	//t.Log(AllResponse())
	c := new(CollectSystemData)
	c.Tag = "test"
	c.AllResponse(runTask, map[string]interface{}{"id": 1})
	t.Log(c.Result)
	time.Sleep(time.Second * 1)
	t.Log("after:", runtime.NumGoroutine())
}