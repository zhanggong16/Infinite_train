package system

type Agent struct {
	*BaseStruct
}

func (agent *Agent) perform() (*ResultBaseStruct, error) {
	return agent.collectSystemMonitorData()
}

func (agent *Agent) collectSystemMonitorData() (*ResultBaseStruct, error) {
	result := new(ResultBaseStruct)
	result.CPULoad = 2.1
	return result, nil
}