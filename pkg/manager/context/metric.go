package context

import "Infinite_train/pkg/common/context"

//GetMetricContext is used to delivery context for get system metric
type SystemMetricContext struct {
	RequestID		string
	InstanceIP		string
	SystemType		string
	Collector		string	// 获取操作系统指标的方式，agent or ansible
}

//GetMetricContext is used to delivery context for get db metric
type DatabaseMetricContext struct {
	RequestID		string
	InstanceIP		string
	DatabaseType	string
	Collector		string
	DatabaseStruct	context.DatabaseStruct
}