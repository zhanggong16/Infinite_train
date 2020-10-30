package utils

import (
	"github.com/juju/errors"
	"os/exec"
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

func ExecWithPolling(cmd string, interval, timeout time.Duration, condition interface{}, args interface{}) (string, error) {
	ticker := time.NewTicker(interval * time.Second)
	if ticker == nil {
		err := errors.New("create ticker error")
		return "", err
	}

	runningCmd := exec.Command("/bin/bash", "-c", cmd)
	err := runningCmd.Start()
	if err != nil {
		return "", err
	}

	status := make(chan error, 1)
	go func() {
		status <- runningCmd.Wait()
	}()

	var tickSum time.Duration
	for {
		select {
		case t := <-status:
			return "", t
		case <-ticker.C:
			tickSum += interval
			if tickSum >= timeout {
				err := errors.New("poll timeout")
				ticker.Stop()
				return "", err
			}
			_, err := condition.(func(interface{}) (bool, error))(args)
			if err != nil {
				//golog.Infof("check error: %s, will kill [%s]", err.Error(), cmd)
				runningCmd.Process.Kill()
				return "", err
			}
		}
	}
}