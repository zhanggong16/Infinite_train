package concurrence

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer func() {
			fmt.Println("goroutine exit")
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Second)
				fmt.Println("default")
			}
		}
	}()
	time.Sleep(time.Second * 2)
	cancel()
	time.Sleep(2 * time.Second)
}
