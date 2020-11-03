package rpc

import (
	"log"
	"strconv"
	"testing"
	"time"
)

type Worker struct {
	Name string
}

func NewWorker() *Worker {
	return &Worker{"test"}
}

func (w *Worker) DoJob(task string, reply *string) error {
	log.Println("Worker: do job", task)
	time.Sleep(time.Second * 3)
	*reply = task
	return nil
}

func Test_RPC(t *testing.T) {
	r := NewRPCServer(":4200", 60)
	r.RegisterService(NewWorker())
	go r.ListenRPC()
	time.Sleep(1 * time.Second)
	N := 5
	c := NewRPCClient("localhost:4200", 15, 15)
	for i := 0; i < N; i++ {

		str := new(string)
		err := c.Call("Worker.DoJob", strconv.Itoa(i), str)
		if err != nil {
			t.Error(err)
		}

		if strconv.Itoa(i) != *str {
			t.Error("call ok, but execute wrong")
		}

	}
	t.Log("ok")

}
