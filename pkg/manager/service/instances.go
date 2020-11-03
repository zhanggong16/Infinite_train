package service

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/api/request"
)

type InstancesService struct {
}

func (is *InstancesService) GetInstancesWithFilter(ctx *request.CustomContext, id string) (string, error) {
	requestID := ctx.CommonContext.RequestID
	tenantID := ctx.CommonContext.TenantID
	golog.Infof(requestID, "GetInstancesWithFilter, tenantID [%s]", tenantID)

	return id, nil
}