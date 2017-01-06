package main

import (
	"errors"
	"net/http"

	"github.com/DarkMetrix/monitor/server/src/config"
	"github.com/DarkMetrix/monitor/server/src/server"

	log "github.com/cihub/seelog"
)

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

func NodeRegister(resp http.ResponseWriter, req *http.Request)  {
	resp.Write([]byte("aaaa"))
}

func main() {
	defer log.Flush()

	//Initialize log using configuration from "../conf/monitor_server_log.config"
	err := InitLog("../conf/monitor_server_log.config")

	if err != nil {
		log.Warnf("Read config failed! error:%s", err)
		return
	}

	log.Info("Starting monitor_server ...")

	//Initialize the configuration from "../conf/monitor_server_config.json"
	log.Info("Initialize monitor_agent configuration from ../conf/monitor_server_config.json ...")
	config, err := InitConfig("../conf/monitor_server_config.json")

	if err != nil {
		log.Warnf("Initialize monitor_agent configuration failed! error:%s", err)
		return
	}

	log.Info("Initialize monitor_agent configuration successed! config:", config)

	monitor_server := server.NewMonitorServer(config)

	err = monitor_server.Run()

	if err != nil {
		log.Warnf("Run monitor server failed! error:%s", err)
		return
	}
}
