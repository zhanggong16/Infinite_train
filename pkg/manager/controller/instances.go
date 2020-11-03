package controller

import (
	"Infinite_train/pkg/common/constant"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/api/request"
	"Infinite_train/pkg/manager/api/view"
	"Infinite_train/pkg/manager/service"
)

type InstancesController struct {
}

func (ic *InstancesController) GetInstances(ctx *request.CustomContext, id string) (*view.CommonGidView, *view.ResponseError) {
	requestID := ctx.CommonContext.RequestID
	tenantID := ctx.CommonContext.TenantID
	golog.Infof(requestID, "GetInstances, tenantID [%s]", tenantID)
	// call service
	resID, err := service.InstancesServiceImpl.GetInstancesWithFilter(ctx, id)
	if err != nil {
		errorResp := view.NewResponseError(constant.SelectDBErrorCode, requestID, err.Error())
		golog.Errorf(requestID, "tenant_id: %s, error message: %s", tenantID, errorResp.Error())
		return nil, errorResp
	}

	resView := &view.CommonGidView{Gid: resID}
	return resView, nil
}