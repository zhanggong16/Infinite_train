package rpc

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"bufio"
	"encoding/gob"
	"net"
	"net/rpc"
	"time"
)

//Client is a rpc class.
type Client struct {
	address      string
	dialTimeout  time.Duration
	codecTimeout time.Duration
}

//NewRPCClient is a constructor.
//Attention: codecTimeout + dialTimeout = the sum of timeout
func NewRPCClient(address string, dialTimeout time.Duration, codecTimeout time.Duration) *Client {
	return &Client{address, dialTimeout, codecTimeout}
}

//Call is a synchronous call method.
func (r *Client) Call(rpcname string, args interface{}, reply interface{}) error {

	conn, err := net.DialTimeout("tcp", r.address, time.Second*r.dialTimeout)
	if err != nil {
		golog.Warnf("0", "net.DialTimeout:%s", err.Error())
		time.Sleep(1 * time.Second)
		var rt int
		for rt = 0; rt < 5; rt++ {
			conn, err = net.DialTimeout("tcp", r.address, time.Second*r.dialTimeout)
			if err != nil {
				golog.Warnf("0", "net.DialTimeout:%s", err.Error())
				time.Sleep(1 * time.Second)
				continue
			} else {
				break
			}
		}
		if rt == 5 {
			golog.Warn("0", "net.DialTimeout: has retried 5 times")
			return err
		}
	}

	encBuf := bufio.NewWriter(conn)
	codec := &gobClientCodec{conn, gob.NewDecoder(conn), gob.NewEncoder(encBuf), encBuf, r.codecTimeout}
	c := rpc.NewClientWithCodec(codec)
	err = c.Call(rpcname, args, reply)
	if err != nil {
		golog.Errorf("0", "rpc.Call err:%s", err.Error())
	}
	clErr := c.Close()
	if clErr != nil {
		golog.Error("0", clErr.Error())
	}

	return err
}

/*// CallRPCInterface ...
func CallRPCInterface(rpcName string, rpcClientCfg *config.RPCClient, req interface{}, reply interface{}, retryOptions ...retry.RetryOption) error {
	c := NewRPCClient(rpcClientCfg.Address, rpcClientCfg.DialTimeout, rpcClientCfg.CodecTimeout)
	op := func() error {
		return c.Call(rpcName, req, reply)
	}
	if len(retryOptions) == 0 {
		//not input retry options from out of this api
		return retry.Do(op, retry.Timeout(0), retry.MaxTries(1), retry.RetryChecker(errgo.Any))
	}
	return retry.Do(op, retryOptions...)
}*/
