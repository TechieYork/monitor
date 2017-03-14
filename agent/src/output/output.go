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

//Output plugin
type OutputPlugin struct {
	node config.NodeInfo                    //Node information
	config config.OutputPluginInfo          //Plugin information
	sendQueue queue.TransferQueue           //Transfer queue

	plugin *plugin.Plugin                   //Plugin pointer

	initFunc plugin.Symbol                  //Init function symbol
	sendFunc plugin.Symbol                  //Send function symbol
}

func NewOutputPlugin(nodeInfo config.NodeInfo, configInfo config.OutputPluginInfo, bufferSize int) *OutputPlugin {
	return &OutputPlugin{
		node: nodeInfo,
		config: configInfo,
		sendQueue: *queue.NewTransferQueue(bufferSize),
	}
}

//Init plugin
func (outputPlugin *OutputPlugin) Init () error {
	//Check input plugin, if no inputs specify then default as all
	if len(outputPlugin.config.Inputs) != 0 {
		inputActiveCount := 0

		for _, isActive := range outputPlugin.config.Inputs {
			if !isActive {
				continue
			}

			inputActiveCount += 1
		}

		if inputActiveCount == 0 {
			return errors.New("No '" + outputPlugin.config.Name + "' output plugin's input plugin active!")
		}
	}
	//Initialize plugin
	p, err := plugin.Open(outputPlugin.config.Path)

	if err != nil {
		return err
	}

	outputPlugin.plugin = p

	InitFunc, err := outputPlugin.plugin.Lookup("Init")

	if err != nil {
		return err
	}

	outputPlugin.initFunc = InitFunc

	SendFunc, err := outputPlugin.plugin.Lookup("Send")

	if err != nil {
		return err
	}

	outputPlugin.sendFunc = SendFunc

	//Call plugin interface to initialize
	err = outputPlugin.initFunc.(func(config.NodeInfo, map[string]string) error)(outputPlugin.node, outputPlugin.config.PluginConfig)

	if err != nil {
		return err
	}

	return nil
}

//Run to send data
func (outputPlugin *OutputPlugin) Run () {
	//Loop to call plugin interface to send data
	for {
		//Pop from transfer queue
		data, err := outputPlugin.sendQueue.Pop(time.Millisecond * 10)

		if err != nil {
			continue
		}

		//log.Infof("send data to %s, data:%s", outputPlugin.Config.Name, data)

		//Call plugin Send function
		err = outputPlugin.sendFunc.(func(*protocol.Proto) error)(data)

		if err != nil {
			log.Warnf("Send data failed! error:%s, data:%s", err, data)
			continue
		}
	}
}

//Output plugin manager
type OutputPluginManager struct {
	node config.NodeInfo                        //Node information
	configs []config.OutputPluginInfo           //All output plugin configs
	plugins map[string]*OutputPlugin            //All plugins
	transferQueue *queue.TransferQueue          //Transfer queue
}

func NewOutputPluginManager(nodeInfo config.NodeInfo, configInfos []config.OutputPluginInfo, transferQueue *queue.TransferQueue) *OutputPluginManager {
	return &OutputPluginManager{
		node: nodeInfo,
		configs: configInfos,
		plugins: map[string]*OutputPlugin{},
		transferQueue: transferQueue,
	}
}

//Init output plugin manager
func (manager *OutputPluginManager) Init () error {
	//Loop to initialize all output plugins
	pluginActiveNumber := 0

	for _, pluginConfig := range manager.configs {
		if !pluginConfig.Active {
			continue
		}

		log.Info("Initialize output plugin, plugin name:", pluginConfig.Name)
		_, ok := manager.plugins[pluginConfig.Name]

		if ok {
			log.Warnf("Initialize output plugin failed! plugin name:%s, error:All ready started", pluginConfig.Name)
			return errors.New("All ready started, plugin name:" + pluginConfig.Name)
		}

		plugin := NewOutputPlugin(manager.node, pluginConfig, 1000)

		err := plugin.Init()

		if err != nil {
			return err
		}

		manager.plugins[pluginConfig.Name] = plugin

		log.Info("Initialize output plugin successed! plugin name:", pluginConfig.Name)

		pluginActiveNumber += 1
	}

	if pluginActiveNumber == 0 {
		return errors.New("No output plugin active!")
	}

	return nil
}

//Run to dispatch data to output plugin
func (manager *OutputPluginManager) Run () {
	//Loop to run all input plugins
	for _, plugin := range manager.plugins {
		if !plugin.config.Active {
			continue
		}

		go plugin.Run()
	}

	//Loop to pop data from transfer queue and push into send queue
	go func(manager *OutputPluginManager) {
		for {
			data, err := manager.transferQueue.Pop(time.Millisecond * 100)

			if err != nil {
				continue
			}

			for _, plugin := range manager.plugins {
				//Check is this input plugin in the output plugin's inputs map
				isActive, ok := plugin.config.Inputs[data.Name]

				if !ok || !isActive {
					continue
				}

				//Send
				err = plugin.sendQueue.Push(data)

				if err != nil {
					log.Warnf("InputPlugin transfer failed! error:%s", err)
				}
			}
		}
	}(manager)
}

