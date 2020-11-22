package collector

import (
	"Infinite_train/pkg/common/constant"
	"Infinite_train/pkg/manager/context"
)

// 操作系统监控项入口
type SystemBaseStruct struct {
	*context.SystemMetricContext
}

type SystemBase interface {
	SysPerform() (resultArr []map[string]interface{}, err error)
}

func SystemCollector(systemMetricContext *context.SystemMetricContext) ([]map[string]interface{}, error) {
	systemBaseStruct := &SystemBaseStruct{systemMetricContext}
	var sb SystemBase
	switch systemMetricContext.Collector {
	case constant.MetricCollectorMethodMySQLAnsible:
		sb = &Ansible{SystemBaseStruct: systemBaseStruct}
	case constant.MetricCollectorMethodSystemAgent:
		sb = &Agent{SystemBaseStruct: systemBaseStruct}
	default:
		sb = &Ansible{SystemBaseStruct: systemBaseStruct}
	}

	result, err := sb.SysPerform()
	return result, err
}

// 数据库监控项入口
type DatabaseBaseStruct struct {
	*context.DatabaseMetricContext
}

type DatabaseBase interface {
	DbPerform() (resultArr []map[string]interface{}, err error)
}

func DatabaseCollector(databaseMetricContext *context.DatabaseMetricContext) ([]map[string]interface{}, error) {
	databaseBaseStruct := &DatabaseBaseStruct{databaseMetricContext}
	var db DatabaseBase
	switch databaseMetricContext.Collector {
	case constant.MetricCollectorMethodMySQLConnect:
		db = &MySQLConnect{DatabaseBaseStruct: databaseBaseStruct}
	default:
		db = &MySQLConnect{DatabaseBaseStruct: databaseBaseStruct}
	}

	result, err := db.DbPerform()
	return result, err
}