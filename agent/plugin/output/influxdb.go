package main

import(
	"fmt"
	"time"
	"errors"

	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/protocol"

	"github.com/influxdata/influxdb/client/v2"
)

var GlobalNodeInfo config.NodeInfo
var GlobalConfig map[string]string

var GlobalInfluxDB client.Client
var GlobalDBName string

func Init(nodeInfo config.NodeInfo, config map[string]string) error {
	var err error

	_, ok := config["influxdb_address"]

	if !ok {
		return errors.New("Missing config 'influxdb_address'")
	}

	GlobalDBName, ok = config["db_name"]

	if !ok {
		return errors.New("Missing config 'db_name'")
	}

	GlobalInfluxDB, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: config["influxdb_address"],
	})

	if err != nil {
		return err
	}

	GlobalConfig = make(map[string]string)
	GlobalConfig = config
	GlobalNodeInfo = nodeInfo

	return nil
}

func Send(proto *protocol.Proto) error {
	batchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: GlobalDBName,
	})

	if err != nil {
		return err
	}

	for _, data := range proto.DataList {
		tags := make(map[string]string)

		curTime, err := time.Parse("2006-01-02 15:04:05", data.Time)

		if err != nil {
			continue
		}

		for key, value := range data.Tag {
			tags[key] = fmt.Sprintf("%v", value)
		}

		for key, value := range data.Field {
			tags["instance"] = key
			field := map[string]interface{}{"value":value}

			point, err := client.NewPoint(proto.Name, tags, field, curTime)

			if err != nil {
				continue
			}

			batchPoints.AddPoint(point)
		}
	}

	err = GlobalInfluxDB.Write(batchPoints)

	if err != nil {
		return err
	}

	return nil
}
