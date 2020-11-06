package controller

import (
	"Infinite_train/pkg/common/constant"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/common/collector/system"
	"github.com/satori/go.uuid"
)

func TestCon() {
	requestId := uuid.NewV4().String()
	golog.Infof(requestId, "start TestCon per 1 min")
	// step 1，采集
	collectSystemModel := &system.BaseModel{InstanceIP: "10.0.0.1", SystemMethod: constant.CollectorSystemMethodAnsible}
	systemData, _ := system.Run(collectSystemModel)
	golog.Infof(requestId, "collect system monitor data: [%s]", systemData.CPULoad)

	// step 2，计算

	// step 3, 存储
}