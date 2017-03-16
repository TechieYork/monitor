# DarkMetrix monitor server

## Introduction

DarkMetrix monitor server is an application that read influxdb and show all the collect data on web.



## Configuration

##### config.json

```json
{
	"server": {
		"addr":":9111"
	},
	"influxdb":{
		"addr":"http://172.16.101.128:8086",
		"db_name":"dark_metrix"
	}
}
```

- **server.addr:** the address that the server to listen.
- **influxdb.addr:** the address of influxdb
- **influxdb.db_name:** the name of influxdb which monitor agent used to save data.



##### log.config

See [cihub/seelog](https://github.com/cihub/seelog) to get more information.