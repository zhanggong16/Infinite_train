package system

import "Infinite_train/pkg/common/constant"

type BaseModel struct {
	InstanceIP		string
	SystemMethod	string
}

type ResultBaseModel struct {
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
	perform() (*ResultBaseModel, error)
}

func Run(BaseModel *BaseModel) (result *ResultBaseModel, err error) {
	baseModel := BaseModel
	var cs CollectSystemBase
	switch baseModel.SystemMethod {
	case constant.CollectorSystemMethodAnsible:
		cs = &Ansible{BaseModel: baseModel}
	case constant.CollectorSystemMethodAgent:
		cs = &Agent{BaseModel: baseModel}
	default:
		cs = &Ansible{BaseModel: baseModel}
	}

	result, err = cs.perform()
	return
}