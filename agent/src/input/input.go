package input

import (
	"time"
	"errors"
	"plugin"

	"github.com/DarkMetrix/monitor/agent/src/protocol"
	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/queue"

	log "github.com/cihub/seelog"
)

//Input plugin
type InputPlugin struct {
	node config.NodeInfo
	config config.InputPluginInfo
	collectQueue queue.TransferQueue

	plugin *plugin.Plugin

	initFunc plugin.Symbol
	collectFunc plugin.Symbol
}

func NewInputPlugin(nodeInfo config.NodeInfo, configInfo config.InputPluginInfo, bufferSize int) *InputPlugin {
	return &InputPlugin{
		node: nodeInfo,
		config: configInfo,
		collectQueue: *queue.NewTransferQueue(bufferSize),
		plugin: nil,
	}
}

//Init plugin
func (inputPlugin *InputPlugin) Init () error {
	//Initialize plugin
	p, err := plugin.Open(inputPlugin.config.Path)

	if err != nil {
		return err
	}

	inputPlugin.plugin = p

	InitFunc, err := inputPlugin.plugin.Lookup("Init")

	if err != nil {
		return err
	}

	inputPlugin.initFunc = InitFunc

	CollectFunc, err := inputPlugin.plugin.Lookup("Collect")

	if err != nil {
		return err
	}

	inputPlugin.collectFunc = CollectFunc

	//Call plugin interface to initialize
	err = inputPlugin.initFunc.(func(config.NodeInfo, map[string]string) error)(inputPlugin.node, inputPlugin.config.PluginConfig)

	if err != nil {
		return err
	}

	return nil
}

//Run to collect data
func (inputPlugin *InputPlugin) Run () {
	//Loop to call plugin interface to collect data
	for {
		select {
		case <- time.After(time.Second * time.Duration(inputPlugin.config.Duration)):
			//Call plugin Collect function
			data, err := inputPlugin.collectFunc.(func()(*protocol.Proto, error))()

			if err != nil {
				log.Warnf("Collect data failed! error:%s", err)
				continue
			}

			//log.Infof("Collect data from %s, data:%s", inputPlugin.Config.Name, data)

			data.Name = inputPlugin.config.Name

			//Push to collect queue
			inputPlugin.collectQueue.Push(data)
		}
	}
}

//Input plugin manager
type InputPluginManager struct {
	node config.NodeInfo
	configs []config.InputPluginInfo
	plugins map[string]*InputPlugin
	transferQueue *queue.TransferQueue
}

func NewInputPluginManager(nodeInfo config.NodeInfo, configInfos []config.InputPluginInfo, transferQueue *queue.TransferQueue) *InputPluginManager {
	return &InputPluginManager{
		node: nodeInfo,
		configs: configInfos,
		plugins: map[string]*InputPlugin{},
		transferQueue: transferQueue,
	}
}

//Init plugin manager
func (manager *InputPluginManager) Init () error {
	//Loop to initialize all input plugins
	pluginActiveNumber := 0

	for _, pluginConfig := range manager.configs {
		if !pluginConfig.Active {
			continue
		}

		log.Info("Initialize input plugin, plugin name:", pluginConfig.Name)
		_, ok := manager.plugins[pluginConfig.Name]

		if ok {
			log.Warnf("Initialize input plugin failed! plugin name:%s, error:All ready started", pluginConfig.Name)
			return errors.New("All ready started, plugin name:" + pluginConfig.Name)
		}

		plugin := NewInputPlugin(manager.node, pluginConfig, 1000)

		err := plugin.Init()

		if err != nil {
			return err
		}

		manager.plugins[pluginConfig.Name] = plugin

		log.Info("Initialize input plugin successed! plugin name:", pluginConfig.Name)

		pluginActiveNumber += 1
	}

	if pluginActiveNumber == 0 {
		return errors.New("No input plugin active!")
	}

	return nil
}

//Run all plugin to collect data
func (manager *InputPluginManager) Run () {
	//Loop to run all input plugins
	for _, plugin := range manager.plugins {
		if !plugin.config.Active {
			continue
		}

		//Begin collect data
		go plugin.Run()

		//Begin transfer data from collect queue to transfer queue
		go func(collectQueue *queue.TransferQueue, transferQueue *queue.TransferQueue) {
			for {
				data, err := collectQueue.Pop(time.Millisecond * 100)

				if err != nil {
					continue
				}

				err = transferQueue.Push(data)

				if err != nil {
					log.Warnf("InputPlugin transfer failed! error:%s", err)
				}
			}
		}(&plugin.collectQueue, manager.transferQueue)
	}
}
