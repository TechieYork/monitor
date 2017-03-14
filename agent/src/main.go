package main

import (
	"os"
	"os/signal"
	"time"
	"errors"
	"plugin"

	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/queue"

	"github.com/DarkMetrix/monitor/agent/src/input"
	"github.com/DarkMetrix/monitor/agent/src/output"

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

//Init agent plugin libraries
func InitPluginLibs(config *config.Config) error {
	inputs := make(map[string]bool)

	for _, pluginConfig := range config.Inputs {
		if !pluginConfig.Active {
			continue
		}

		//Initialize plugin
		_, err := plugin.Open(pluginConfig.Path)

		if err != nil {
			return err
		}

		inputs[pluginConfig.Name] = true
	}

	for _, pluginConfig := range config.Outputs {
		if !pluginConfig.Active {
			continue
		}

		//Initialize plugin
		_, err := plugin.Open(pluginConfig.Path)

		if err != nil {
			return err
		}

		if len(pluginConfig.Inputs) == 0 {
			continue
		}

		for inputName, isActive := range pluginConfig.Inputs {
			if !isActive {
				continue
			}

			_, ok := inputs[inputName]

			if !ok {
				return errors.New("'" + pluginConfig.Name + "' output plugin's input plugin '" + inputName + "' not found or not active!")
			}
		}
	}

	return nil
}

//Init transfer queue
func InitTransferQueue(bufferSize int) (*queue.TransferQueue, error) {
	transfer := queue.NewTransferQueue(bufferSize)

	return transfer, nil
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

	log.Info(time.Now().String(), "Starting monitor agent ... ")

	//Initialize the configuration from "../conf/config.json"
	log.Info("Initialize monitor_agent configuration from ../conf/config.json ...")
	config, err := InitConfig("../conf/config.json")

	if err != nil {
		log.Warnf("Initialize monitor agent configuration failed! error:%s", err)
		return
	}

	log.Info("Initialize monitor agent configuration successed! config:", config)

	//Initialize all plugin libs
	log.Info("Initialize monitor agent plugin libs ...")
	err = InitPluginLibs(config)

	if err != nil {
		log.Warnf("Initialize monitor agent plugin libs failed! error:%s", err)
		return
	}

	log.Info("Initialize monitor agent plugin libs successed!")

	//Init queue between input plugin and output plugin
	log.Info("Initialize monitor agent transfer queue ...")
	transfer, err := InitTransferQueue(config.Node.TransferQueue.BufferSize)

	if err != nil {
		log.Warnf("Initialize monitor agent transfer queue failed! error:%s", err)
		return
	}

	log.Info("Initialize monitor agent transfer queue successed! buffer size:", config.Node.TransferQueue.BufferSize)

	//Start output plugins
	log.Info("Initialize monitor agent output plugin ...")
	outputPluginManager := output.NewOutputPluginManager(config.Node, config.Outputs, transfer)

	err = outputPluginManager.Init()

	if err != nil {
		log.Warnf("Initialize monitor agent output plugin failed! error:%s", err)
		return
	}

	log.Info("Initialize monitor agent output plugin successed!")

	outputPluginManager.Run()

	//Start input plugins
	log.Info("Initialize monitor agent input plugin ...")
	inputPluginManager := input.NewInputPluginManager(config.Node, config.Inputs, transfer)

	err = inputPluginManager.Init()

	if err != nil {
		log.Warnf("Initialize monitor agent input plugin failed! error:%s", err)
		return
	}

	log.Info("Initialize monitor agent input plugin successed!")

	inputPluginManager.Run()

	//Deal with signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)

	signalOccur := <- signalChannel

	log.Info("Signal occured, signal:", signalOccur.String())
}
