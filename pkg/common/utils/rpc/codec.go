package rpc

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net/rpc"
	"time"
)

//timeoutCoder: for both server and client here.
func timeoutCoder(f func(interface{}) error, e interface{}, codecTimeout time.Duration, msg string) error {
	echan := make(chan error, 1)
	go func() { echan <- f(e) }()
	select {
	case e := <-echan:
		return e
	case <-time.After(time.Second * codecTimeout):
		return fmt.Errorf("Coder Timeout %s", msg)
	}
}

//server codec here.
type gobServerCodec struct {
	rwc          io.ReadWriteCloser
	dec          *gob.Decoder
	enc          *gob.Encoder
	encBuf       *bufio.Writer
	codecTimeout time.Duration
	closed       bool
}

func (c *gobServerCodec) ReadRequestHeader(r *rpc.Request) error {
	return timeoutCoder(c.dec.Decode, r, c.codecTimeout, "server read request header")
}

func (c *gobServerCodec) ReadRequestBody(body interface{}) error {
	return timeoutCoder(c.dec.Decode, body, c.codecTimeout, "server read request body")
}

func (c *gobServerCodec) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	if err = timeoutCoder(c.enc.Encode, r, c.codecTimeout, "server write response"); err != nil {
		if c.encBuf.Flush() == nil {
			golog.Info("0", "rpc: gob error encoding response:", "err", err)
			clErr := c.Close()
			if clErr != nil {
				golog.Error("0", clErr.Error())
			}

		}
		return
	}
	if err = timeoutCoder(c.enc.Encode, body, c.codecTimeout, "server write response body"); err != nil {
		if c.encBuf.Flush() == nil {
			golog.Info("0", "rpc: gob error encoding body:", "err", err)
			clErr := c.Close()
			if clErr != nil {
				golog.Error("0", clErr.Error())
			}

		}
		return
	}
	return c.encBuf.Flush()
}

func (c *gobServerCodec) Close() error {
	if c.closed {
		// Only call c.rwc.Close once; otherwise the semantics are undefined.
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

//for client codec here.
type gobClientCodec struct {
	rwc          io.ReadWriteCloser
	dec          *gob.Decoder
	enc          *gob.Encoder
	encBuf       *bufio.Writer
	codecTimeout time.Duration
}

func (c *gobClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	if err = timeoutCoder(c.enc.Encode, r, c.codecTimeout, "client write request"); err != nil {
		return
	}
	if err = timeoutCoder(c.enc.Encode, body, c.codecTimeout, "client write request body"); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *gobClientCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.dec.Decode(r)
}

func (c *gobClientCodec) ReadResponseBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *gobClientCodec) Close() error {
	return c.rwc.Close()
}
