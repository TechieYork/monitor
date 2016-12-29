package main

import(
	"time"
	"errors"

	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/protocol"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var GlobalNodeInfo config.NodeInfo
var GlobalConfig map[string]string

var GlobalSession *mgo.Session
var GlobalMongoDB *mgo.Database

func Init(nodeInfo config.NodeInfo, config map[string]string) error {
	var err error

	mongoAddress, ok := config["mongodb_address"]

	if !ok {
		return errors.New("Missing config 'mongodb_address'")
	}

	_, ok = config["db_name"]

	if !ok {
		return errors.New("Missing config 'db_name'")
	}

	GlobalSession, err = mgo.DialWithTimeout(mongoAddress, time.Second * 5)

	if err != nil {
		return err
	}

	GlobalSession.SetSyncTimeout(time.Second * 5)
	GlobalSession.SetSocketTimeout(time.Second * 5)

	GlobalConfig = make(map[string]string)
	GlobalConfig = config
	GlobalNodeInfo = nodeInfo

	GlobalMongoDB = GlobalSession.DB(GlobalConfig["db_name"])

	return nil
}

func Send(proto *protocol.Proto) error {
	collection := GlobalMongoDB.C(proto.Name)

	for _, data := range proto.DataList {
		//_, err := collection.Upsert(
		//	bson.M{"node":GlobalNodeInfo.Name, "ip":GlobalNodeInfo.IP},
		//	bson.M{"node":GlobalNodeInfo.Name, "ip":GlobalNodeInfo.IP, "info":data})

		err := collection.Insert(bson.M{"node":GlobalNodeInfo.Name, "ip":GlobalNodeInfo.IP, "info":data})

		if err != nil {
			GlobalSession.Refresh()
			return err
		}
	}

	return nil
}
