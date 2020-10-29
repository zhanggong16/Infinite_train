package context

import "Infinite_train/pkg/manager/config"

type Context struct {
	Config		*config.Config
	ManagerIP	string
}

var Instance = &Context{}