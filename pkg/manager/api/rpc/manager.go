package rpc

import (
	"Infinite_train/pkg/common/context"
	"Infinite_train/pkg/common/utils/log/golog"
	"fmt"
)

type ManagerRPC struct {
}

func (rpc *ManagerRPC) ReceiveHeartBeat(req *context.ReportHeartBeatRequest, reply *string) error {
	*reply = fmt.Sprintf("received agent [%s] heart beat", req.AgentIP)
	requestID := req.RequestID

	golog.Infof(requestID, "ReceiveHeartBeat, req [%+v], reply [%+v]", *req, *reply)
	return nil
}