package rpc

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"runtime"
	"time"
)

//Server is a rpc servr side class.
type Server struct {
	listenAddress string
	closed        chan bool
	stackSize     int
	codecTimeout  time.Duration
	netListener   net.Listener
}

//NewRPCServer is a constructor.
func NewRPCServer(addr string, codecTimeout time.Duration) *Server {
	return &Server{addr, make(chan bool, 1), 1024 * 8, codecTimeout, nil}
}

//RegisterService is to register service api, which is to be called by client side.
func (r *Server) RegisterService(svc interface{}) error {
	return rpc.Register(svc)
}

//Close is to give up listen on server sice.
func (r *Server) Close() {
	fmt.Println(" 1Close chan")
	if r.closed != nil {
		fmt.Println("Close chan")
		close(r.closed)
	}

	if r.netListener != nil {
		golog.Info("0", "rpc server stoping listening on:", "address", r.netListener.Addr())
		clErr := r.netListener.Close()
		if clErr != nil {
			golog.Error("0", clErr.Error())
		}
		r.netListener = nil
	}

}

//ListenRPC is to listen on server side for providing service call.
func (r *Server) ListenRPC() {
	l, e := net.Listen("tcp", r.listenAddress)
	if e != nil {
		golog.Info("0", "rpc server listen error ", "err", e)
		return
	}

	r.netListener = l
	for {
		select {
		case <-r.closed:
			golog.Info("0", "rpc server stoping listening ")
			return
		default:
		}
		conn, err := l.Accept()
		if err != nil {
			golog.Info("0", "rpc server side , accept rpc conn err:", "err", err)
			continue
		}
		go func(conn net.Conn) {
			//add recover protect mechanism
			defer func() {
				if err := recover(); err != nil {

					stack := make([]byte, r.stackSize)
					stack = stack[:runtime.Stack(stack, true)]

					//f := "PANIC: %s\n%s"
					//rec.Logger.Printf(f, err, stack)
					//golog.Info("0", "rpc server panic:", "err", err, "stack", stack)
					golog.Errorf("0", "rpc server panic, err:%s, stack:%s", err, stack)
				}
			}()

			buf := bufio.NewWriter(conn)
			srv := &gobServerCodec{
				rwc:          conn,
				dec:          gob.NewDecoder(conn),
				enc:          gob.NewEncoder(buf),
				encBuf:       buf,
				codecTimeout: r.codecTimeout,
			}
			err = rpc.ServeRequest(srv)
			if err != nil {
				if err != io.EOF && err != io.ErrUnexpectedEOF {
					golog.Errorf("0", "rpc server side, rpc request error:", "err", err.Error())
				} else {
					golog.Tracef("0", "rpc server side, rpc request empty:", "err", err)
				}
			}
			clErr := srv.Close()
			if clErr != nil {
				golog.Error("0", clErr.Error())
			}

		}(conn)
	}
}
