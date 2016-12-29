package main

import(
	"time"
	"errors"

	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/protocol"

	"github.com/akhenakh/statgo"
)

var GlobalNodeInfo config.NodeInfo
var GlobalConfig map[string]string

var GlobalStat *statgo.Stat

func Init(nodeInfo config.NodeInfo, config map[string]string) error {
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
	cpuStat := GlobalStat.CPUStats()

	if cpuStat == nil {
		return nil, errors.New("CPUStats failed! error:got nil")
	}

	proto := protocol.NewProto(1)

	curTime := time.Now()
	currentTime := curTime.Local().Format("2006-01-02 15:04:05")

	data := protocol.NewData()
	data.Time = currentTime

	data.Tag["node_name"] = GlobalNodeInfo.Name
	data.Tag["node_ip"] = GlobalNodeInfo.IP
	data.Field["user"] = cpuStat.User
	data.Field["kernel"] = cpuStat.Kernel
	data.Field["idle"] = cpuStat.Idle
	data.Field["iowait"] = cpuStat.IOWait
	data.Field["swap"] = cpuStat.Swap
	data.Field["nice"] = cpuStat.Nice
	data.Field["loadmin1"] = cpuStat.LoadMin1
	data.Field["loadmin5"] = cpuStat.LoadMin5
	data.Field["loadmin15"] = cpuStat.LoadMin15

	proto.DataList = append(proto.DataList, *data)

	return proto, nil
}
