package config

import (
	"Infinite_train/pkg/common/utils/log"
	"Infinite_train/pkg/common/utils/log/golog"
)

type LogConfig struct {
	Target          string `toml:"target"`
	Level           string `toml:"level"`
	Path            string `toml:"path, omitempty"`
	RotateMethod    string `toml:"rotate_method, omitempty"`
	RotateFileSize  uint64 `toml:"rotate_file_size, omitempty"`
	RotateFileCount uint64 `toml:"rotate_file_count, omitempty"`
}

// InitConfig is to init global log module.
func InitConfig(logConfigs map[string]*LogConfig) error {
	golog.GlobalSysLoggers = golog.GlobalSysLoggers[:0]
	for logKey, logConfig := range logConfigs {
		if logConfig.Target == "file" {
			sysLogPath := logConfig.Path
			if logConfig.RotateMethod == "rotate_by_count" {
				sysLogFile, err := golog.NewRotatingFileHandler(sysLogPath,
					int(logConfig.RotateFileSize),
					int(logConfig.RotateFileCount))
				if err != nil {
					return err
				}
				logInstance := golog.New(sysLogFile, golog.Lfile|golog.Ltime|golog.Llevel, logKey)
				log.SetLogLevel(logInstance, logConfig.Level)
				golog.GlobalSysLoggers = append(golog.GlobalSysLoggers, logInstance)
			}
		}
	}
	return nil
}