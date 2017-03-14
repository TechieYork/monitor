package server

import (
	"fmt"
	"time"

	"github.com/DarkMetrix/monitor/server/src/config"
	"github.com/DarkMetrix/monitor/server/src/protocol"

	log "github.com/cihub/seelog"

	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb/client/v2"
)

type MonitorServer struct {
	config config.Config                                //Config information
	influxdbClient client.Client                        //Influxdb client connection
}

//New Config
func NewMonitorServer(config *config.Config) *MonitorServer {
	return &MonitorServer{
		config:*config,
	}
}

//Init influxdb
func (server *MonitorServer) InitInfluxdb() error {
	var err error

	server.influxdbClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: server.config.Influxdb.Address,
	})

	if err != nil {
		return err
	}

	return nil
}

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
	var err error

	//Unmarshal json
	var request protocol.GetNodesRequest

	request.IP = context.Query("ip")

	if len(request.IP) == 0 {
		context.JSON(400, gin.H{"message":"param invalid!"})
		return
	}

	//Get node list
	var query client.Query

	if request.IP == "all" {
		query = client.Query{
			Command: "SELECT * FROM node GROUP BY node_ip ORDER BY time desc LIMIT 1",
			Database: server.config.Influxdb.DBName,
		}
	} else {
		query = client.Query{
			Command: fmt.Sprintf("SELECT * FROM node WHERE node_ip = '%s' GROUP BY node_ip ORDER BY time desc LIMIT 1", request.IP),
			Database: server.config.Influxdb.DBName,
		}
	}

	response, err := server.influxdbClient.Query(query)

	if err != nil {
		log.Warn("Query show measurements failed! error:", err)
		context.JSON(500, gin.H{"message":err.Error()})
		return
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", response.Error().Error())
		context.JSON(500, gin.H{"message":response.Error().Error()})
		return
	}

	nodes := []protocol.Node{}

	for _, result := range response.Results {
		for _, serie := range result.Series {
			if serie.Name == "node" {
				node := protocol.NewNode()

				for index, val := range serie.Columns {
					column := fmt.Sprintf("%s", val)

					if len(serie.Values) <= 0 {
						log.Warn("Get tag value failed! error:val array len is 0")
						context.JSON(404, gin.H{"message":"serie not found!"})
						return
					}

					if len(serie.Columns) != len(serie.Values[0]) {
						log.Warn("Get tag value failed! error:columns and values len not match")
						context.JSON(500, gin.H{"message":"columns and values not match!"})
						return
					}

					value := fmt.Sprintf("%s", serie.Values[0][index])

					node.Info[column] = value
				}

				for key, value := range serie.Tags {
					node.Info[key] = value
				}

				nodes = append(nodes, *node)
			}
		}
	}

	context.JSON(200, gin.H{"message":"success",
		"data":gin.H{"nodes":nodes, "server_time":time.Now().UTC().Format("2006-01-02T15:04:05.000000000Z")}})
}

//Http interface get_node_instances
func (server *MonitorServer) GetNodeInstances(context *gin.Context) {
	var err error

	//Unmarshal json
	var request protocol.GetNodeInstancesRequest

	request.IP = context.Query("ip")

	if len(request.IP) == 0 {
		context.JSON(400, gin.H{"message":"param invalid!"})
		return
	}

	//Get all measurements
	query := client.Query{
		Command: "SHOW MEASUREMENTS",
		Database: server.config.Influxdb.DBName,
	}

	response, err := server.influxdbClient.Query(query)

	if err != nil {
		log.Warn("Query show measurements failed! error:", err)
		context.JSON(500, gin.H{"message":err.Error()})
		return
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", response.Error().Error())
		context.JSON(500, gin.H{"message":response.Error().Error()})
		return
	}

	//Get all instances except application
	collection := protocol.NewNodeInstance()

	for _, result := range response.Results {
		for _, serie := range result.Series {
			if serie.Name == "measurements" {
				for _, val := range serie.Values {
					for _, m := range val {
						measurement := fmt.Sprintf("%s", m)

						if measurement== "application" {
							continue
						}

						collection.Measurements[measurement] = []string{}
					}
				}
			}
		}
	}

	for key := range collection.Measurements {
		//Get node list
		var query client.Query

		if request.IP == "all" {
			query = client.Query{
				Command: fmt.Sprintf("SHOW TAG VALUES FROM \"%s\" WITH KEY = \"instance\"", key),
				Database: server.config.Influxdb.DBName,
			}
		} else {
			query = client.Query{
				Command: fmt.Sprintf("SHOW TAG VALUES FROM \"%s\" WITH KEY = \"instance\" WHERE node_ip = '%s'", key, request.IP),
				Database: server.config.Influxdb.DBName,
			}
		}

		log.Info("Query string:", query.Database, "-> ", query.Command)

		response, err := server.influxdbClient.Query(query)

		if err != nil {
			log.Warn("Query show tag values failed! error:", err)
			context.JSON(500, gin.H{"message":err.Error()})
			return
		}

		if response.Error() != nil {
			log.Warn("Query response failed! error:", response.Error().Error())
			context.JSON(500, gin.H{"message":response.Error().Error()})
			return
		}

		for _, result := range response.Results {
			for _, serie := range result.Series {
				if serie.Name != "application" {
					for _, value := range serie.Values {
						for _, inst := range value {
							instance := fmt.Sprintf("%s", inst)

							if instance == "instance" {
								continue
							}

							collection.Measurements[key] = append(collection.Measurements[key], instance)
						}
					}
				}
			}
		}
	}

	context.JSON(200, gin.H{"message":"success", "data":collection})
}

//Http interface get_node_metrix
func (server *MonitorServer) GetNodeMetrix(context *gin.Context) {
	var err error

	//Unmarshal json
	var request protocol.GetNodeMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")

	if len(request.IP) == 0 || len(request.Time) == 0 {
		context.JSON(400, gin.H{"message":"param invalid!"})
		return
	}

	//Get all measurements
	query := client.Query{
		Command: "SHOW MEASUREMENTS",
		Database: server.config.Influxdb.DBName,
	}

	response, err := server.influxdbClient.Query(query)

	if err != nil {
		log.Warn("Query show measurements failed! error:", err)
		context.JSON(500, gin.H{"message":err.Error()})
		return
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", err)
		context.JSON(500, gin.H{"message":response.Error().Error()})
		return
	}

	//Get all instances except application
	collection := protocol.NewNodeInstance()

	for _, result := range response.Results {
		for _, serie := range result.Series {
			if serie.Name == "measurements" {
				for _, val := range serie.Values {
					for _, m := range val {
						measurement := fmt.Sprintf("%s", m)

						if measurement == "application" {
							continue
						}

						collection.Measurements[measurement] = []string{}
					}
				}
			}
		}
	}

	metrixes := make(map[string][]client.Result)

	//Get all metrix except application
	for measurement := range collection.Measurements {
		var method string
		var groupby string
		groupby = ""

		switch measurement {
		case "cpu":
			method = "MEAN"
			break
		case "memory":
			method = "MEAN"
			break
		case "filesystem":
			method = "MAX"
			groupby = ", mount_point, device_name "
			break
		case "net":
			method = "SUM"
			break
		case "page":
			method = "SUM"
			break
		case "process":
			method = "MAX"
			break
		case "interfaces":
			groupby = ", interface "
			method = "MEAN"
			break
		default:
			method = "SUM"
		}

		var interval string

		switch request.Time {
		case "1h":
			interval = "1m"
		case "1d":
			interval = "1m"
		case "7d":
			interval = "5m"
		case "30d":
			interval = "10m"
		case "90d":
			interval = "30m"
		default:
			interval = "5m"
		}

		query := client.Query{
			Command: fmt.Sprintf("SELECT %s(value) FROM %s WHERE time > now() - %s AND node_ip = '%s' GROUP BY time(%s), instance%s ORDER BY time desc",
				method, measurement, request.Time, request.IP, interval, groupby),
			Database: server.config.Influxdb.DBName,
		}

		log.Info("Query string:", query.Database, "-> ", query.Command)

		response, err := server.influxdbClient.Query(query)

		if err != nil {
			log.Warn("Query select failed! error:", err)
			context.JSON(500, gin.H{"message":err.Error()})
			return
		}

		if response.Error() != nil {
			log.Warn("Query response failed! error:", err)
			context.JSON(500, gin.H{"message":response.Error().Error()})
			return
		}

		metrixes[measurement] = response.Results
	}

	context.JSON(200, gin.H{"message":"success", "data":gin.H{"metrix":metrixes}})
}

//Http interface get_application_instances
func (server *MonitorServer) GetApplicationInstances(context *gin.Context) {
	var err error

	//Unmarshal json
	var request protocol.GetApplicationInstancesRequest

	request.IP = context.Query("ip")

	if len(request.IP) == 0 {
		context.JSON(400, gin.H{"message":"param invalid!"})
		return
	}

	//Get node list
	var query client.Query

	if request.IP == "all" {
		query = client.Query{
			Command: "SHOW TAG VALUES FROM application WITH KEY = \"instance\"",
			Database: server.config.Influxdb.DBName,
		}
	} else {
		query = client.Query{
			Command: fmt.Sprintf("SHOW TAG VALUES FROM application WITH KEY = \"instance\" WHERE node_ip = '%s'", request.IP),
			Database: server.config.Influxdb.DBName,
		}
	}

	log.Info("Query string:", query.Database, "-> ", query.Command)

	response, err := server.influxdbClient.Query(query)

	if err != nil {
		log.Warn("Query show tag values failed! error:", err)
		context.JSON(500, gin.H{"message":err.Error()})
		return
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", response.Error().Error())
		context.JSON(500, gin.H{"message":response.Error().Error()})
		return
	}

	//Get all instances except application
	collection := protocol.NewNodeInstance()

	for _, result := range response.Results {
		for _, serie := range result.Series {
			if serie.Name == "application" {
				for _, value := range serie.Values {
					for _, inst := range value {
						instance := fmt.Sprintf("%s", inst)

						if instance == "instance" {
							continue
						}

						collection.Measurements["application"] = append(collection.Measurements["application"], instance)
					}
				}
			}
		}
	}

	context.JSON(200, gin.H{"message":"success", "data":collection})
}

//Http interface get_application_metrix
func (server *MonitorServer) GetApplicationMetrix(context *gin.Context) {
	var err error

	//Unmarshal json
	var request protocol.GetApplicationMetrixRequest

	request.IP = context.Query("ip")
	request.Time = context.Query("time")
	request.Instance = context.Query("instance")

	if len(request.Time) == 0 || len(request.IP) == 0 || len(request.Instance) == 0{
		context.JSON(400, gin.H{"message":"param invalid!"})
		return
	}

	//Get all instances except application
	metrixes := make(map[string][]client.Result)

	//Get all metrix except application
	var interval string

	switch request.Time {
	case "1h":
		interval = "1m"
	case "1d":
		interval = "1m"
	case "7d":
		interval = "5m"
	case "30d":
		interval = "10m"
	case "90d":
		interval = "30m"
	default:
		interval = "5m"
	}

	query := client.Query{}

	if request.IP == "all" {
		if request.Instance == "all" {
			query = client.Query{
				Command: fmt.Sprintf("SELECT SUM(value) FROM application WHERE time > now() - %s GROUP BY time(%s), instance ORDER BY time desc",
					request.Time, interval),
				Database: server.config.Influxdb.DBName,
			}
		} else {
			query = client.Query{
				Command: fmt.Sprintf("SELECT SUM(value) FROM application WHERE time > now() - %s AND instance = '%s' GROUP BY time(%s), instance ORDER BY time desc",
					request.Time, request.Instance, interval),
				Database: server.config.Influxdb.DBName,
			}
		}

	} else {
		if request.Instance == "all" {
			query = client.Query{
				Command: fmt.Sprintf("SELECT SUM(value) FROM application WHERE time > now() - %s AND node_ip = '%s' GROUP BY time(%s), instance ORDER BY time desc",
					request.Time, request.IP, interval),
				Database: server.config.Influxdb.DBName,
			}
		} else {
			query = client.Query{
				Command: fmt.Sprintf("SELECT SUM(value) FROM application WHERE time > now() - %s AND node_ip = '%s' AND instance = '%s' GROUP BY time(%s), instance ORDER BY time desc",
					request.Time, request.IP, request.Instance, interval),
				Database: server.config.Influxdb.DBName,
			}
		}
	}

	log.Info("Query string:", query.Database, "-> ", query.Command)

	response, err := server.influxdbClient.Query(query)

	if err != nil {
		log.Warn("Query select failed! error:", err)
		context.JSON(500, gin.H{"message":err.Error()})
		return
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", err)
		context.JSON(500, gin.H{"message":response.Error().Error()})
		return
	}

	metrixes["application"] = response.Results

	context.JSON(200, gin.H{"message":"success", "data":gin.H{"metrix":metrixes}})
}

