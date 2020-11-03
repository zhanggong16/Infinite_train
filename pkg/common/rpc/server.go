package rpc

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/common/utils/rpc"
	"Infinite_train/pkg/manager/config"
)

type Server struct {
	server *rpc.Server
	service interface{}
}

func NewServer(config *config.Config, service interface{}) (*Server, error) {
	s := new(Server)
	s.service = service
	s.server = rpc.NewRPCServer(config.RPCServer.Address, config.RPCServer.CodecTimeout)
	return s, nil
}

func (s *Server) Run() {
	err := s.RegisterRPCService()
	if err != nil {
		golog.Errorf("RegisterRPCService occurs error:%s", err.Error())
		return
	}
	s.server.ListenRPC()
	return
}

func (s *Server) RegisterRPCService() error {
	err := s.server.RegisterService(s.service)
	if err != nil {
		golog.Errorf("Register ScheduleTaskRPC error:%s", err.Error())
		return err
	}
	return nil
}

func (s *Server) Close() {
	s.server.Close()
}