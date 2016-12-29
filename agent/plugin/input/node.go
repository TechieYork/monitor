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
var GlobalHostInfo *statgo.HostInfos

func Init(nodeInfo config.NodeInfo, config map[string]string) error {
	GlobalConfig = make(map[string]string)
	GlobalConfig = config
	GlobalNodeInfo = nodeInfo

	GlobalStat = statgo.NewStat()

	if GlobalStat == nil {
		return errors.New("NewStat failed! error:got nil")
	}

	GlobalHostInfo = GlobalStat.HostInfos()

	if GlobalHostInfo == nil {
		return errors.New("HostInfos failed! error:got nil")
	}

	return nil
}

func Collect()(*protocol.Proto, error) {
	proto := protocol.NewProto(1)

	data := protocol.NewData()

	curTime := time.Now()
	currentTime := curTime.Local().Format("2006-01-02 15:04:05")

	data.Time = currentTime
	data.Tag["node_name"] = GlobalNodeInfo.Name
	data.Tag["node_ip"] = GlobalNodeInfo.IP
	data.Tag["os"] = GlobalHostInfo.OSName
	data.Tag["os_release"] = GlobalHostInfo.OSRelease
	data.Tag["os_version"] = GlobalHostInfo.OSVersion
	data.Tag["platform"] = GlobalHostInfo.Platform
	data.Tag["host_name"] = GlobalHostInfo.HostName
	data.Tag["ncpus"] = GlobalHostInfo.NCPUs
	data.Tag["max_cpus"] = GlobalHostInfo.MaxCPUs
	data.Tag["bitwidth"] = GlobalHostInfo.BitWidth
	data.Tag["type"] = "heartbeat"
	data.Field["value"] = 1

	proto.DataList = append(proto.DataList, *data)

	return proto, nil
}
