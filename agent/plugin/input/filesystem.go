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
	fsInfos := GlobalStat.FSInfos()

	if fsInfos == nil {
		return nil, errors.New("FSInfos failed! error:got nil")
	}

	proto := protocol.NewProto(1)

	curTime := time.Now()
	currentTime := curTime.Local().Format("2006-01-02 15:04:05")

	for _, info := range fsInfos {
		if len(GlobalIncludes) != 0 {
			_, ok := GlobalIncludes[info.MountPoint]

			if !ok {
				continue
			}
		}

		data := protocol.NewData()
		data.Time = currentTime

		data.Tag["node_name"] = GlobalNodeInfo.Name
		data.Tag["node_ip"] = GlobalNodeInfo.IP

		data.Tag["device_name"] = info.DeviceName
		data.Tag["fs_type"] = info.FSType
		data.Tag["mount_point"] = info.MountPoint

		data.Field["size"] = info.Size
		data.Field["used"] = info.Used
		data.Field["free"] = info.Free
		data.Field["available"] = info.Available
		data.Field["inodes_size"] = info.TotalInodes
		data.Field["inodes_used"] = info.UsedInodes
		data.Field["inodes_free"] = info.FreeInodes
		data.Field["inodes_available"] = info.AvailableInodes

		proto.DataList = append(proto.DataList, *data)
	}

	return proto, nil
}
