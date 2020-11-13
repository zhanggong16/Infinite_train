package main

import (
	"Infinite_train/pkg/common/utils/linux"
	"fmt"
)

func main() {
	cmd := "ls"
	ret, err := linux.ExecWithTimeout("0", cmd, 1)
	fmt.Printf("ret [%s], err [%s]", ret, err.Error())
}
