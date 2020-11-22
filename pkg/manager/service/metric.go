package service

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/context"
	"Infinite_train/pkg/manager/service/metric_collector/collector"
)

type MetricService struct {
}

func (ms *MetricService) AsyncSystemMetricsUpload(requestID string, systemMetricArrContext []*context.SystemMetricContext) {
	for _, systemMetric := range systemMetricArrContext {
		go func(systemMetric *context.SystemMetricContext) {
			// 获取metric
				resultArr, _ := collector.SystemCollector(systemMetric)
				golog.Infof(requestID, "systemMetric result %s", resultArr)
			// upload storage

		}(systemMetric)
	}

	return
}

func (ms *MetricService) AsyncDatabaseMetricsUpload(requestID string, databaseMetricArrContext []*context.DatabaseMetricContext) {
	for _, databaseMetric := range databaseMetricArrContext {
		go func(databaseMetric *context.DatabaseMetricContext) {
			// 获取metric
			resultArr, _ := collector.DatabaseCollector(databaseMetric)
			golog.Infof(requestID, "databaseMetric result %s", resultArr)
		}(databaseMetric)
	}
	return
}