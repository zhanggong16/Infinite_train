package system

type Ansible struct {
	*BaseModel
}

func (ansible *Ansible) perform() (*ResultBaseModel, error) {
	return ansible.collectSystemMonitorData()
}

func (ansible *Ansible) collectSystemMonitorData() (*ResultBaseModel, error) {
	result := new(ResultBaseModel)
	result.CPULoad = 1.1
	return result, nil
}