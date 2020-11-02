package controller

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"github.com/satori/go.uuid"
)

func TestCon() {
	requestId := uuid.NewV4().String()
	golog.Infof(requestId, "start TestCon")
}