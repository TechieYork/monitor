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
	pageStat := GlobalStat.PageStats()

	if pageStat == nil {
		return nil, errors.New("PageStats failed! error:got nil")
	}

	proto := protocol.NewProto(1)

	curTime := time.Now()
	currentTime := curTime.Local().Format("2006-01-02 15:04:05")

	data := protocol.NewData()
	data.Time = currentTime

	data.Tag["node_name"] = GlobalNodeInfo.Name
	data.Tag["node_ip"] = GlobalNodeInfo.IP
	data.Field["page_in"] = pageStat.PageIn
	data.Field["page_out"] = pageStat.PageOut

	proto.DataList = append(proto.DataList, *data)

	return proto, nil
}
