package collector

type Agent struct {
	*SystemBaseStruct
}

func (ag *Agent) SysPerform() ([]map[string]interface{}, error) {
	var metric map[string]interface{}
	metric["system"] = ag.InstanceIP
	var results []map[string]interface{}
	results = append(results, metric)

	return results, nil
}
