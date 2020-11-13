package service

import "Infinite_train/pkg/manager/api/restful/request"

type InstancesServiceBase interface {
	GetInstancesWithFilter(cc *request.CommonContext, id string) (string, error)
}

var InstancesServiceImpl InstancesServiceBase = new(InstancesService)

func InitServiceLayer() {
	InstancesServiceImpl = new(InstancesService)
}