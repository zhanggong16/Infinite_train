package cron

import (
	"Infinite_train/pkg/manager/context"
	"Infinite_train/pkg/manager/controller"
	"github.com/go-co-op/gocron"
	"time"
)

var scheduler = gocron.NewScheduler(time.UTC)

func registerCron() {
	scheduler.Every(context.Instance.Config.CronInterval.IntervalPingPong).Seconds().Do(controller.TestCon)
	scheduler.Every(1).Day().At("10:10:10").Do(controller.TestScheduler)
}

func Start() (chan struct{}, error) {
	registerCron()
	exitCh := scheduler.StartAsync()
	return exitCh, nil
}