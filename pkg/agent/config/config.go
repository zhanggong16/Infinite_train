package config

import (
	"Infinite_train/pkg/common/config"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"time"
)

type AgentCfg struct {
	IPPrefix     string   `toml:"ip_prefix"`
	Region       string   `toml:"region"`
	AdminRoles   []string `toml:"admin_roles"`
	AZones       string   `toml:"availability_zones"`
	ResourcePool string   `toml:"resource_pool"`
}

type RPCServer struct {
	Address      string
	CodecTimeout time.Duration `toml:"codec_timeout"`
}

type CronInterval struct {
	IntervalEveryMinute		uint64 `toml:"interval_every_minute"`
}

type ManagerRPCServer struct {
	Address      string
	CodecTimeout time.Duration `toml:"codec_timeout"`
	DialTimeout  time.Duration `toml:"dial_timeout"`
}

type Config struct {
	AgentConfig			*AgentCfg                  		`toml:"manager"`
	ManagerRPCServer	*ManagerRPCServer				`toml:"manager_rpc_server"`
	RPCServer     		*RPCServer                   	`toml:"rpc_server"`
	LogConfigs   		map[string]*config.LogConfig 	`toml:"logs"`
	RetryConfig   		*config.RetryConfig          	`toml:"retry"`
	PollingConfig 		*config.PollingConfig        	`toml:"polling"`
	CronInterval  		*CronInterval                	`toml:"cron_interval"`
}

// func...
func ParseConfig(file string) (*Config, error) {
	conf := new(Config)
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return conf, nil
}

func (conf *Config) String() string {
	b, err := json.Marshal(*conf)
	if err != nil {
		return fmt.Sprintf("\n%+v\n", *conf)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("\n%+v\n", *conf)
	}
	return out.String()
}
