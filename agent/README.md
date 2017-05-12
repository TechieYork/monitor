# DarkMetrix monitor agent

## Introduction

DarkMetrix monitor agent is an application that collect informations about the host(e.g.:cpu usage, memory, network etc.) and application information which reported by the application itself over UDP, then send to backend services like influxdb, mongodb, nsq and so on periodically. 



## Interfaces

All input or output plugin needs to implement the following interfaces.

#### Input plugin

*Init function*

```go
//Params:
//    nodeInfo: gives you all the information about in "node" of config.json
//    config: all the config key-value configurations in each "input_plugin.config" of config.json
//Return:
//    error: error information, return nil for success
func Init(nodeInfo config.NodeInfo, config map[string]string) error {
    //Do your initial work here(eg:initial some libs etc.)
}
```

*Collect function*

```Go
//Params:
//Return:
//    *protocol.Proto: the information that collected
//    error: error information
func Collect()(*protocol.Proto, error) {
    //Do your collect work here(eg:get cpu usage etc.)
}
```

#### Output plugin

*Init function*

```Go
//Params:
//    nodeInfo: gives you all the information about in "node" of config.json
//    config: all the config key-value configurations in each "input_plugin.config" of config.json
//Return:
//    error: error information, return nil for success
func Init(nodeInfo config.NodeInfo, config map[string]string) error {
    //Do your initial work here(eg:initial connection to influxdb or nsq etc.)
}
```

*Send function*

```go
//Params:
//    *protocol.Proto: the information to send
//Return:
//    error: error information, return nil for success
func Send(proto *protocol.Proto) error {
    //Do your send work here(eg:send data to influxdb or nsq etc.)
}
```



## Configuration

##### config.json

```json
{
	"node":
	{
		"name":"test",
		"ip":"10.0.0.1",
		"transfer_queue":
		{
			"buffer_size":10000
		}
	},

	"input_plugin":
	[
		{
			"plugin_name": "node",
			"plugin_path": "../plugin/input/node.so",
			"duration": 10,
			"active":true,
			"config":
			{
			}
		},
		{
			"plugin_name": "cpu",
			"plugin_path": "../plugin/input/cpu.so",
			"duration": 10,
			"active":true,
			"config":
			{
			}
		},
		{
			"plugin_name": "memory",
			"plugin_path": "../plugin/input/memory.so",
			"duration": 10,
			"active":true
		},
		{
			"plugin_name": "filesystem",
			"plugin_path": "../plugin/input/filesystem.so",
			"duration": 10,
			"active":true,
			"config":
			{
				"include":"/dev/sda.*;/dev/mapper/centos-root.*"
			}
		},
		{
			"plugin_name": "net",
			"plugin_path": "../plugin/input/net.so",
			"duration": 10,
			"active":true
		},
		{
			"plugin_name": "page",
			"plugin_path": "../plugin/input/page.so",
			"duration": 10,
			"active":true
		},
		{
			"plugin_name": "process",
			"plugin_path": "../plugin/input/process.so",
			"duration": 10,
			"active":true
		},
		{
			"plugin_name": "interfaces",
			"plugin_path": "../plugin/input/interfaces.so",
			"duration": 10,
			"active":true,
			"config":
			{
			}
		},
		{
			"plugin_name": "application",
			"plugin_path": "../plugin/input/application.so",
			"duration": 10,
			"active":true,
			"config":
			{
				"udp_address":"127.0.0.1:5656",
				"unix_address":"/var/tmp/monitor.sock"
			}
		}
	],
	"output_plugin":
	[
		{
			"plugin_name": "console",
			"plugin_path": "../plugin/output/console.so",
			"active":true,
			"inputs":
			{
				"node":false,
				"cpu":false,
				"memory":false,
				"filesystem":true,
				"net":false,
				"page":false,
				"process":false,
				"interfaces":false,
				"application":false
			},
			"config":
			{
				"type":"stdout"
			}
		},
		{
			"plugin_name": "nsq",
			"plugin_path": "../plugin/output/nsq.so",
			"active":false,
			"inputs":
			{
				"node":true,
				"cpu":false,
				"memory":false,
				"filesystem":false,
				"net":false,
				"page":false,
				"process":false,
				"interfaces":false,
				"application":true
			},
			"config":
			{
				"nsqd_address":"172.16.101.128:4150",
				"topic":"dark_metrix_monitor"
			}
		},
		{
			"plugin_name": "mongodb",
			"plugin_path": "../plugin/output/mongodb.so",
			"active":false,
			"inputs":
			{
				"node":true
			},
			"config":
			{
				"mongodb_address":"172.16.101.128:27017",
				"db_name":"dark_metrix"
			}
		},
		{
			"plugin_name": "influxdb",
			"plugin_path": "../plugin/output/influxdb.so",
			"active":true,
			"inputs":
			{
				"node":true,
				"cpu":true,
				"memory":true,
				"filesystem":true,
				"net":true,
				"page":true,
				"process":true,
				"interfaces":true,
				"application":true
			},
			"config":
			{
				"influxdb_address":"http://172.16.101.128:8086",
				"db_name":"dark_metrix"
			}
		}
	]
}
```

- **node.name:** the name of the host, the agent won't get the host name, you need to specify youself.
- **node.ip:** the ip of the host.
- **node.transfer_queue:** the size of queue to buffer monitored information which waiting to send.



- **input_plugin.plugin_name:** the plugin name.
- **input_plugin.plugin_path:** the plugin **.so** file path.
- **input_plugin.duration:** the duration which the agent would call **Collect** function.
- **input_plugin.active:** *true* or *false* to activated or deactivated the plugin.
- **input_plugin.config:** the configuration for each plugin in key-value style(map[string] string).



- **output_plugin.plugin_name:** the plugin name.
- **output_plugin.plugin_path:** the plugin **.so** file path.
- **output_plugin.active:** *true* or *false* to activated or deactivated the plugin.
- **output_plugin.inputs:** indicate which input plugin's data will be sent to this output plugin.
- **output_plugin.config:** the configuration for each plugin in key-value style(map[string] string).



##### log.config

See [cihub/seelog](https://github.com/cihub/seelog) to get more information.

## Notice

When collecting application report, if there is high concurrency demand, udp's receive buffer is needed to set a bigger value(In linux you should set the kernel limitation, such as rmem_max etc.). In application plugin, the receive buffer is set to 16M.





