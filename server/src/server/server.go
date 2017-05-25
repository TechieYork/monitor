package server

import (
	"time"

	"github.com/DarkMetrix/monitor/server/src/config"
	"github.com/DarkMetrix/monitor/server/src/protocol"

	log "github.com/cihub/seelog"

	"github.com/DarkMetrix/monitor/server/src/utils"
	"github.com/gin-gonic/gin"
)

type MonitorServer struct {
	config   config.Config       //Config information
	influxdb utils.InfluxDBUtils //Influxdb utils
	mongodb  utils.MongoDBUtils  //Mongodb utils
}

//New Config
func NewMonitorServer(config *config.Config) *MonitorServer {
	return &MonitorServer{
		config: *config,
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

func (server *MonitorServer) InitMongodb() error {
	err := server.mongodb.Init(server.config.Mongodb.Address, server.config.Mongodb.DBName)

	if err != nil {
		return err
	}

	return nil
}

//Run http server
func (server *MonitorServer) Run() error {
	var err error

	//Init influxdb
	err = server.InitInfluxdb()

	if err != nil {
		return err
	}

	//Init mongodb
	err = server.InitMongodb()

	if err != nil {
		return err
	}

	go server.SyncMetaDataToMongodb()

	//Begin serve
	http := gin.Default()

	http.GET("/monitor/nodes", server.GetNodes)
	http.GET("/monitor/nodes/metrix/cpu", server.GetNodeMerixCpu)
	http.GET("/monitor/nodes/metrix/memory", server.GetNodeMerixMemory)
	http.GET("/monitor/nodes/metrix/net", server.GetNodeMerixNet)
	http.GET("/monitor/nodes/metrix/page", server.GetNodeMerixPage)
	http.GET("/monitor/nodes/metrix/process", server.GetNodeMerixProcess)
	http.GET("/monitor/nodes/metrix/filesystem", server.GetNodeMerixFileSystem)
	http.GET("/monitor/nodes/metrix/interfaces", server.GetNodeMerixInterfaces)
	http.GET("/monitor/application/instances/mapping", server.GetApplicationInstancesMapping)
	http.GET("/monitor/application/nodes/mapping", server.GetNodesMapping)
	http.GET("/monitor/application/metrix", server.GetApplicationMetrix)

	http.POST("/monitor/view/project", server.AddProject)
	http.GET("/monitor/view/project", server.GetProject)
	http.PUT("/monitor/view/project", server.SetProject)
	http.DELETE("/monitor/view/project", server.DelProject)

	http.POST("/monitor/view/project/service", server.AddService)
	http.GET("/monitor/view/project/service", server.GetService)
	http.PUT("/monitor/view/project/service", server.SetService)
	http.DELETE("/monitor/view/project/service", server.DelService)

	http.POST("/monitor/view/project/service/module", server.AddModule)
	http.GET("/monitor/view/project/service/module", server.GetModule)
	http.PUT("/monitor/view/project/service/module", server.SetModule)
	http.DELETE("/monitor/view/project/service/module", server.DelModule)

	http.POST("/monitor/view/project/service/module/instance", server.AddApplicationInstanceToModule)
	http.PUT("/monitor/view/project/service/module/instance", server.DelApplicationInstanceFromModule)

	http.GET("/monitor/instance/search", server.SearchApplicationInstance)

	err = http.Run(server.config.Server.Address)

	if err != nil {
		return err
	}

	return nil
}

//Http interface /monitor/nodes
func (server *MonitorServer) GetNodes(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodesRequest

	request.IP = context.Query("ip")

	if len(request.IP) == 0 {
		context.JSON(400, gin.H{"message": "param invalid!"})
		return
	}

	hosts, err := server.mongodb.GetNode(request.IP)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message": "success", "data": gin.H{"nodes": hosts, "server_time": time.Now().UTC().Format("2006-01-02T15:04:05.000000000Z")}})
}

//Http interface /monitor/nodes/metrix/cpu
func (server *MonitorServer) GetNodeMerixCpu(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodeMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")

	if len(request.IP) == 0 || len(request.Time) == 0 {
		context.JSON(400, gin.H{"message": "param invalid!"})
		return
	}

	//Get node metrix
	metrixes, err := server.influxdb.GetNodeMetrixCpu(request.IP, request.Time)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "success", "data": gin.H{"metrix": metrixes}})
}

//Http interface /monitor/nodes/metrix/memory
func (server *MonitorServer) GetNodeMerixMemory(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodeMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")

	if len(request.IP) == 0 || len(request.Time) == 0 {
		context.JSON(400, gin.H{"message": "param invalid!"})
		return
	}

	//Get node metrix
	metrixes, err := server.influxdb.GetNodeMetrixMemory(request.IP, request.Time)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "success", "data": gin.H{"metrix": metrixes}})
}

//Http interface /monitor/nodes/metrix/net
func (server *MonitorServer) GetNodeMerixNet(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodeMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")

	if len(request.IP) == 0 || len(request.Time) == 0 {
		context.JSON(400, gin.H{"message": "param invalid!"})
		return
	}

	//Get node metrix
	metrixes, err := server.influxdb.GetNodeMetrixNet(request.IP, request.Time)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "success", "data": gin.H{"metrix": metrixes}})
}

//Http interface /monitor/nodes/metrix/page
func (server *MonitorServer) GetNodeMerixPage(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodeMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")

	if len(request.IP) == 0 || len(request.Time) == 0 {
		context.JSON(400, gin.H{"message": "param invalid!"})
		return
	}

	//Get node metrix
	metrixes, err := server.influxdb.GetNodeMetrixPage(request.IP, request.Time)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "success", "data": gin.H{"metrix": metrixes}})
}

//Http interface /monitor/nodes/metrix/process
func (server *MonitorServer) GetNodeMerixProcess(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodeMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")

	if len(request.IP) == 0 || len(request.Time) == 0 {
		context.JSON(400, gin.H{"message": "param invalid!"})
		return
	}

	//Get node metrix
	metrixes, err := server.influxdb.GetNodeMetrixProcess(request.IP, request.Time)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "success", "data": gin.H{"metrix": metrixes}})
}

//Http interface /monitor/nodes/metrix/filesystem
func (server *MonitorServer) GetNodeMerixFileSystem(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodeMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")

	if len(request.IP) == 0 || len(request.Time) == 0 {
		context.JSON(400, gin.H{"message": "param invalid!"})
		return
	}

	//Get node metrix
	metrixes, err := server.influxdb.GetNodeMetrixFileSystem(request.IP, request.Time)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "success", "data": gin.H{"metrix": metrixes}})
}

//Http interface /monitor/nodes/metrix/interfaces
func (server *MonitorServer) GetNodeMerixInterfaces(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodeMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")

	if len(request.IP) == 0 || len(request.Time) == 0 {
		context.JSON(400, gin.H{"message": "param invalid!"})
		return
	}

	//Get node metrix
	metrixes, err := server.influxdb.GetNodeMetrixInterfaces(request.IP, request.Time)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "success", "data": gin.H{"metrix": metrixes}})
}

//Http interface /monitor/application/instances/mapping
func (server *MonitorServer) GetApplicationInstancesMapping(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetApplicationInstancesMappingRequest

	request.Instance = context.Query("instance")

	if len(request.Instance) == 0 {
		context.JSON(400, gin.H{"message": "param invalid!"})
		return
	}

	//Get application instances
	mapping, err := server.mongodb.GetApplicationInstanceNodeMappingByInstance(request.Instance)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "success", "data": mapping})
}

//Http interface /monitor/application/nodes/mapping
func (server *MonitorServer) GetNodesMapping(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetNodesMappingRequest

	request.IP = context.Query("ip")

	if len(request.IP) == 0 {
		context.JSON(400, gin.H{"message": "param invalid!"})
		return
	}

	//Get application instances
	mapping, err := server.mongodb.GetApplicationInstanceNodeMappingByIP(request.IP)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "success", "data": mapping})
}

//Http interface get_application_metrix
func (server *MonitorServer) GetApplicationMetrix(context *gin.Context) {
	//Unmarshal json and check params
	var request protocol.GetApplicationMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")
	request.Instance = context.Query("instance")

	if len(request.Time) == 0 || len(request.IP) == 0 || len(request.Instance) == 0 {
		context.JSON(400, gin.H{"message": "param invalid!"})
		return
	}

	//Get application metrix
	metrixes, err := server.influxdb.GetApplicationMetrix(request.IP, request.Time, request.Instance)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "success", "data": gin.H{"metrix": metrixes}})
}

//Sync meta data from influxdb to mongodb
func (server *MonitorServer) SyncMetaDataToMongodb() {
	go server.SyncMetaDataNode(10 * time.Second)
	go server.SyncMetaDataApplication(10 * time.Second, "1m")
}

//Sync node info
func (server *MonitorServer) SyncMetaDataNode(duration time.Duration) {
	for {
		select {
		case <-time.After(duration):
			nodes, err := server.influxdb.GetNodes("all")

			if err != nil {
				log.Warn("GetNodes from influxdb failed! error:", err.Error())
				continue
			}

			for _, node := range nodes {
				var nodeInMongo protocol.Node

				nodeInMongo.NodeName = node.Info["node_name"]
				nodeInMongo.NodeIP = node.Info["node_ip"]
				nodeInMongo.HostName = node.Info["host_name"]
				nodeInMongo.Platform = node.Info["platform"]
				nodeInMongo.OS = node.Info["os"]
				nodeInMongo.OSVersion = node.Info["os_version"]
				nodeInMongo.OSRelease = node.Info["os_release"]

				nodeInMongo.MaxCPUs = node.Info["max_cpus"]
				nodeInMongo.NCPUs = node.Info["ncpus"]

				nodeInMongo.Bitwith = node.Info["bitwidth"]

				nodeInMongo.Time = node.Info["time"]

				err = server.mongodb.AddNode(nodeInMongo)

				if err != nil {
					log.Warn("AddNode to mongo db failed! error:", err.Error())
					continue
				}
			}
		}
	}
}

//Sync application info
func (server *MonitorServer) SyncMetaDataApplication(duration time.Duration, period string) {
	for {
		select {
		case <-time.After(duration):
			mapping, err := server.influxdb.GetApplicationInstancesNodeMapping(period)

			log.Info(mapping)

			if err != nil {
				log.Warn("GetApplicationInstancesNodeMapping from influxdb failed! error:", err.Error())
				continue
			}

			for ip, instance := range mapping {
				var instanceInMongo protocol.ApplicationInstance

				instanceInMongo.Name = instance

				//Add application instance information (Using upsert)
				err := server.mongodb.AddApplicationInstance(instanceInMongo)

				if err != nil {
					log.Warn("AddApplicationInstance to mongo db failed! error:", err.Error())
					continue
				}

				//Add instance and node mapping
				var instanceNodeMapping protocol.ApplicationInstanceNodeMapping

				instanceNodeMapping.NodeIP = ip
				instanceNodeMapping.Instance = instance

				err = server.mongodb.AddApplicationInstanceNodeMapping(instanceNodeMapping)

				if err != nil {
					log.Warn("AddApplicationInstanceNodeMapping to mongo db failed! error:", err.Error())
					continue
				}
			}
		}
	}
}

//Add project
func (server *MonitorServer) AddProject(context *gin.Context) {
	//Unmarshal json body and check params
	var project protocol.Project

	if context.BindJSON(&project) != nil {
		context.JSON(500, gin.H{"message": "json body invalid!"})
		return
	}

	err := server.mongodb.AddProject(project)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message": "success"})
}

//Get project
func (server *MonitorServer) GetProject(context *gin.Context) {
	//Unmarshal json and check params
	project := context.Query("project")

	projects, err := server.mongodb.GetProject(project)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
	}

	//Reply
	context.JSON(200, gin.H{"message": "success", "data": gin.H{"projects":projects}})
}

//Set project
func (server *MonitorServer) SetProject(context *gin.Context) {
	//Unmarshal json body and check params
	var project protocol.Project

	if context.BindJSON(&project) != nil {
		context.JSON(500, gin.H{"message": "json body invalid!"})
		return
	}

	err := server.mongodb.SetProject(project)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message": "success"})
}

//Del project
func (server *MonitorServer) DelProject(context *gin.Context) {
	//Unmarshal json body and check params
	project := context.Query("project")

	err := server.mongodb.DelProject(project)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
	}

	//Reply
	context.JSON(200, gin.H{"message": "success"})
}

//Add service
func (server *MonitorServer) AddService(context *gin.Context) {
	//Unmarshal json body and check params
	var service protocol.Service

	if context.BindJSON(&service) != nil {
		context.JSON(500, gin.H{"message": "json body invalid!"})
		return
	}

	err := server.mongodb.AddService(service.Project, service)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message": "success"})
}

//Get service
func (server *MonitorServer) GetService(context *gin.Context) {
	//Unmarshal json and check params
	project := context.Query("project")
	service := context.Query("service")

	services, err := server.mongodb.GetService(project, service)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
	}

	//Reply
	context.JSON(200, gin.H{"message": "success", "data": gin.H{"services":services}})
}

//Set service
func (server *MonitorServer) SetService(context *gin.Context) {
	//Unmarshal json body and check params
	var service protocol.Service

	if context.BindJSON(&service) != nil {
		context.JSON(500, gin.H{"message": "json body invalid!"})
		return
	}

	err := server.mongodb.SetService(service.Project, service)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message": "success"})
}

//Del service
func (server *MonitorServer) DelService(context *gin.Context) {
	//Unmarshal json body and check params
	project := context.Query("project")
	service := context.Query("service")

	err := server.mongodb.DelService(project, service)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
	}

	//Reply
	context.JSON(200, gin.H{"message": "success"})
}

//Add module
func (server *MonitorServer) AddModule(context *gin.Context) {
	//Unmarshal json body and check params
	var module protocol.Module

	if context.BindJSON(&module) != nil {
		context.JSON(500, gin.H{"message": "json body invalid!"})
		return
	}

	err := server.mongodb.AddModule(module.Project, module.Service, module)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message": "success"})
}

//Get module
func (server *MonitorServer) GetModule(context *gin.Context) {
	//Unmarshal json and check params
	project := context.Query("project")
	service := context.Query("service")
	module := context.Query("module")

	modules, err := server.mongodb.GetModule(project, service, module)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
	}

	//Reply
	context.JSON(200, gin.H{"message": "success", "data": gin.H{"modules":modules}})
}

//Set module
func (server *MonitorServer) SetModule(context *gin.Context) {
	//Unmarshal json body and check params
	var module protocol.Module

	if context.BindJSON(&module) != nil {
		context.JSON(500, gin.H{"message": "json body invalid!"})
		return
	}

	err := server.mongodb.SetModule(module.Project, module.Service, module)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message": "success"})
}

//Del module
func (server *MonitorServer) DelModule(context *gin.Context) {
	//Unmarshal json body and check params
	project := context.Query("project")
	service := context.Query("service")
	module := context.Query("module")

	err := server.mongodb.DelModule(project, service, module)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
	}

	//Reply
	context.JSON(200, gin.H{"message": "success"})
}

//Add application instance
func (server *MonitorServer) AddApplicationInstanceToModule(context *gin.Context) {
	//Unmarshal json body and check params
	var module protocol.Module

	if context.BindJSON(&module) != nil {
		context.JSON(500, gin.H{"message": "json body invalid!"})
		return
	}

	err := server.mongodb.AddApplicationInstancesToModule(module.Project, module.Service, module.Name, module.Instances)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message": "success"})
}

//Del application instance
func (server *MonitorServer) DelApplicationInstanceFromModule(context *gin.Context) {
	//Unmarshal json body and check params
	var module protocol.Module

	if context.BindJSON(&module) != nil {
		context.JSON(500, gin.H{"message": "json body invalid!"})
		return
	}

	err := server.mongodb.DelApplicationInstancesFromModule(module.Project, module.Service, module.Name, module.Instances)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message": "success"})
}

//Search instance
func (server *MonitorServer) SearchApplicationInstance(context *gin.Context) {
	//Unmarshal json body and check params
	text := context.Query("text")

	instances, err := server.mongodb.SearchApplicationInstance(text)

	if err != nil {
		context.JSON(500, gin.H{"message": err.Error()})
		return
	}

	//Reply
	context.JSON(200, gin.H{"message": "success", "data":instances})
}
