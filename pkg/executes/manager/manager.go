package manager

import (
	logCommonConfig "Infinite_train/pkg/common/config"
	"Infinite_train/pkg/common/utils"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/config"
	"Infinite_train/pkg/manager/context"
	"flag"
	"fmt"
	"runtime"
	"time"
)

func Main(versionInfo *utils.VersionInfo) {
	var banner = "welcome to infinite train manager !"

	var cfgFile = flag.String("config", "/etc/manager.toml", "manager configure file absolute path!")
	var isShowVersion = false
	flag.BoolVar(&isShowVersion, "version", false, "Show version")
	flag.Parse()

	if isShowVersion {
		utils.ShowVersion(versionInfo, banner)
		return
	}

	// init config file
	if len(*cfgFile) <= 0 {
		fmt.Printf("the absolute path of the config file lost!\n")
		return
	}
	conf, err := config.ParseConfig(*cfgFile)
	if err != nil {
		fmt.Printf("parse config file failed!\n")
		return
	}
	context.Instance.Config = conf
	err = logCommonConfig.InitConfig(conf.LogConfigs)
	if err != nil {
		golog.Error("", err.Error())
		return
	}
	golog.Info("0", banner, "time: ", time.Now())
	golog.Infof("init config", "Config: %s", conf.String())

	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 1<<20)
			stack = stack[:runtime.Stack(stack, true)]
			golog.Error("0", "manager panic, err: %s, stack: %s", err, stack)
		}
	}()

	time.Sleep(1*time.Second)

	return
}