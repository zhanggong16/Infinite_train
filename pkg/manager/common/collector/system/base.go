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

func Run(baseStruct *BaseStruct) (result *ResultBaseStruct, err error) {
	var cs CollectSystemBase
	switch baseStruct.SystemMethod {
	case constant.CollectorSystemMethodAnsible:
		cs = &Ansible{BaseStruct: baseStruct}
	case constant.CollectorSystemMethodAgent:
		cs = &Agent{BaseStruct: baseStruct}
	default:
		cs = &Ansible{BaseStruct: baseStruct}
	}

	result, err = cs.perform()
	return
}