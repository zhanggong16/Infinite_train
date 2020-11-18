package controller

import (
	"Infinite_train/pkg/common/constant"
	lock2 "Infinite_train/pkg/common/utils/lock"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/service/metric_collector/system"
	"github.com/satori/go.uuid"
	"time"
)

func MetricCollectorTask(lock *lock2.Lock) {
	requestId := uuid.NewV4().String()
	golog.Infof(requestId, "start TestCon per 1 min")
	// 内部抢锁执行
	if ok := lock.GetLock(); !ok {
		golog.Warnf(requestId, lock2.GetSchedulerLockFailed)
		return
	} else {
		defer lock.ReleaseLock()
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