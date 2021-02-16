package plugin

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	Waiting = iota
	Running
)

type WorkersError struct {
	WorkerErrors []error
}

func (we WorkersError) Error() string {
	var strs []string
	for _, err := range we.WorkerErrors {
		strs = append(strs, err.Error())
	}
	return strings.Join(strs, ";")
}

type Event struct {
	Source  string
	Content string
}

type EventReceiver interface {
	OnEvent(evt Event)
}

type Worker interface {
	Init(evtReceiver EventReceiver) error
	Start(agtCtx context.Context) error
	Stop() error
	Destroy() error
}

type Operator struct {
	workers map[string]Worker
	//evtBuf 主进程读取协程回传信息用
	evtBuf chan Event
	//cancel ctx 用于context.WithCancel
	cancel context.CancelFunc
	ctx    context.Context
	state  int
}

func NewOperator(sizeEvtBuf int) *Operator {
	agt := Operator{
		workers: map[string]Worker{},
		evtBuf:  make(chan Event, sizeEvtBuf),
		state:   Waiting,
	}

	return &agt
}

func (op *Operator) RegisterWorker(name string, worker Worker) error {
	if op.state != Waiting {
		fmt.Printf("can not take the operation in the current state")
		return errors.New("can not take the operation in the current state")
	}

	op.workers[name] = worker
	return worker.Init(op)
}

func (op *Operator) OnEvent(evt Event) {
	op.evtBuf <- evt
}

func (op *Operator) EventProcessGoroutine() {
	var evtSeg [10]Event
	for {
		for i := 0; i < 10; i++ {
			select {
			case evtSeg[i] = <-op.evtBuf:
			case <-op.ctx.Done():
				return
			}
		}
		fmt.Println(evtSeg)
	}

}

func (op *Operator) Start() error {
	if op.state != Waiting {
		fmt.Printf("can not take the operation in the current state")
		return errors.New("can not take the operation in the current state")
	}
	op.state = Running
	op.ctx, op.cancel = context.WithCancel(context.Background())
	//go op.EventProcessGoroutine()
	return op.startWorkers()
}

func (op *Operator) startWorkers() error {
	var err error
	var errs WorkersError
	var mutex sync.Mutex

	ch := make(chan error, len(op.workers))

	for name, worker := range op.workers {
		go func(name string, worker Worker, ctx context.Context) {
			defer func() {
				mutex.Unlock()
			}()
			err = worker.Start(ctx)
			mutex.Lock()
			ch <- err
			if err != nil {
				errs.WorkerErrors = append(errs.WorkerErrors,
					errors.New(name+":"+err.Error()))
			}
		}(name, worker, op.ctx)
	}

	for _, worker := range op.workers {
		select {
		case <-ch:
		case <-time.After(time.Second * 1):
			fmt.Println("timeout !!!")
			worker.Stop()
			worker.Destroy()
		}
	}

	if len(errs.WorkerErrors) == 0 {
		return nil
	}
	return errs
}
