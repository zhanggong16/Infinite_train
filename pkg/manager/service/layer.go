package service

import "Infinite_train/pkg/manager/api/request"

type InstancesServiceBase interface {
	GetInstancesWithFilter(ctx *request.CustomContext, id string) (string, error)
}

var InstancesServiceImpl InstancesServiceBase = new(InstancesService)

func InitServiceLayer() {
	InstancesServiceImpl = new(InstancesService)
}