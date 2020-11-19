package controller

import (
	"Infinite_train/pkg/common/constant"
	"Infinite_train/pkg/common/utils/lock"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/service/metric_collector/system"
	"github.com/satori/go.uuid"
	"time"
)

func MetricCollectorTask(lock *lock.Mutex) {
	requestId := uuid.NewV4().String()
	golog.Infof(requestId, "start TestCon per 1 min")
	// 内部抢锁执行
	if ok := lock.TryLock(); !ok {
		golog.Warnf(requestId, "try get lock failed, skip this task")
		return
	} else {
		defer lock.Unlock()
	}

	// step 1，采集
	collectSystemModel := &system.BaseStruct{InstanceIP: "10.0.0.1", SystemMethod: constant.CollectorSystemMethodAnsible}
	systemData, _ := system.Run(collectSystemModel)
	golog.Infof(requestId, "collect system monitor data: [%f]", systemData.CPULoad)
	// step 2，计算
	time.Sleep(time.Second*120)
	// step 3, 存储
	golog.Infof(requestId, "finish metric collector task")
	return
}