package rpc

import (
	"Infinite_train/pkg/common/context"
	"Infinite_train/pkg/common/utils/log/golog"
)

type ManagerRPC struct {
}

func (rpc *ManagerRPC) PingPong(req *context.PingPongRequest, reply *string) error {
	*reply = "okay"
	requestId := req.RequestId

	golog.Infof(requestId, "ReportHeatBeat, req [%+v], reply [%+v]", *req, *reply)
	return nil
}