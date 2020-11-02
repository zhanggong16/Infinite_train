package config

import (
	"Infinite_train/pkg/common/config"
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)
import "github.com/BurntSushi/toml"

type ManagerCfg struct {
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

type DataBase struct {
	Account  string `toml:"account"`
	Password string `toml:"password"`
	Port     int    `toml:"port"`
	Schema   string `toml:"schema"`
	IP       string `toml:"ip"`
	Charset  string `toml:"charset"`
	MaxIdle  int    `toml:"maxIdle"`
	MaxOpen  int    `toml:"maxOpen"`
}

type CronInterval struct {
	IntervalPingPong		uint64 `toml:"interval_ping_pong"`
}

type Config struct {
	WebAddr       string                       `toml:"web_addr"`
	ManagerConfig *ManagerCfg                  `toml:"manager"`
	RPCServer     *RPCServer                   `toml:"rpc_server"`
	DataBase      *DataBase                    `toml:"database"`
	LogConfigs    map[string]*config.LogConfig `toml:"logs"`
	RetryConfig   *config.RetryConfig          `toml:"retry"`
	PollingConfig *config.PollingConfig        `toml:"polling"`
	CronInterval  *CronInterval                `toml:"cron_interval"`
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
