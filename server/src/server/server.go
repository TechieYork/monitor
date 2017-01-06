package server

import (
	"time"
	"encoding/json"

	"github.com/DarkMetrix/monitor/common/error_code"
	"github.com/DarkMetrix/monitor/server/src/config"
	"github.com/DarkMetrix/monitor/server/src/protocol"

	log "github.com/cihub/seelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/gin-gonic/gin"
)

type MonitorServer struct {
	Config config.Config
	MongodbClient *mgo.Session
	InfluxdbClient client.Client
}

//New Config
func NewMonitorServer(config *config.Config) *MonitorServer {
	return &MonitorServer{
		Config:*config,
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

	//Begin serve
	http := gin.Default()

	http.POST("/monitor/node_register", server.NodeRegister)

	err = http.Run(server.Config.Server.Address)

	if err != nil {
		return err
	}

	return nil
}

func (server *MonitorServer) NodeRegister(context *gin.Context) {
	log.Info("Register node:", context.Request.Body)

	var err error

	//Unmarshal json
	node := protocol.NewNodeConfig()

	decoder := json.NewDecoder(context.Request.Body)

	err = decoder.Decode(&node)

	log.Info("Node:", node)

	if err != nil {
		context.JSON(200, gin.H{"code":error_code.ParamError, "desc":error_code.GetErrorString(error_code.ParamError)})
		return
	}

	//Upsert to mongo db
	db := server.MongodbClient.DB(server.Config.Mongodb.DBName)
	collection := db.C("node")

	_, err = collection.Upsert(
		bson.M{"node_config.node.name":node.Node.Name, "node_config.node.ip":node.Node.IP},
		bson.M{"time":time.Now().String(), "node_config":node})

	if err != nil {
		server.MongodbClient.Refresh()

		context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
		return
	}

	context.JSON(200, gin.H{"code":error_code.Success, "desc":error_code.GetErrorString(error_code.Success)})
}

