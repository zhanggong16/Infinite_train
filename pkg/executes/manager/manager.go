package manager

import (
	"Infinite_train/pkg/common/utils/tools"
	"flag"
	"fmt"
)

func Main(versionInfo *tools.VersionInfo) {
	var banner = "welcome to infinite train manager !"

	//var cfgFile = flag.String("config", "/etc/manager.toml", "manager configure file absolute path!")
	var isShowVersion = true
	flag.BoolVar(&isShowVersion, "version", false, "Show version")
	flag.Parse()

	if isShowVersion {
		tools.ShowVersion(versionInfo, banner)
		return
	}
	fmt.Printf("start...")


	return
}