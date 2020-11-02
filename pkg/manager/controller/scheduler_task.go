package controller

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"github.com/satori/go.uuid"
	"time"
)

func TestCon() {
	requestId := uuid.NewV4().String()
	golog.Infof(requestId, "start TestCon")
	time.Sleep(120 * time.Second)
}

func TestScheduler() {
	requestId := uuid.NewV4().String()
	golog.Infof(requestId, "start TestScheduler")
}