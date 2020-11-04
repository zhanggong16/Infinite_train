package main

import (
	"Infinite_train/pkg/common/utils"
	"Infinite_train/pkg/executes/agent"
	"Infinite_train/pkg/hack"
)

func main() {
	versionInfo := &utils.VersionInfo{ReleaseVersion: hack.ReleaseVersion, Version: hack.Version, Compile: hack.Compile}
	agent.Main(versionInfo)
}