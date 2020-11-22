package service

import (
	"Infinite_train/pkg/manager/api/restful/request"
	"Infinite_train/pkg/manager/context"
)

type InstancesServiceBase interface {
	GetInstancesWithFilter(cc *request.CommonContext, id string) (string, error)
}

type MetricServiceBase interface {
	AsyncSystemMetricsUpload(requestID string, systemMetricArrContext []*context.SystemMetricContext)
	AsyncDatabaseMetricsUpload(requestID string, databaseMetricArrContext []*context.DatabaseMetricContext)
}


var InstancesServiceImpl InstancesServiceBase = new(InstancesService)
var MetricServiceImpl MetricServiceBase = new(MetricService)

func InitServiceLayer() {
	InstancesServiceImpl = new(InstancesService)
	MetricServiceImpl = new(MetricService)
}