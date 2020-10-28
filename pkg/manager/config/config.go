package config

import "fmt"
import "github.com/BurntSushi/toml"


func ParseConfig(file string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &conf, nil
}