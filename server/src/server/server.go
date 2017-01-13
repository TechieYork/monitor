package server

import (
	"fmt"
	"io/ioutil"
	"encoding/json"

	"github.com/DarkMetrix/monitor/common/error_code"
	"github.com/DarkMetrix/monitor/server/src/config"
	"github.com/DarkMetrix/monitor/server/src/protocol"

	log "github.com/cihub/seelog"

	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb/client/v2"
)

type MonitorServer struct {
	Config config.Config
	InfluxdbClient client.Client
}

//New Config
func NewMonitorServer(config *config.Config) *MonitorServer {
	return &MonitorServer{
		Config:*config,
	}
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

func (server *MonitorServer) Run () error {
	var err error

	//Init influxdb
	err = server.InitInfluxdb()

	if err != nil {
		return err
	}

	//Begin serve
	http := gin.Default()

	http.POST("/monitor/get_nodes", server.GetNodes)
	http.POST("/monitor/get_node_instances", server.GetNodeInstances)
	http.POST("/monitor/get_node_metrix", server.GetNodeMetrix)
	http.POST("/monitor/get_application_instances", server.GetApplicationInstances)
	http.POST("/monitor/get_application_metrix", server.GetApplicationMetrix)

	err = http.Run(server.Config.Server.Address)

	if err != nil {
		return err
	}

	return nil
}

//Http interface get_nodes
func (server *MonitorServer) GetNodes(context *gin.Context) {
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
	var request protocol.GetNodesRequest

	err = json.Unmarshal(body, &request)

	if err != nil {
		log.Warn("Decode body failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ParamError, "desc":error_code.GetErrorString(error_code.ParamError)})
		return
	}

	//Get node list
	var query client.Query

	if request.IP == "all" {
		query = client.Query{
			Command: "SELECT * FROM node GROUP BY node_ip LIMIT 1",
			Database: server.Config.Influxdb.DBName,
		}
	} else {
		query = client.Query{
			Command: fmt.Sprintf("SELECT * FROM node WHERE node_ip = '%s' GROUP BY node_ip LIMIT 1", request.IP),
			Database: server.Config.Influxdb.DBName,
		}
	}

	response, err := server.InfluxdbClient.Query(query)

	if err != nil {
		log.Warn("Query show measurements failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
		return
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
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
						context.JSON(200, gin.H{"code":error_code.Fail, "desc":error_code.GetErrorString(error_code.Fail)})
						return
					}

					if len(serie.Columns) != len(serie.Values[0]) {
						log.Warn("Get tag value failed! error:columns and valuse len not match")
						context.JSON(200, gin.H{"code":error_code.Fail, "desc":error_code.GetErrorString(error_code.Fail)})
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

	context.JSON(200, gin.H{"code":error_code.Success, "desc":error_code.GetErrorString(error_code.Success), "data":gin.H{"node_list":nodes}})
}

//Http interface get_node_instances
func (server *MonitorServer) GetNodeInstances(context *gin.Context) {
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
	var request protocol.GetNodeInstancesRequest

	err = json.Unmarshal(body, &request)

	if err != nil {
		log.Warn("Decode body failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ParamError, "desc":error_code.GetErrorString(error_code.ParamError)})
		return
	}

	//Get all measurements
	query := client.Query{
		Command: "SHOW MEASUREMENTS",
		Database: server.Config.Influxdb.DBName,
	}

	response, err := server.InfluxdbClient.Query(query)

	if err != nil {
		log.Warn("Query show measurements failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
		return
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
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
				Database: server.Config.Influxdb.DBName,
			}
		} else {
			query = client.Query{
				Command: fmt.Sprintf("SHOW TAG VALUES FROM \"%s\" WITH KEY = \"instance\" WHERE node_ip = '%s'", key, request.IP),
				Database: server.Config.Influxdb.DBName,
			}
		}

		log.Info("Query string:", query.Database, "-> ", query.Command)

		response, err := server.InfluxdbClient.Query(query)

		if err != nil {
			log.Warn("Query show tag values failed! error:", err)
			context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
			return
		}

		if response.Error() != nil {
			log.Warn("Query response failed! error:", err)
			context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
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

	context.JSON(200, gin.H{"code":error_code.Success, "desc":error_code.GetErrorString(error_code.Success), "data":collection})
}

//Http interface get_node_metrix
func (server *MonitorServer) GetNodeMetrix(context *gin.Context) {
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
	var request protocol.GetNodeMetrixRequest

	err = json.Unmarshal(body, &request)

	if err != nil {
		log.Warn("Decode body failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ParamError, "desc":error_code.GetErrorString(error_code.ParamError)})
		return
	}

	//Get all measurements
	query := client.Query{
		Command: "SHOW MEASUREMENTS",
		Database: server.Config.Influxdb.DBName,
	}

	response, err := server.InfluxdbClient.Query(query)

	if err != nil {
		log.Warn("Query show measurements failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
		return
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", err)
		context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
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

	metrixes := make(map[string][]client.Result)

	//Get all metrix except application
	for measurement := range collection.Measurements {
		var method string

		switch measurement {
		case "cpu":
			method = "MEAN"
			break
		case "memory":
			method = "MEAN"
			break
		case "filesystem":
			method = "MAX"
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
			Command: fmt.Sprintf("SELECT %s(value) FROM %s WHERE time > now() - %s AND node_ip = '%s' GROUP BY time(%s), instance ORDER BY time desc",
				method, measurement, request.Time, request.IP, interval),
			Database: server.Config.Influxdb.DBName,
		}

		log.Info("Query string:", query.Database, "-> ", query.Command)

		response, err := server.InfluxdbClient.Query(query)

		if err != nil {
			log.Warn("Query select failed! error:", err)
			context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
			return
		}

		if response.Error() != nil {
			log.Warn("Query response failed! error:", err)
			context.JSON(200, gin.H{"code":error_code.ServerBusy, "desc":error_code.GetErrorString(error_code.ServerBusy)})
			return
		}

		metrixes[measurement] = response.Results
	}

	context.JSON(200, gin.H{"code":error_code.Success, "desc":error_code.GetErrorString(error_code.Success), "data":gin.H{"metrix":metrixes}})
}

//Http interface get_application_instances
func (server *MonitorServer) GetApplicationInstances(context *gin.Context) {
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

	context.JSON(200, gin.H{"code":error_code.Success, "desc":error_code.GetErrorString(error_code.Success), "data":gin.H{"node_list":"1"}})
}

//Http interface get_application_metrix
func (server *MonitorServer) GetApplicationMetrix(context *gin.Context) {
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

	context.JSON(200, gin.H{"code":error_code.Success, "desc":error_code.GetErrorString(error_code.Success), "data":gin.H{"node_list":"1"}})
}

