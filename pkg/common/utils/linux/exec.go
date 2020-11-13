package linux

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"bytes"
	"github.com/juju/errors"
	"os/exec"
	"reflect"
	"strings"
	"syscall"
	"time"
)

func Exec(requestID, cmd string) (string, error) {
	golog.Infof(requestID, "Exec cmd=[%s]", cmd)
	out, err := exec.Command("/bin/bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	cnt := strings.TrimRight(string(out), "\n")
	//golog.Debugf("cmd", "Exec cnt=[%s]", cnt)
	return cnt, nil
}

func ExecWithTimeout(requestID, cmd string, timeout time.Duration) (string, error) {
	//golog.Infof(requestID, "Exec with timeout cmd=[%s], timeout [%s]", cmd, string(timeout))
	execErrChan := make(chan error)
	defer close(execErrChan)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	runningCmd := exec.Command("/bin/bash", "-c", cmd)
	runningCmd.Stdout = &stdout
	runningCmd.Stderr = &stderr

	err := runningCmd.Start()
	if err != nil {
		return "", err
	}

	go func() {
		execErrChan <- runningCmd.Wait()
	}()

	select {
	case err = <-execErrChan:
		if err != nil {
			return "", err
		} else {
			stdoutStr := strings.TrimSpace(string(bytes.TrimRight(stdout.Bytes(), "\x00")))
			stderrStr := strings.TrimSpace(string(bytes.TrimRight(stderr.Bytes(), "\x00")))
			if stderrStr != "" {
				return stderrStr, nil
			} else {
				return stdoutStr, nil
			}
		}
	case <-time.After(time.Second * timeout):
		runningCmd.Process.Signal(syscall.SIGINT)
		time.Sleep(time.Second)
		runningCmd.Process.Kill()
		err = errors.New("Polling timeout")
		return "", err
	}
}

func ExecWithPollingCondition(requestID, cmd string, interval, timeout time.Duration, condition interface{}, args interface{}) error {
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