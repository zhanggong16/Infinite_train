package collector

type Ansible struct {
	*SystemBaseStruct
}

func (an *Ansible) SysPerform() ([]map[string]interface{}, error) {
	var metric map[string]interface{}
	metric["system"] = an.InstanceIP
	var results []map[string]interface{}
	results = append(results, metric)

	return results, nil
}