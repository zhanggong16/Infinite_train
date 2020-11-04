package rpc

import (
	"Infinite_train/pkg/agent/context"
	commonContext "Infinite_train/pkg/common/context"
	"Infinite_train/pkg/common/utils/retry"
	"Infinite_train/pkg/common/utils/rpc"
	"github.com/juju/errgo"
)

func CallRPCInterface(rpcName string, req interface{}, reply interface{}, retryOptions ...retry.RetryOption) error {
	rpcClientCfg := context.Agent.Config.ManagerRPCServer
	c := rpc.NewRPCClient(rpcClientCfg.Address, rpcClientCfg.DialTimeout, rpcClientCfg.CodecTimeout)
	op := func() error {
		return c.Call(rpcName, req, reply)
	}
	if len(retryOptions) == 0 {
		return retry.Do(op, retry.Timeout(0), retry.MaxTries(1), retry.RetryChecker(errgo.Any))
	}
	return retry.Do(op, retryOptions...)
}

func CallReportHeartBeat(req *commonContext.ReportHeartBeatRequest, reply *string, retryOptions ...retry.RetryOption) error {
	return CallRPCInterface("ManagerRPC.ReceiveHeartBeat", req, reply, retryOptions...)
}
