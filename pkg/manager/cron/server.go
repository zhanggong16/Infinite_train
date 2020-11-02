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
}

func Start() (chan struct{}, error) {
	registerCron()
	exitCh := scheduler.StartAsync()
	return exitCh, nil
}