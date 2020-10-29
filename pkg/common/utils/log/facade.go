package log

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"strings"
)

//SetLogLevel is to set log level.
func SetLogLevel(GlobalSysLogger *golog.Logger, level string) {
	switch strings.ToLower(level) {
	case "debug":
		GlobalSysLogger.SetLevel(golog.LevelDebug)
	case "info":
		GlobalSysLogger.SetLevel(golog.LevelInfo)
	case "warn":
		GlobalSysLogger.SetLevel(golog.LevelWarn)
	case "error":
		GlobalSysLogger.SetLevel(golog.LevelError)
	case "fatal":
		GlobalSysLogger.SetLevel(golog.LevelFatal)
	default:
		GlobalSysLogger.SetLevel(golog.LevelError)

	}
}
