package collector

type MySQLConnect struct {
	*DatabaseBaseStruct
}

func (mc *MySQLConnect) DbPerform() ([]map[string]interface{}, error) {
	var metric map[string]interface{}
	metric["system"] = mc.InstanceIP
	var results []map[string]interface{}
	results = append(results, metric)

	return results, nil
}