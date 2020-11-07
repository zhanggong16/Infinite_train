package controller

import (
	"Infinite_train/pkg/manager/api/request"
	"Infinite_train/pkg/manager/api/view"
)

type InstancesControllerBase interface {
	GetInstances(cc *request.CommonContext, instanceID string) (*view.CommonGidView, *view.ResponseError)
	ChangeInstanceName(mc *request.ManagerCommonContext, newInstanceName string) *view.ResponseError
}

var InstancesControllerImpl InstancesControllerBase = new(InstancesController)

func InitControllerLayer() {
	InstancesControllerImpl = new(InstancesController)
}