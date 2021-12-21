package main

import (
	"github.com/daria/PortMicroserviceClient/cmd/client"
	cnfg "github.com/daria/PortMicroserviceClient/data/config"
)

var (
	configPath = "configs/dataConfig.toml"
)

func configService() *cnfg.Config {
	config, _ := cnfg.NewConfigPath(configPath)
	client.JsonPath = config.JsonPath
	client.Connect.InitConn(config.BindAddrServer)
	return config
}

func main() {
	config := configService()
	err := client.Start(config)
	if err != nil {
		return
	}
	defer client.Connect.Conn.Close()
}
