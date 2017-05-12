package main

import (
	"os"
	"os/signal"
	"net"
	"time"
	"errors"
	"flag"
	"plugin"

	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/queue"

	"github.com/DarkMetrix/monitor/agent/src/input"
	"github.com/DarkMetrix/monitor/agent/src/output"

	log "github.com/cihub/seelog"
)

//Init log
func InitLog(path string) {
	logger, err := log.LoggerFromConfigAsFile(path)

	if err != nil {
		panic(err)
	}

	err = log.ReplaceLogger(logger)

	if err != nil {
		panic(err)
	}
}

//Get local machine ip
func getLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Warn("Get local ip failed! err:" + err.Error())
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", nil
}

//Init config
func InitConfig(path string) *config.Config{
	//Parse flag
	configPath := flag.String("config_path", "", "The config file path, default:''(empty)")

	flag.Parse()

	if len(*configPath) != 0 {
		path = *configPath
	}

	log.Info("Initialize monitor_agent configuration from " + path + " ...")

	globalConfig := config.GetConfig()

	if globalConfig == nil {
		panic(errors.New("Get global config failed!"))
	}

	err := globalConfig.Init(path)

	if err != nil {
		panic(err)
	}

	//Check ip and name, if empty use host name as the name and use one of the local ip as the ip
	if len(globalConfig.Node.Name) == 0 {
		globalConfig.Node.Name, err = os.Hostname()

		if err != nil {
			panic(err)
		}
	}

	if len(globalConfig.Node.IP) == 0 {
		globalConfig.Node.IP, err = getLocalIp()

		if err != nil {
			panic(err)
		}
	}

	return globalConfig
}

//Init agent plugin libraries
func InitPluginLibs(config *config.Config) {
	log.Info("Initialize monitor agent plugin libs ...")

	inputs := make(map[string]bool)

	for _, pluginConfig := range config.Inputs {
		if !pluginConfig.Active {
			continue
		}

		//Initialize plugin
		_, err := plugin.Open(pluginConfig.Path)

		if err != nil {
			panic(err)
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
			panic(err)
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
				panic(errors.New("'" + pluginConfig.Name + "' output plugin's input plugin '" + inputName + "' not found or not active!"))
			}
		}
	}
}

//Init transfer queue
func InitTransferQueue(bufferSize int) *queue.TransferQueue {
	log.Info("Initialize monitor agent transfer queue ...")

	transfer := queue.NewTransferQueue(bufferSize)

	return transfer
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
	InitLog("../conf/log.config")

	log.Info(time.Now().String(), "Starting monitor agent ... ")
	log.Info("Version: " + config.Version)

	//Initialize the configuration from "../conf/config.json"
	config := InitConfig("../conf/config.json")

	//Initialize all plugin libs
	InitPluginLibs(config)

	//Init queue between input plugin and output plugin
	transfer := InitTransferQueue(config.Node.TransferQueue.BufferSize)

	//Start output plugins
	log.Info("Initialize monitor agent output plugin ...")
	outputPluginManager := output.NewOutputPluginManager(config.Node, config.Outputs, transfer)

	err := outputPluginManager.Init()

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
