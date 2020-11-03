package controller

import (
	"Infinite_train/pkg/manager/api/request"
	"Infinite_train/pkg/manager/api/view"
)

type InstancesControllerBase interface {
	GetInstances(ctx *request.CustomContext, id string) (*view.CommonGidView, *view.ResponseError)
}

var InstancesControllerImpl InstancesControllerBase = new(InstancesController)

func InitControllerLayer() {
	InstancesControllerImpl = new(InstancesController)
}