package main

import(
	"errors"
	"encoding/json"

	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/protocol"

	"github.com/nsqio/go-nsq"
)

var GlobalNodeInfo config.NodeInfo
var GlobalConfig map[string]string
var GlobalProducer *nsq.Producer

func Init(nodeInfo config.NodeInfo, config map[string]string) error {
	_, ok := config["nsqd_address"]

	if !ok {
		return errors.New("Missing config 'nsqd_address'")
	}

	_, ok = config["topic"]

	if !ok {
		return errors.New("Missing config 'topic'")
	}

	GlobalConfig = make(map[string]string)
	GlobalConfig = config
	GlobalNodeInfo = nodeInfo

	//Init nsq producer
	producer, err := nsq.NewProducer(config["nsqd_address"], nsq.NewConfig())

	if err != nil {
		return err
	}

	GlobalProducer = producer

	return nil
}

func Send(proto *protocol.Proto) error {
	body, err := json.MarshalIndent(proto, "", "    ")

	if err != nil {
		return err
	}

	err = GlobalProducer.Publish(GlobalConfig["topic"], []byte(body))

	if err != nil {
		return err
	}

	return nil
}
