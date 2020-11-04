package cron

import (
	"Infinite_train/pkg/agent/context"
	"Infinite_train/pkg/agent/controller"
	"github.com/go-co-op/gocron"
	"time"
)

var scheduler = gocron.NewScheduler(time.UTC)

func registerCron() {
	scheduler.Every(context.Agent.Config.CronInterval.IntervalEveryMinute).Seconds().Do(controller.ReportHeartBeat)
}

func Start() (chan struct{}, error) {
	registerCron()
	exitCh := scheduler.StartAsync()
	return exitCh, nil
}