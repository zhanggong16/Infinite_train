package main

import (
	"Infinite_train/pkg/common/utils/tools"
	"Infinite_train/pkg/executes/manager"
	"Infinite_train/pkg/hack"
)

func main() {
	versionInfo := &tools.VersionInfo{ReleaseVersion: hack.ReleaseVersion, Version: hack.Version, Compile: hack.Compile}
	manager.Main(versionInfo)
	return
}
