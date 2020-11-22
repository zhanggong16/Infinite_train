package controller

import (
	"Infinite_train/pkg/common/constant"
	commonContext "Infinite_train/pkg/common/context"
	"Infinite_train/pkg/common/utils/lock"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/context"
	"Infinite_train/pkg/manager/service"
	"github.com/satori/go.uuid"
)

func MetricCollectorTask(lock *lock.Mutex) {
	requestID := uuid.NewV4().String()
	systemRequestID := requestID + "_sys"
	databaseRequestID := requestID + "_db"
	// 内部抢锁执行
	if ok := lock.TryLock(); !ok {
		golog.Warnf(requestID, "try get local lock failed, skip this task")
		return
	} else {
		defer lock.Unlock()
	}
	// 获取目标和属性
	var systemMetricArrContext []*context.SystemMetricContext
	systemMetricContext1 := &context.SystemMetricContext{
		InstanceIP: "10.0.0.1",
		SystemType: constant.MetricTypeLinux,
		Collector:  constant.MetricCollectorMethodMySQLAnsible,
	}
	systemMetricContext2 := &context.SystemMetricContext{
		InstanceIP: "10.0.0.2",
		SystemType: constant.MetricTypeLinux,
		Collector:  constant.MetricCollectorMethodMySQLAnsible,
	}
	systemMetricArrContext = append(systemMetricArrContext, systemMetricContext1, systemMetricContext2)

	var databaseMetricArrContext []*context.DatabaseMetricContext
	databaseInfo := &commonContext.DatabaseStruct{
		Host: "10.0.0.1",
		User: "zhg",
		Pwd:  "zhg",
		Port: 3306,
	}
	databaseMetricContext1 := &context.DatabaseMetricContext{
		InstanceIP:     "10.0.0.1",
		DatabaseType:   constant.MetricTypeMySQL,
		Collector:      constant.MetricCollectorMethodMySQLConnect,
		DatabaseStruct:	*databaseInfo,
	}
	databaseMetricContext2 := &context.DatabaseMetricContext{
		InstanceIP:     "10.0.0.2",
		DatabaseType:   constant.MetricTypeMySQL,
		Collector:      constant.MetricCollectorMethodMySQLConnect,
		DatabaseStruct:	*databaseInfo,
	}
	databaseMetricArrContext = append(databaseMetricArrContext, databaseMetricContext1, databaseMetricContext2)
	// call service
	service.MetricServiceImpl.AsyncSystemMetricsUpload(systemRequestID, systemMetricArrContext)
	service.MetricServiceImpl.AsyncDatabaseMetricsUpload(databaseRequestID, databaseMetricArrContext)

	return
}