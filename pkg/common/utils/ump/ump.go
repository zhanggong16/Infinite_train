package ump

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"encoding/json"
)

//Record ...
type Record struct {
	StartTime    string `json:"time"`
	Key          string `json:"key"`
	AppName      string `json:"appName"`
	Hostname     string `json:"hostname"`
	ProcessState string `json:"processState"`
	ElapsedTime  string `json:"elapsedTime"`
	RequestID    string `json:"request_id"`
}

// WriteToFile ...
func (UmpR *Record) WriteToFile() {
	toJSON, _ := json.Marshal(UmpR)
	golog.Infof("0", "UmpRecord: %+s", string(toJSON))
	golog.Ump(string(toJSON))
}

