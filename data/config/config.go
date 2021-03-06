package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	BindAddrServer string `toml:"bind_addr_server"`
	BindAddrHost   string `toml:"bind_addr_host"`
	BindAddrOuter  string `toml:"bind_addr_outer"`
}

func NewConfig() *Config {
	return &Config{
		BindAddrServer: ":9080",
		BindAddrHost:   "localhost",
		BindAddrOuter:  ":3000",
	}
}
func NewConfigPath(configPath string) (*Config, error) {
	config := NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return config, err
	}
	return config, nil
}
