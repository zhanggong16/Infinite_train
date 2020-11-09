package version

import "fmt"

type VersionInfo struct {
	ReleaseVersion string
	Version        string
	Compile        string
}

func ShowVersion(versionInfo *VersionInfo, banner string) {
	fmt.Printf("%s\n", banner)
	fmt.Printf("Release version:%s\n", versionInfo.ReleaseVersion)
	fmt.Printf("Git commit:%s\n", versionInfo.Version)
	fmt.Printf("Build time:%s\n", versionInfo.Compile)
}