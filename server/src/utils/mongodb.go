package utils

import (
	"time"
	"errors"

	"github.com/DarkMetrix/monitor/server/src/protocol"

	//log "github.com/cihub/seelog"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBUtils struct {
	session *mgo.Session							//Mongo db session
	database *mgo.Database							//Mongo db database
}

//Init mongodb
func (db *MongoDBUtils) Init (address string, database string) error {
	var err error

	if len(address) == 0 || len(database) == 0 {
		return errors.New("Param error, address or database empty")
	}

	db.session, err = mgo.DialWithTimeout(address, time.Second * 5)

	if err != nil {
		return err
	}

	db.session.SetSyncTimeout(time.Second * 5)
	db.session.SetSocketTimeout(time.Second * 5)

	db.database = db.session.DB(database)

	return nil
}

//******************Node*****************//
//Add node
func (db *MongoDBUtils) AddNode(node protocol.NodeInMongo) error {
	collection := db.database.C("nodes")

	if len(node.Info.NodeIP) == 0 {
		return errors.New("Param Invalid! node ip is empty!")
	}

	_, err := collection.Upsert(bson.M{"info.node_ip":node.Info.NodeIP}, bson.M{"info":node.Info})

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Get node
func (db *MongoDBUtils) GetNode(ip string) ([]protocol.NodeInMongo, error) {
	collection := db.database.C("nodes")

	var nodes []protocol.NodeInMongo

	if ip == "all" {
		err := collection.Find(bson.M{}).All(&nodes)

		if err != nil {
			return nil, err
		}

	} else {
		err := collection.Find(bson.M{"info.node_ip":ip}).All(&nodes)

		if err != nil {
			return nil, err
		}
	}

	return nodes, nil
}

//Del node
func (db *MongoDBUtils) DelNode(ip string) error {
	collection := db.database.C("nodes")

	if ip == "all" {
		_, err := collection.RemoveAll(bson.M{})

		if err != nil {
			return err
		}
	} else {
		err := collection.Remove(bson.M{"info.node_ip":ip})

		if err != nil {
			return err
		}
	}

	return nil
}

//******************Application instance*****************//
//Add application
func (db *MongoDBUtils) AddApplicationInstance(instance protocol.ApplicationInstanceInMongo) error {
	collection := db.database.C("applications")

	if len(instance.Info.Name) == 0 {
		return errors.New("Param Invalid! node ip is empty!")
	}

	_, err := collection.Upsert(bson.M{"info.name":instance.Info.Name}, bson.M{"info":instance})

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Get application instance
func (db *MongoDBUtils) GetApplicationInstance(instance string) ([]protocol.ApplicationInstanceInMongo, error) {
	collection := db.database.C("applications")

	var instances []protocol.ApplicationInstanceInMongo

	if instance == "all" {
		err := collection.Find(bson.M{}).All(&instances)

		if err != nil {
			return nil, err
		}

	} else {
		err := collection.Find(bson.M{"info.name":instance}).All(&instances)

		if err != nil {
			return nil, err
		}
	}

	return instances, nil
}

//Del application instance
func (db *MongoDBUtils) DelApplicationInstance(instance string) error {
	collection := db.database.C("applications")

	if instance == "all" {
		_, err := collection.RemoveAll(bson.M{})

		if err != nil {
			return err
		}
	} else {
		err := collection.Remove(bson.M{"info.name":instance})

		if err != nil {
			return err
		}
	}

	return nil
}

//******************Application instance & node mapping*****************//
//Add application instance & node mapping
func (db *MongoDBUtils) AddApplicationInstanceNodeMapping(mapping protocol.ApplicationInstanceNodeMapping) error {
	collection := db.database.C("mapping")

	if len(mapping.Info.NodeIP) == 0 || len(mapping.Info.Instance) == 0 {
		return errors.New("Param Invalid! ip or instance is empty")
	}

	mapping.Info.Key = mapping.Info.Instance + "@" + mapping.Info.NodeIP

	_, err := collection.Upsert(bson.M{"info.key":mapping.Info.Key}, bson.M{"info":mapping.Info})

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Get application instance by ip
func (db *MongoDBUtils) GetApplicationInstanceNodeMappingByIP(ip string) ([]protocol.ApplicationInstanceNodeMapping, error) {
	collection := db.database.C("mapping")

	if len(ip) == 0 {
		return nil, errors.New("Param Invalid! ip is empty")
	}

	var mapping []protocol.ApplicationInstanceNodeMapping

	err := collection.Find(bson.M{"info.node_ip":ip}).All(&mapping)

	if err != nil {
		db.session.Refresh()
		return nil, err
	}

	return mapping, nil
}

//Get application instance by instance
func (db *MongoDBUtils) GetApplicationInstanceNodeMappingByInstance(instance string) ([]protocol.ApplicationInstanceNodeMapping, error) {
	collection := db.database.C("mapping")

	if len(instance) == 0 {
		return nil, errors.New("Param Invalid! instance is empty")
	}

	var mapping []protocol.ApplicationInstanceNodeMapping

	err := collection.Find(bson.M{"info.instance":instance}).All(&mapping)

	if err != nil {
		db.session.Refresh()
		return nil, err
	}

	return mapping, nil
}

//******************View*****************//

//******************User Setting*****************//