package utils

import (
	"fmt"
	"errors"

	"github.com/DarkMetrix/monitor/server/src/protocol"

	log "github.com/cihub/seelog"
	"github.com/influxdata/influxdb/client/v2"
)

type InfluxDBUtils struct {
	client client.Client                        //Influxdb client connection
	name string									//DB name
}

//Init influxdb
func (db *InfluxDBUtils) Init (address string, DBName string) error {
	var err error

	db.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: address,
	})

	if err != nil {
		return err
	}

	db.name = DBName

	return nil
}

//Get nodes information
func (db *InfluxDBUtils) GetNodes(ip string) ([]protocol.Node, error) {
	var err error

	//Get node list
	var query client.Query

	if ip == "all" {
		query = client.Query{
			Command: "SELECT * FROM node GROUP BY node_ip ORDER BY time desc LIMIT 1",
			Database: db.name,
		}
	} else {
		query = client.Query{
			Command: fmt.Sprintf("SELECT * FROM node WHERE node_ip = '%s' GROUP BY node_ip ORDER BY time desc LIMIT 1", ip),
			Database: db.name,
		}
	}

	response, err := db.client.Query(query)

	if err != nil {
		log.Warn("Query node failed! error:", err)
		return nil, err
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", response.Error().Error())
		return nil, response.Error()
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
						return nil, errors.New("serie not found!")
					}

					if len(serie.Columns) != len(serie.Values[0]) {
						log.Warn("Get tag value failed! error:columns and values len not match")
						return nil, errors.New("columns and values not match!")
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

	return nodes, nil
}

//Get node instances
func (db *InfluxDBUtils) GetNodeInstances(ip string) (*protocol.NodeInstance, error) {
	//Get all measurements
	query := client.Query{
		Command: "SHOW MEASUREMENTS",
		Database: db.name,
	}

	response, err := db.client.Query(query)

	if err != nil {
		log.Warn("Query show measurements failed! error:", err)
		return nil, err
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", response.Error().Error())
		return nil, response.Error()
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

		if ip == "all" {
			query = client.Query{
				Command: fmt.Sprintf("SHOW TAG VALUES FROM \"%s\" WITH KEY = \"instance\"", key),
				Database: db.name,
			}
		} else {
			query = client.Query{
				Command: fmt.Sprintf("SHOW TAG VALUES FROM \"%s\" WITH KEY = \"instance\" WHERE node_ip = '%s'", key, ip),
				Database: db.name,
			}
		}

		log.Info("Query string:", query.Database, "-> ", query.Command)

		response, err := db.client.Query(query)

		if err != nil {
			log.Warn("Query show tag values failed! error:", err)
			return nil, err
		}

		if response.Error() != nil {
			log.Warn("Query response failed! error:", response.Error().Error())
			return nil, response.Error()
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

	return collection, nil
}

//Get node metrix
func (db *InfluxDBUtils) GetNodeMetrix(ip string, time string) (map[string][]client.Result, error) {
	//Get all measurements
	query := client.Query{
		Command: "SHOW MEASUREMENTS",
		Database: db.name,
	}

	response, err := db.client.Query(query)

	if err != nil {
		log.Warn("Query show measurements failed! error:", err)
		return nil, err
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", err)
		return nil, response.Error()
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

		switch time {
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
				method, measurement, time, ip, interval, groupby),
			Database: db.name,
		}

		log.Info("Query string:", query.Database, "-> ", query.Command)

		response, err := db.client.Query(query)

		if err != nil {
			log.Warn("Query select failed! error:", err)
			return nil, err
		}

		if response.Error() != nil {
			log.Warn("Query response failed! error:", err)
			return nil, response.Error()
		}

		metrixes[measurement] = response.Results
	}

	return metrixes, nil
}

//Get application instances
func (db *InfluxDBUtils) GetApplicationInstances(ip string) (*protocol.ApplicationInstance, error) {
	var query client.Query

	if ip == "all" {
		query = client.Query{
			Command: "SHOW TAG VALUES FROM application WITH KEY = \"instance\"",
			Database: db.name,
		}
	} else {
		query = client.Query{
			Command: fmt.Sprintf("SHOW TAG VALUES FROM application WITH KEY = \"instance\" WHERE node_ip = '%s'", ip),
			Database: db.name,
		}
	}

	log.Info("Query string:", query.Database, "-> ", query.Command)

	response, err := db.client.Query(query)

	if err != nil {
		log.Warn("Query show tag values failed! error:", err)
		return nil, err
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", response.Error().Error())
		return nil, response.Error()
	}

	//Get all instances except application
	collection := protocol.NewApplicationInstance()

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

	return collection, nil
}

//Get application metrix
func (db *InfluxDBUtils) GetApplicationMetrix(ip string, time string, instance string) (map[string][]client.Result, error) {
	//Get all instances except application
	metrixes := make(map[string][]client.Result)

	//Get all metrix except application
	var interval string

	switch time {
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

	if ip == "all" {
		if instance == "all" {
			query = client.Query{
				Command: fmt.Sprintf("SELECT SUM(value) FROM application WHERE time > now() - %s GROUP BY time(%s), instance ORDER BY time desc",
					time, interval),
				Database: db.name,
			}
		} else {
			query = client.Query{
				Command: fmt.Sprintf("SELECT SUM(value) FROM application WHERE time > now() - %s AND instance = '%s' GROUP BY time(%s), instance ORDER BY time desc",
					time, instance, interval),
				Database: db.name,
			}
		}

	} else {
		if instance == "all" {
			query = client.Query{
				Command: fmt.Sprintf("SELECT SUM(value) FROM application WHERE time > now() - %s AND node_ip = '%s' GROUP BY time(%s), instance ORDER BY time desc",
					time, ip, interval),
				Database: db.name,
			}
		} else {
			query = client.Query{
				Command: fmt.Sprintf("SELECT SUM(value) FROM application WHERE time > now() - %s AND node_ip = '%s' AND instance = '%s' GROUP BY time(%s), instance ORDER BY time desc",
					time, ip, instance, interval),
				Database: db.name,
			}
		}
	}

	log.Info("Query string:", query.Database, "-> ", query.Command)

	response, err := db.client.Query(query)

	if err != nil {
		log.Warn("Query select failed! error:", err)
		return nil, err
	}

	if response.Error() != nil {
		log.Warn("Query response failed! error:", err)
		return nil, response.Error()
	}

	metrixes["application"] = response.Results

	return metrixes, nil
}
