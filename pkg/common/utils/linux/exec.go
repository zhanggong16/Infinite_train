package linux

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"github.com/juju/errors"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

func Exec(cmd string) (string, error) {
	//golog.Debugf("cmd", "Exec cmd=[%s]", cmd)
	out, err := exec.Command("/bin/bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	cnt := strings.TrimRight(string(out), "\n")
	//golog.Debugf("cmd", "Exec cnt=[%s]", cnt)
	return cnt, nil
}

func ExecWithPolling(requestID, cmd string, interval, timeout time.Duration, condition interface{}, args interface{}) error {
	// 检查传入参数condition是否符合规范
	typ := reflect.TypeOf(condition)
	if typ.Kind() != reflect.Func {
		panic("Only function can be condition")
	}
	rc, err := condition.(func(interface{}) (bool, error))(args)
	if rc == true {
		golog.Info(requestID, "Success check finishing polling")
		return nil
	}
	if err != nil {
		msg := "Failed to check polling condition"
		golog.Info(requestID, msg)
		err := errors.New(msg)
		return err
	}

	ticker := time.NewTicker(interval * time.Second)
	if ticker == nil {
		err := errors.New("Create ticker error")
		return err
	}

	runningCmd := exec.Command("/bin/bash", "-c", cmd)
	err = runningCmd.Start()
	if err != nil {
		return err
	}

	status := make(chan error, 1)
	go func() {
		status <- runningCmd.Wait()
	}()

	var tickSum time.Duration
	for {
		select {
		case t := <-status:
			return t
		case <-ticker.C:
			tickSum += interval
			if tickSum >= timeout {
				err := errors.New("Polling timeout")
				ticker.Stop()
				return err
			}
			rc, err := condition.(func(interface{}) (bool, error))(args)
			if err != nil {
				golog.Infof("Condition check error: [%s], will kill [%s]", err.Error(), cmd)
				runningCmd.Process.Kill()
				return err
			}
			if rc == false {
				golog.Infof("Condition check failed, will kill [%s]", cmd)
				runningCmd.Process.Kill()
				return err
			}
		}
	}
}