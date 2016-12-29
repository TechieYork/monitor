package main

import(
	"time"
	"errors"
	"strings"

	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/protocol"

	"github.com/akhenakh/statgo"
)

var GlobalNodeInfo config.NodeInfo
var GlobalConfig map[string]string

var GlobalStat *statgo.Stat
var GlobalIncludes map[string]bool

func Init(nodeInfo config.NodeInfo, config map[string]string) error {
	GlobalIncludes = make(map[string]bool)

	_, ok := config["include"]

	if ok {
		includes := strings.Split(config["include"], ";")

		for _, include := range includes {
			GlobalIncludes[include] = true
		}
	}

	GlobalConfig = make(map[string]string)
	GlobalConfig = config
	GlobalNodeInfo = nodeInfo

	GlobalStat = statgo.NewStat()

	if GlobalStat == nil {
		return errors.New("NewStat failed! error:got nil")
	}

	return nil
}

func Collect()(*protocol.Proto, error) {
	interfaceInfos := GlobalStat.InteraceInfos()

	if interfaceInfos == nil {
		return nil, errors.New("InterfaceInfos failed! error:got nil")
	}

	proto := protocol.NewProto(1)

	curTime := time.Now()
	currentTime := curTime.Local().Format("2006-01-02 15:04:05")

	for _, info := range interfaceInfos {
		if len(GlobalIncludes) != 0 {
			_, ok := GlobalIncludes[info.Name]

			if !ok {
				continue
			}
		}

		data := protocol.NewData()
		data.Time = currentTime

		data.Tag["node_name"] = GlobalNodeInfo.Name
		data.Tag["node_ip"] = GlobalNodeInfo.IP

		data.Tag["name"] = info.Name
		data.Tag["factor"] = info.Factor
		data.Tag["duplex"] = info.Duplex
		data.Tag["state"] = info.State

		data.Field["speed"] = info.Speed

		proto.DataList = append(proto.DataList, *data)
	}

	return proto, nil
}
