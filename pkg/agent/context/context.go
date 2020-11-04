package context

import "Infinite_train/pkg/agent/config"

type Context struct {
	Config		*config.Config
	LocalIP		string
}

var Agent = &Context{}
