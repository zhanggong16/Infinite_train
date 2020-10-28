package manager

import (
	"Infinite_train/pkg/common/utils/tools"
	"Infinite_train/pkg/manager/config"
	"flag"
	"fmt"
)

func Main(versionInfo *tools.VersionInfo) {
	var banner = "welcome to infinite train manager !"

	var cfgFile = flag.String("config", "/etc/manager.toml", "manager configure file absolute path!")
	var isShowVersion = false
	flag.BoolVar(&isShowVersion, "version", false, "Show version")
	flag.Parse()

	if isShowVersion {
		tools.ShowVersion(versionInfo, banner)
		return
	}
	if len(*cfgFile) <= 0 {
		fmt.Printf("the absolute path of the config file lost!\n")
		return
	}

	conf, err := config.ParseConfig(*cfgFile)
	if err != nil {
		fmt.Printf("parse config file failed!\n")
		return
	}


	return
}