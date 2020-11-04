package controller

import (
	agentContext "Infinite_train/pkg/agent/context"
	"Infinite_train/pkg/agent/rpc"
	"Infinite_train/pkg/common/context"
	"Infinite_train/pkg/common/utils/log/golog"
	"github.com/satori/go.uuid"
)

func ReportHeartBeat() {
	var reply string
	req := new(context.ReportHeartBeatRequest)
	RequestID := uuid.NewV4().String()
	req.RequestID = RequestID
	req.AgentIP = agentContext.Agent.LocalIP

	golog.Infof(RequestID, "start report heart beat")
	err := rpc.CallReportHeartBeat(req, &reply)
	if err != nil {
		golog.Errorf(RequestID, "err smg %s", err.Error())
	}

	golog.Infof(RequestID, "report heart beat result [%s]", reply)
}
