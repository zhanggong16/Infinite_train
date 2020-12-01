package history

import (
	"encoding/json"
	"fmt"
	"testing"
)

type BaseInfo struct {
	Name	string	`json:"name"`
	Age		int		`json:"age"`
}

type JobInfo struct {
	Skills []string	`json:"skills"`
}

type EmployeeInfo struct {
	BaseInfo BaseInfo `json:"basic_info"`
	JobInfo  JobInfo  `json:"job_info"`
}

var jsonStr = `{
	"basic_info": {
		"name": "Mike",
		"age": 30
	},
	"job_info": {
		"skills": ["Java", "Go"]
	}
}`

func TestJsonCommon(t *testing.T) {
	e := new(EmployeeInfo)
	// 字符串到struct的填充
	if err := json.Unmarshal([]byte(jsonStr), e); err !=nil {
		t.Error(err)
	} else {
		fmt.Println(*e, e.JobInfo.Skills)
	}


	if v, err := json.Marshal(e); err == nil {
		fmt.Println(string(v))
	} else {
		t.Error(err)
	}
}