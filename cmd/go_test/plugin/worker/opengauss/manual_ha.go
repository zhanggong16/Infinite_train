package opengauss

import (
	"Infinite_train/cmd/go_test/plugin"
	"context"
	"errors"
	"fmt"
	"time"
)

type ManualHa struct {
	evtReceiver plugin.EventReceiver
	agtCtx      context.Context

	stopChan    chan struct{}
	name        string
	content     string
}

func NewManualHa(name string, content string) *ManualHa {
	return &ManualHa{
		stopChan: make(chan struct{}),
		name:     name,
		content:  content,
	}
}

func (m *ManualHa) Init(evtReceiver plugin.EventReceiver) error {
	fmt.Println("initialize worker", m.name)
	m.evtReceiver = evtReceiver
	return nil
}

func (m *ManualHa) Start(agtCtx context.Context) error {
	fmt.Println("start worker", m.name)
	select {
	case <-agtCtx.Done():
		m.stopChan <- struct{}{}
		break
	default:
		time.Sleep(time.Second * 3)
		fmt.Println("work finished")
		//m.evtReceiver.OnEvent(plugin.Event{m.name, m.content})
	}
	return nil
}

func (m *ManualHa) Stop() error {
	fmt.Println("stop work", m.name)
	select {
	case <-m.stopChan:
		return nil
	case <-time.After(time.Second * 1):
		return errors.New("failed to stop for timeout")
	}
}

func (m *ManualHa) Destroy() error {
	fmt.Println(m.name, "released resources.")
	return nil
}
