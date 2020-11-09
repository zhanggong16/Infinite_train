package main

import (
	"Infinite_train/pkg/common/utils/version"
	"Infinite_train/pkg/executes/agent"
	"Infinite_train/pkg/hack"
)

func main() {
	versionInfo := &version.VersionInfo{ReleaseVersion: hack.ReleaseVersion, Version: hack.Version, Compile: hack.Compile}
	agent.Main(versionInfo)
}