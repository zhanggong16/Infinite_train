package system

import "Infinite_train/pkg/common/constant"

type BaseStruct struct {
	InstanceIP		string
	SystemMethod	string
}

type ResultBaseStruct struct {
	CPUUtil			float64
	CPULoad			float64
	MemoryUsed		float64
	NetworkIncoming	float64
	NetworkOutgoing	float64
	DiskIOPSRead	float64
	DiskIOPSWrite	float64
	DiskUsed		float64
	DiskTotal		float64
}

type CollectSystemBase interface {
	perform() (*ResultBaseStruct, error)
}

func Run(BaseModel *BaseStruct) (result *ResultBaseStruct, err error) {
	baseModel := BaseModel
	var cs CollectSystemBase
	switch baseModel.SystemMethod {
	case constant.CollectorSystemMethodAnsible:
		cs = &Ansible{BaseStruct: baseModel}
	case constant.CollectorSystemMethodAgent:
		cs = &Agent{BaseModel: baseModel}
	default:
		cs = &Ansible{BaseStruct: baseModel}
	}

	result, err = cs.perform()
	return
}