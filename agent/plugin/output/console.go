package main

import(
	"os"
	"errors"
	"encoding/json"

	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/protocol"
)

var GlobalNodeInfo config.NodeInfo
var GlobalConfig map[string]string

func Init(nodeInfo config.NodeInfo, config map[string]string) error {
	value, ok := config["type"]

	if !ok {
		return errors.New("Missing config 'type'")
	}

	if value != "stdout" && value != "stderr" {
		return errors.New("Config 'type' error, should be 'stdout' or 'stderr'")
	}

	GlobalConfig = make(map[string]string)
	GlobalConfig = config
	GlobalNodeInfo = nodeInfo

	return nil
}

func Send(proto *protocol.Proto) error {
	body, err := json.MarshalIndent(proto, "", "    ")

	if err != nil {
		return err
	}

	if GlobalConfig["type"] == "stdout" {
		os.Stdout.WriteString(string(body) + "\r\n")
	}

	if GlobalConfig["type"] == "stderr" {
		os.Stderr.WriteString(string(body) + "\r\n")
	}

	return nil
}
