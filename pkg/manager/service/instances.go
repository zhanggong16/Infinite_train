package service

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/api/restful/request"
)

type InstancesService struct {
}

func (is *InstancesService) GetInstancesWithFilter(cc *request.CommonContext, id string) (string, error) {
	requestID := cc.RequestID
	tenantID := cc.TenantID
	golog.Infof(requestID, "GetInstancesWithFilter, tenantID [%s]", tenantID)

	return id, nil
}

