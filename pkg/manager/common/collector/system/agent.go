package system

type Agent struct {
	*BaseModel
}

func (agent *Agent) perform() (*ResultBaseModel, error) {
	return agent.collectSystemMonitorData()
}

func (agent *Agent) collectSystemMonitorData() (*ResultBaseModel, error) {
	result := new(ResultBaseModel)
	result.CPULoad = 2.1
	return result, nil
}