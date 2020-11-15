package linux

import (
	"fmt"
	"testing"
)

func TestExecWithTimeout(t *testing.T) {
	cmd := "ls"
	ret, err := ExecWithTimeout("0", cmd, 1)
	if err != nil {
		fmt.Printf("ret [%s], err [%s]", ret, err.Error())
	} else {
		fmt.Printf("ret [%s]", ret)
	}
}