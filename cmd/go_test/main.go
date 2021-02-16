package main

import (
	"Infinite_train/cmd/go_test/plugin"
	"Infinite_train/cmd/go_test/plugin/worker/opengauss"
	"fmt"
)

func main() {
	op := plugin.NewOperator(100)
	action := opengauss.NewManualHa("zhg", "MogHa")
	op.RegisterWorker("zhg", action)
	if err := op.Start(); err != nil {
		fmt.Printf("start error %v\n", err)
	}
}

