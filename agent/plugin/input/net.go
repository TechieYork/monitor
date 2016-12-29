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
	netStats := GlobalStat.NetIOStats()

	if netStats == nil {
		return nil, errors.New("NetIOStats failed! error:got nil")
	}

	proto := protocol.NewProto(1)

	curTime := time.Now()
	currentTime := curTime.Local().Format("2006-01-02 15:04:05")

	for _, info := range netStats{
		data := protocol.NewData()
		data.Time = currentTime

		data.Tag["node_name"] = GlobalNodeInfo.Name
		data.Tag["node_ip"] = GlobalNodeInfo.IP
		data.Tag["instance"] = info.IntName
		data.Field["tx"] = info.TX
		data.Field["rx"] = info.RX
		data.Field["ipackets"] = info.IPackets
		data.Field["opackets"] = info.OPackets
		data.Field["ierrors"] = info.IErrors
		data.Field["oerrors"] = info.OErrors
		data.Field["collisions"] = info.Collisions

		proto.DataList = append(proto.DataList, *data)
	}

	return proto, nil
}
