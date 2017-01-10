package server

import (
	"time"
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/DarkMetrix/monitor/common/error_code"
	"github.com/DarkMetrix/monitor/server/src/config"
	"github.com/DarkMetrix/monitor/server/src/protocol"

	log "github.com/cihub/seelog"

	"github.com/gin-gonic/gin"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/nsqio/go-nsq"
)

type MonitorServer struct {
	Config config.Config
	MongodbClient *mgo.Session
	InfluxdbClient client.Client
	NsqConsumer *nsq.Consumer

	//Application info
	PointMap map[string]string
	Lock sync.Mutex
}

//New Config
func NewMonitorServer(config *config.Config) *MonitorServer {
	return &MonitorServer{
		Config:*config,
		PointMap:make(map[string]string),
	}
}

//Init mongodb
func (server *MonitorServer) InitMongodb() error {
	var err error

	server.MongodbClient, err = mgo.DialWithTimeout(server.Config.Mongodb.Address, time.Second * 5)

	if err != nil {
		return err
	}

	server.MongodbClient.SetSyncTimeout(time.Second * 5)
	server.MongodbClient.SetSocketTimeout(time.Second * 5)

	return nil
}

//Init influxdb
func (server *MonitorServer) InitInfluxdb() error {
	var err error

	server.InfluxdbClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: server.Config.Influxdb.Address,
	})

	if err != nil {
		return err
	}

	return nil
}

//Nsq consumer
func (server *MonitorServer) HandleMessage(message *nsq.Message) error {
	//log.Info("Message:", string(message.Body))

	var err error

	//Unmarshal json
	proto := protocol.NewProto(1)
	err = json.Unmarshal(message.Body, proto)

	if err != nil {
		log.Warn("Decode body failed! error:", err)
		return err
	}

	switch proto.Name {
	case "node":
		server.HandleNodeMessage(proto)
	case "application":
		server.HandleApplicationMessage(proto)
	}

	return nil
}

//Init nsq
func (server *MonitorServer) InitNsq() error {
	var err error

	//Init nsq producer
	server.NsqConsumer, err = nsq.NewConsumer(server.Config.Nsq.TopicName, "monitor_server", nsq.NewConfig())

	if err != nil {
		return err
	}

	server.NsqConsumer.AddHandler(server)

	err = server.NsqConsumer.ConnectToNSQD(server.Config.Nsq.Address)

	if err != nil {
		return err
	}

	go server.ApplicationMessageProcessor()

	return nil
}

//Handle node message
func (server *MonitorServer) HandleNodeMessage (proto *protocol.Proto) {
	var err error

	//Update node information in mongodb
	db := server.MongodbClient.DB(server.Config.Mongodb.DBName)
	collection := db.C("node")

	for _, data := range proto.DataList {
		err = collection.Update(
			bson.M{"node_config.node.name":data.Tag["node_name"], "node_config.node.ip":data.Tag["node_ip"]},
			bson.M{"$set":bson.M{"node_info":data}})

		if err != nil {
			log.Warn("Update to mongodb failed! error:", err)
			server.MongodbClient.Refresh()

			continue
		}
	}
}

//Handle application message
func (server *MonitorServer) HandleApplicationMessage (proto *protocol.Proto) {
	//Update application node information in mongodb
	server.Lock.Lock()
	defer server.Lock.Unlock()

	curTime := time.Now()
	currentTime := curTime.Local().Format("2006-01-02 15:04:05")

	for _, data := range proto.DataList {
		for key := range data.Field {
			server.PointMap[key] = currentTime
		}
	}
}

//Application message processor
func (server *MonitorServer) ApplicationMessageProcessor () {
	var err error

	//Upsert report point to mongodb every 5 seconds
	for {
		select {
		case <- time.After(time.Second * 5):
			server.Lock.Lock()

			db := server.MongodbClient.DB(server.Config.Mongodb.DBName)
			collection := db.C("application")

			for key, value := range server.PointMap {
				_, err = collection.Upsert(
					bson.M{"name":key},
					bson.M{"name":key, "info":bson.M{"update_time":value}})

				if err != nil {
					log.Warn("Update to mongodb failed! error:", err)
					server.MongodbClient.Refresh()

					continue
				}
			}

			server.PointMap = make(map[string]string)

			server.Lock.Unlock()
		}
	}
}

//Run
func (server *MonitorServer) Run () error {
	var err error

	//Init mongodb
	err = server.InitMongodb()

	if err != nil {
		return err
	}

	//Init influxdb
	err = server.InitInfluxdb()

	if err != nil {
		return err
	}

	//Init nsq
	err = server.InitNsq()

	if err != nil {
		return err
	}

	//Begin serve
	http := gin.Default()

	http.POST("/monitor/node_register", server.NodeRegister)
	http.POST("/monitor/get_node_list", server.GetNodeList)
	http.POST("/monitor/get_point_list", server.GetPointList)

	err = http.Run(server.Config.Server.Address)

	if err != nil {
		return err
	}

	return nil
}

func (server *MonitorServer) NodeRegister(context *gin.Context) {
	var err error

	//Get json body
	body, err := ioutil.ReadAll(context.Request.Body)

	if err != nil {
		log.Warn("Get request body failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ParamError, "desc":error_code.GetErrorString(error_code.ParamError)})
		return
	}

	log.Info("Request body:", string(body))

	//Unmarshal json
	node := protocol.NewNodeConfig()
	err = json.Unmarshal(body, node)

	if err != nil {
		log.Warn("Decode body failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ParamError, "desc":error_code.GetErrorString(error_code.ParamError)})
		return
	}

	//Upsert to mongo db
	db := server.MongodbClient.DB(server.Config.Mongodb.DBName)
	collection := db.C("node")

	curTime := time.Now()
	currentTime := curTime.Local().Format("2006-01-02 15:04:05")

	_, err = collection.Upsert(
		bson.M{"node_config.node.name":node.Node.Name, "node_config.node.ip":node.Node.IP},
		bson.M{"time":currentTime, "node_config":node})

	if err != nil {
		log.Warn("Upsert to mongodb failed! error:", err)
		server.MongodbClient.Refresh()

		context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
		return
	}

	context.JSON(200, gin.H{"code":error_code.Success, "desc":error_code.GetErrorString(error_code.Success)})
}

func (server *MonitorServer) GetNodeList(context *gin.Context) {
	var err error

	//Get json body
	body, err := ioutil.ReadAll(context.Request.Body)

	if err != nil {
		log.Warn("Get request body failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ParamError, "desc":error_code.GetErrorString(error_code.ParamError)})
		return
	}

	log.Info("Request body:", string(body))

	//Unmarshal json
	page := protocol.NewPageInfo()
	err = json.Unmarshal(body, page)

	if err != nil {
		log.Warn("Decode body failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ParamError, "desc":error_code.GetErrorString(error_code.ParamError)})
		return
	}

	//Get node list
	db := server.MongodbClient.DB(server.Config.Mongodb.DBName)
	collection := db.C("node")

	nodes := []protocol.Node{}

	if page.Begin == 0 && page.Number == 0 {
		err = collection.Find(&bson.M{}).All(&nodes)
	} else {
		err = collection.Find(&bson.M{}).Sort("node_ip").Skip(page.Begin * page.Number).Limit(page.Number).All(&nodes)
	}

	if err != nil {
		log.Warn("Find from mongodb failed! error:", err)
		server.MongodbClient.Refresh()

		context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
		return
	}

	context.JSON(200, gin.H{"code":error_code.Success, "desc":error_code.GetErrorString(error_code.Success), "data":gin.H{"node_list":nodes}})
}

func (server *MonitorServer) GetPointList(context *gin.Context) {
	var err error

	//Get json body
	body, err := ioutil.ReadAll(context.Request.Body)

	if err != nil {
		log.Warn("Get request body failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ParamError, "desc":error_code.GetErrorString(error_code.ParamError)})
		return
	}

	log.Info("Request body:", string(body))

	//Unmarshal json
	page := protocol.NewPageInfo()
	err = json.Unmarshal(body, page)

	if err != nil {
		log.Warn("Decode body failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ParamError, "desc":error_code.GetErrorString(error_code.ParamError)})
		return
	}

	//Get node list
	db := server.MongodbClient.DB(server.Config.Mongodb.DBName)
	collection := db.C("application")

	points := []protocol.Point{}

	if page.Begin == 0 && page.Number == 0 {
		err = collection.Find(&bson.M{}).All(&points)
	} else {
		err = collection.Find(&bson.M{}).Sort("name").Skip(page.Begin * page.Number).Limit(page.Number).All(&points)
	}

	if err != nil {
		log.Warn("Find from mongodb failed! error:", err)
		server.MongodbClient.Refresh()

		context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
		return
	}

	context.JSON(200, gin.H{"code":error_code.Success, "desc":error_code.GetErrorString(error_code.Success), "data":gin.H{"node_list":points}})
}
