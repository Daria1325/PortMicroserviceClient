package main

import (
	"fmt"
	"github.com/daria/PortMicroserviceClient/cmd/client"
	cnfg "github.com/daria/PortMicroserviceClient/data/config"
)

var (
	configPath = "configs/dataConfig.toml"
)

func configService() (*cnfg.Config, error) {
	config, err := cnfg.NewConfigPath(configPath)
	if err != nil {
		return nil, err
	}
	err = client.Connect.InitConn(config.BindAddrHost, config.BindAddrServer)
	if err != nil {
		return nil, err
	}
	return config, nil
}

//check errors and handle
//[]port is not good in domain DONE
//grpc stream
//файл загрузка Done

func main() {
	config, err := configService()
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}
	err = client.Start(config)
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}
	defer client.Connect.Conn.Close()
}
