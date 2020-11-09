package system

type Ansible struct {
	*BaseStruct
}

func (ansible *Ansible) perform() (*ResultBaseStruct, error) {
	return ansible.collectSystemMonitorData()
}

func (ansible *Ansible) collectSystemMonitorData() (*ResultBaseStruct, error) {
	result := new(ResultBaseStruct)
	result.CPULoad = 1.1
	return result, nil
}