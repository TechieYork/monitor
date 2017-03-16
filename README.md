# DarkMetrix Monitor

A monitor system that collect host information(e.g.:cpu usage, memory, network etc.) and application information which reported by the application itself over UDP, then send to backend services, such as influxdb, mongodb, nsq and so on periodically. And if the data is send to influxdb, the server will supply web and server to show all the information.



## Architecture

![image](https://github.com/DarkMetrix/monitor/blob/master/doc/architecture.png)

## Quick start

#### Preperation

There are two parts, agent and server. The agent will work fine without server. The server only works with influxdb and supply a couple of http interfaces for web pages. If influxdb is used as one of the output, Grafana could also used to show the information.

A **golang version 1.8** and above is needed.

#### Install from source

```shell
git clone git@github.com:DarkMetrix/monitor.git 
```

#### Build agent & run

```shell
$cd DarkMetrix/monitor/agent/src
$go build -o ../bin/dm_monitor_agent

$cd DarkMetrix/monitor/agent/plugin/input
$go build -buildmode=plugin node.go
$go build -buildmode=plugin cpu.go
$go build -buildmode=plugin filesystem.go
$go build -buildmode=plugin interfaces.go
$go build -buildmode=plugin memory.go
$go build -buildmode=plugin net.go
$go build -buildmode=plugin process.go
$go build -buildmode=plugin page.go
$go build -buildmode=plugin application.go

$cd DarkMetrix/monitor/agent/plugin/output
$go build -buildmode=plugin console.go
$go build -buildmode=plugin nsq.go
$go build -buildmode=plugin mongodb.go
$go build -buildmode=plugin influxdb.go

$../../admin/start.sh
```

#### Build server & run

```shell
$cd DarkMetrix/monitor/server/src
$go build -o ../bin/dm_log_server
$../admin/start.sh
```



## Configuration

All configuration file is in json.

#### agent

See [DarkMetrix/monitor/agent](https://github.com/DarkMetrix/monitor/blob/master/agent/README.md) to get more information.

#### server

See [DarkMetrix/monitor/server](https://github.com/DarkMetrix/monitor/blob/master/server/README.md) to get more information.



## TODO

* Monitor server's web
* Alarm system
* Some other monitor agent plugin, such as nginx, mysql etc.



## Lisense

#### DarkMetrix Monitor

MIT license

#### Dependencies

* github.com/cihub/seelog [BSD License](https://github.com/cihub/seelog/blob/master/LICENSE.txt)
* github.com/akhenakh/statgo [MIT License](https://github.com/akhenakh/statgo/blob/master/LICENSE)
* github.com/influxdata/influxdb/client/v2 [BSD License](https://github.com/influxdata/influxdb/blob/master/LICENSE)
* gopkg.in/mgo.v2 [MIT License](https://github.com/go-mgo/mgo/blob/v2/LICENSE)
* github.com/nsqio/go-nsq [MIT License](https://github.com/nsqio/go-nsq/blob/master/LICENSE)
* github.com/gin-gonic/gin [MIT License](https://github.com/gin-gonic/gin/blob/master/LICENSE)
* github.com/spf13/viper [MIT License](https://github.com/spf13/viper/blob/master/LICENSE)
