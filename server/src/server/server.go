package server

import (
	"time"

	"github.com/DarkMetrix/monitor/server/src/config"
	"github.com/DarkMetrix/monitor/server/src/protocol"

	"github.com/gin-gonic/gin"
	"github.com/DarkMetrix/monitor/server/src/utils"
)

type MonitorServer struct {
	config config.Config                                //Config information
	influxdb utils.InfluxDBUtils						//Influxdb utils
}

//New Config
func NewMonitorServer(config *config.Config) *MonitorServer {
	return &MonitorServer{
		config:*config,
	}
}

//Init influxdb
func (server *MonitorServer) InitInfluxdb() error {
	err := server.influxdb.Init(server.config.Influxdb.Address, server.config.Influxdb.DBName)

	if err != nil {
		return err
	}

	return nil
}

//Run http server
func (server *MonitorServer) Run () error {
	var err error

	//Init influxdb
	err = server.InitInfluxdb()

	if err != nil {
		return err
	}

	//Begin serve
	http := gin.Default()

	http.GET("/monitor/nodes", server.GetNodes)
	http.GET("/monitor/nodes/instances", server.GetNodeInstances)
	http.GET("/monitor/nodes/metrix", server.GetNodeMetrix)
	http.GET("/monitor/application/instances", server.GetApplicationInstances)
	http.GET("/monitor/application/metrix", server.GetApplicationMetrix)

	err = http.Run(server.config.Server.Address)

	if err != nil {
		return err
	}

	return nil
}

//Http interface get_nodes
func (server *MonitorServer) GetNodes(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodesRequest

	request.IP = context.Query("ip")

	if len(request.IP) == 0 {
		context.JSON(400, gin.H{"message":"param invalid!"})
		return
	}

	//Get nodes informations
	nodes, err := server.influxdb.GetNodes(request.IP)

	if err != nil {
		context.JSON(500, gin.H{"message":err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message":"success",
		"data":gin.H{"nodes":nodes, "server_time":time.Now().UTC().Format("2006-01-02T15:04:05.000000000Z")}})
}

//Http interface get_node_instances
func (server *MonitorServer) GetNodeInstances(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodeInstancesRequest

	request.IP = context.Query("ip")

	if len(request.IP) == 0 {
		context.JSON(400, gin.H{"message":"param invalid!"})
		return
	}

	//Get node instances
	collection, err := server.influxdb.GetNodeInstances(request.IP)

	if err != nil {
		context.JSON(500, gin.H{"message":err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message":"success", "data":collection})
}

//Http interface get_node_metrix
func (server *MonitorServer) GetNodeMetrix(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodeMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")

	if len(request.IP) == 0 || len(request.Time) == 0 {
		context.JSON(400, gin.H{"message":"param invalid!"})
		return
	}

	//Get node metrix
	metrixes, err := server.influxdb.GetNodeMetrix(request.IP, request.Time)

	if err != nil {
		context.JSON(500, gin.H{"message":err.Error()})
		return
	}

	context.JSON(200, gin.H{"message":"success", "data":gin.H{"metrix":metrixes}})
}

//Http interface get_application_instances
func (server *MonitorServer) GetApplicationInstances(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetApplicationInstancesRequest

	request.IP = context.Query("ip")

	if len(request.IP) == 0 {
		context.JSON(400, gin.H{"message":"param invalid!"})
		return
	}

	//Get application instances
	collection, err := server.influxdb.GetApplicationInstances(request.IP)

	if err != nil {
		context.JSON(500, gin.H{"message":err.Error()})
		return
	}

	context.JSON(200, gin.H{"message":"success", "data":collection})
}

//Http interface get_application_metrix
func (server *MonitorServer) GetApplicationMetrix(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetApplicationMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")
	request.Instance = context.Query("instance")

	if len(request.Time) == 0 || len(request.IP) == 0 || len(request.Instance) == 0{
		context.JSON(400, gin.H{"message":"param invalid!"})
		return
	}

	//Get application metrix
	metrixes, err := server.influxdb.GetApplicationMetrix(request.IP, request.Time, request.Instance)

	if err != nil {
		context.JSON(500, gin.H{"message":err.Error()})
		return
	}

	context.JSON(200, gin.H{"message":"success", "data":gin.H{"metrix":metrixes}})
}

