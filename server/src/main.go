package main

import (
	"errors"

	"github.com/DarkMetrix/monitor/server/src/config"
	"github.com/DarkMetrix/monitor/server/src/server"

	log "github.com/cihub/seelog"
)

//Init log
func InitLog(path string) error {
	logger, err := log.LoggerFromConfigAsFile(path)

	if err != nil {
		return err
	}

	err = log.ReplaceLogger(logger)

	if err != nil {
		return err
	}

	return nil
}

//Init config
func InitConfig(path string) (*config.Config, error) {
	globalConfig := config.GetConfig()

	if globalConfig == nil {
		return nil, errors.New("Get global config failed!")
	}

	err := globalConfig.Init(path)

	if err != nil {
		return nil, err
	}

	return globalConfig, nil
}

func main() {
	defer log.Flush()

	defer func() {
		err := recover()

		if err != nil {
			log.Critical("Got panic, err:", err)
		}
	} ()

	//Initialize log using configuration from "../conf/log.config"
	err := InitLog("../conf/log.config")

	if err != nil {
		log.Warnf("Read config failed! error:%s", err)
		return
	}

	log.Info("Starting monitor_server ...")

	//Initialize the configuration from "../conf/config.json"
	log.Info("Initialize monitor server configuration from ../conf/config.json ...")
	config, err := InitConfig("../conf/config.json")

	if err != nil {
		log.Warnf("Initialize monitor_server configuration failed! error:%s", err)
		return
	}

	log.Info("Initialize monitor_server configuration successed! config:", config)

	monitor_server := server.NewMonitorServer(config)

	err = monitor_server.Run()

	if err != nil {
		log.Warnf("Run monitor_server failed! error:%s", err)
		return
	}
}
