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

type InputPlugin struct {
	Node config.NodeInfo
	Config config.InputPluginInfo
	CollectQueue queue.TransferQueue

	Plugin *plugin.Plugin

	InitFunc plugin.Symbol
	CollectFunc plugin.Symbol
}

func NewInputPlugin(nodeInfo config.NodeInfo, configInfo config.InputPluginInfo, bufferSize int) *InputPlugin {
	return &InputPlugin{
		Node: nodeInfo,
		Config: configInfo,
		CollectQueue: *queue.NewTransferQueue(bufferSize),
		Plugin: nil,
	}
}

func (inputPlugin *InputPlugin) Init () error {
	//Initialize plugin
	p, err := plugin.Open(inputPlugin.Config.Path)

	if err != nil {
		return err
	}

	inputPlugin.Plugin = p

	InitFunc, err := inputPlugin.Plugin.Lookup("Init")

	if err != nil {
		return err
	}

	inputPlugin.InitFunc = InitFunc

	CollectFunc, err := inputPlugin.Plugin.Lookup("Collect")

	if err != nil {
		return err
	}

	inputPlugin.CollectFunc = CollectFunc

	//Call plugin interface to initialize
	err = inputPlugin.InitFunc.(func(config.NodeInfo, map[string]string) error)(inputPlugin.Node, inputPlugin.Config.PluginConfig)

	if err != nil {
		return err
	}

	return nil
}

func (inputPlugin *InputPlugin) Run () {
	//Loop to call plugin interface to collect data
	for {
		select {
		case <- time.After(time.Second * time.Duration(inputPlugin.Config.Duration)):
			//Call plugin Collect function
			data, err := inputPlugin.CollectFunc.(func()(*protocol.Proto, error))()

			if err != nil {
				log.Warnf("Collect data failed! error:%s", err)
				continue
			}

			//log.Infof("Collect data from %s, data:%s", inputPlugin.Config.Name, data)

			data.Name = inputPlugin.Config.Name

			//Push to collect queue
			inputPlugin.CollectQueue.Push(data)
		}
	}
}

type InputPluginManager struct {
	Node config.NodeInfo
	Configs []config.InputPluginInfo
	Plugins map[string]*InputPlugin
	TransferQueue *queue.TransferQueue
}

func NewInputPluginManager(nodeInfo config.NodeInfo, configInfos []config.InputPluginInfo, transferQueue *queue.TransferQueue) *InputPluginManager {
	return &InputPluginManager{
		Node: nodeInfo,
		Configs: configInfos,
		Plugins: map[string]*InputPlugin{},
		TransferQueue: transferQueue,
	}
}

func (manager *InputPluginManager) Init () error {
	//Loop to initialize all input plugins
	pluginActiveNumber := 0

	for _, pluginConfig := range manager.Configs {
		if !pluginConfig.Active {
			continue
		}

		log.Info("Initialize input plugin, plugin name:", pluginConfig.Name)
		_, ok := manager.Plugins[pluginConfig.Name]

		if ok {
			log.Warnf("Initialize input plugin failed! plugin name:%s, error:All ready started", pluginConfig.Name)
			return errors.New("All ready started, plugin name:" + pluginConfig.Name)
		}

		plugin := NewInputPlugin(manager.Node, pluginConfig, 1000)

		err := plugin.Init()

		if err != nil {
			return err
		}

		manager.Plugins[pluginConfig.Name] = plugin

		log.Info("Initialize input plugin successed! plugin name:", pluginConfig.Name)

		pluginActiveNumber += 1
	}

	if pluginActiveNumber == 0 {
		return errors.New("No input plugin active!")
	}

	return nil
}

func (manager *InputPluginManager) Run () {
	//Loop to run all input plugins
	for _, plugin := range manager.Plugins {
		if !plugin.Config.Active {
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
		}(&plugin.CollectQueue, manager.TransferQueue)
	}
}
