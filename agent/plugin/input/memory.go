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
	memoryStat := GlobalStat.MemStats()

	if memoryStat == nil {
		return nil, errors.New("MemStats failed! error:got nil")
	}

	proto := protocol.NewProto(1)

	curTime := time.Now()
	currentTime := curTime.Local().Format("2006-01-02 15:04:05")

	data := protocol.NewData()
	data.Time = currentTime

	data.Tag["node_name"] = GlobalNodeInfo.Name
	data.Tag["node_ip"] = GlobalNodeInfo.IP
	data.Field["total"] = memoryStat.Total
	data.Field["free"] = memoryStat.Free
	data.Field["used"] = memoryStat.Used
	data.Field["cache"] = memoryStat.Cache
	data.Field["swap_total"] = memoryStat.SwapTotal
	data.Field["swap_free"] = memoryStat.SwapFree
	data.Field["swap_used"] = memoryStat.SwapUsed

	proto.DataList = append(proto.DataList, *data)

	return proto, nil
}
