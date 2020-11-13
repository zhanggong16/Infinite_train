package controller

import (
	"Infinite_train/pkg/common/constant"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/api/restful/request"
	"Infinite_train/pkg/manager/api/view"
	"Infinite_train/pkg/manager/service"
)

type InstancesController struct {
}

func (ic *InstancesController) GetInstances(cc *request.CommonContext, instanceID string) (*view.CommonGidView, *view.ResponseError) {
	requestID := cc.RequestID
	tenantID := cc.TenantID
	golog.Infof(requestID, "GetInstances, tenantID [%s]", tenantID)
	// call service
	resID, err := service.InstancesServiceImpl.GetInstancesWithFilter(cc, instanceID)
	if err != nil {
		errorResp := view.NewResponseError(constant.SelectDBErrorCode, requestID, err.Error())
		golog.Errorf(requestID, "tenant_id: %s, error message: %s", tenantID, errorResp.Error())
		return nil, errorResp
	}

	resView := &view.CommonGidView{GID: resID}
	return resView, nil
}
func (ic *InstancesController) ChangeInstanceName(mc *request.ManagerCommonContext, newInstanceName string) *view.ResponseError {
	requestID := mc.CommonContext.RequestID
	tenantID := mc.CommonContext.TenantID
	golog.Infof(requestID, "ChangeInstanceName, tenantID [%s], new instance name [%s]", tenantID, newInstanceName)

	return nil
}