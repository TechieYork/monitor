package output

import (
	"time"
	"errors"
	"plugin"

	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/queue"
	"github.com/DarkMetrix/monitor/agent/src/protocol"

	log "github.com/cihub/seelog"
)

type OutputPlugin struct {
	Node config.NodeInfo
	Config config.OutputPluginInfo
	SendQueue queue.TransferQueue

	Plugin *plugin.Plugin

	InitFunc plugin.Symbol
	SendFunc plugin.Symbol
}

func NewOutputPlugin(nodeInfo config.NodeInfo, configInfo config.OutputPluginInfo, bufferSize int) *OutputPlugin {
	return &OutputPlugin{
		Node: nodeInfo,
		Config: configInfo,
		SendQueue: *queue.NewTransferQueue(bufferSize),
	}
}

func (outputPlugin *OutputPlugin) Init () error {
	//Check input plugin, if no inputs specify then default as all
	if len(outputPlugin.Config.Inputs) != 0 {
		inputActiveCount := 0

		for _, isActive := range outputPlugin.Config.Inputs {
			if !isActive {
				continue
			}

			inputActiveCount += 1
		}

		if inputActiveCount == 0 {
			return errors.New("No '" + outputPlugin.Config.Name + "' output plugin's input plugin active!")
		}
	}
	//Initialize plugin
	p, err := plugin.Open(outputPlugin.Config.Path)

	if err != nil {
		return err
	}

	outputPlugin.Plugin = p

	InitFunc, err := outputPlugin.Plugin.Lookup("Init")

	if err != nil {
		return err
	}

	outputPlugin.InitFunc = InitFunc

	SendFunc, err := outputPlugin.Plugin.Lookup("Send")

	if err != nil {
		return err
	}

	outputPlugin.SendFunc = SendFunc

	//Call plugin interface to initialize
	err = outputPlugin.InitFunc.(func(config.NodeInfo, map[string]string) error)(outputPlugin.Node, outputPlugin.Config.PluginConfig)

	if err != nil {
		return err
	}

	return nil
}

func (outputPlugin *OutputPlugin) Run () {
	//Loop to call plugin interface to send data
	for {
		//Pop from transfer queue
		data, err := outputPlugin.SendQueue.Pop(time.Millisecond * 10)

		if err != nil {
			continue
		}

		//log.Infof("send data to %s, data:%s", outputPlugin.Config.Name, data)

		//Call plugin Send function
		err = outputPlugin.SendFunc.(func(*protocol.Proto) error)(data)

		if err != nil {
			log.Warnf("Send data failed! error:%s, data:%s", err, data)
			continue
		}
	}
}

type OutputPluginManager struct {
	Node config.NodeInfo
	Configs []config.OutputPluginInfo
	Plugins map[string]*OutputPlugin
	TransferQueue *queue.TransferQueue
}

func NewOutputPluginManager(nodeInfo config.NodeInfo, configInfos []config.OutputPluginInfo, transferQueue *queue.TransferQueue) *OutputPluginManager {
	return &OutputPluginManager{
		Node: nodeInfo,
		Configs: configInfos,
		Plugins: map[string]*OutputPlugin{},
		TransferQueue: transferQueue,
	}
}

func (manager *OutputPluginManager) Init () error {
	//Loop to initialize all output plugins
	pluginActiveNumber := 0

	for _, pluginConfig := range manager.Configs {
		if !pluginConfig.Active {
			continue
		}

		log.Info("Initialize output plugin, plugin name:", pluginConfig.Name)
		_, ok := manager.Plugins[pluginConfig.Name]

		if ok {
			log.Warnf("Initialize output plugin failed! plugin name:%s, error:All ready started", pluginConfig.Name)
			return errors.New("All ready started, plugin name:" + pluginConfig.Name)
		}

		plugin := NewOutputPlugin(manager.Node, pluginConfig, 1000)

		err := plugin.Init()

		if err != nil {
			return err
		}

		manager.Plugins[pluginConfig.Name] = plugin

		log.Info("Initialize output plugin successed! plugin name:", pluginConfig.Name)

		pluginActiveNumber += 1
	}

	if pluginActiveNumber == 0 {
		return errors.New("No output plugin active!")
	}

	return nil
}

func (manager *OutputPluginManager) Run () {
	//Loop to run all input plugins
	for _, plugin := range manager.Plugins {
		if !plugin.Config.Active {
			continue
		}

		go plugin.Run()
	}

	//Loop to pop data from transfer queue and push into send queue
	go func(manager *OutputPluginManager) {
		for {
			data, err := manager.TransferQueue.Pop(time.Millisecond * 100)

			if err != nil {
				continue
			}

			for _, plugin := range manager.Plugins {
				//Check is this input plugin in the output plugin's inputs map
				isActive, ok := plugin.Config.Inputs[data.Name]

				if !ok || !isActive {
					continue
				}

				//Send
				err = plugin.SendQueue.Push(data)

				if err != nil {
					log.Warnf("InputPlugin transfer failed! error:%s", err)
				}
			}
		}
	}(manager)
}

