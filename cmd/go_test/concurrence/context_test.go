package concurrence

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// 关联任务的取消
// 根context， context.Background()创建；子Context， context.WithCancel(context.Background())创建
// ctx, cancel := context.WithCancel(context.Background())
// cancel 是个方法，可以取消Context的，基于他的子context都会被取消。ctx是可以传入子任务的。 <- ctx.Done()

// 外部控制子任务取消

func isCancelled(ctx context.Context) bool {
	select {
	case <- ctx.Done():
		return true
	default:
		return false
	}
}


func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	for i:=0;i<5;i++{
		go func(i int, ctx context.Context) {
			for {
				if isCancelled(ctx) {
					break
				}
				time.Sleep(time.Millisecond*5)
			}
			fmt.Println(i, "Cancelled")

		}(i, ctx)
	}
	cancel() // 执行取消方法，这是传入子任务的ctx会有个消息通知，ctx.Done()。
	time.Sleep(time.Second*1)
}


