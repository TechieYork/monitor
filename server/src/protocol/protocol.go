package protocol

//Transfer queue information
type TransferQueueInfo struct{
	BufferSize int `json:"buffer_size" bson:"buffer_size"`
}

//Node information
type NodeInfo struct {
	Name string `json:"name" bson:"name"`
	IP string `json:"ip" bson:"ip"`
	TransferQueue TransferQueueInfo `json:"transfer_queue" bson:"transfer_queue"`
}

//Admin information
type AdminInfo struct {
	Address string `json:"addr" bson:"addr"`
}

//Registry information
type RegistryInfo struct {
	Address string `json:"addr" bson:"addr"`
}

//Input plugin information
type InputPluginInfo struct {
	Name string `json:"plugin_name" bson:"plugin_name"`
	Path string `json:"plugin_path" bson:"plugin_path"`
	Duration int `json:"duration" bson:"duration"`
	Active bool `json:"active" bson:"active"`
	PluginConfig map[string]string `json:"config" bson:"config"`
}

//Output plugin information
type OutputPluginInfo struct {
	Name string `json:"plugin_name" bson:"plugin_name"`
	Path string `json:"plugin_path" bson:"plugin_path"`
	Active bool `json:"active" bson:"active"`
	Inputs map[string]bool `json:"inputs" bson:"inputs"`
	PluginConfig map[string]string `json:"config" bson:"config"`
}

//Node config sturcture
type NodeConfig struct {
	Node NodeInfo `json:"node" bson:"node"`
	Admin AdminInfo `json:"admin" bson:"admin"`
	Registry RegistryInfo `json:"registry" bson:"registry"`
	Inputs []InputPluginInfo `json:"input_plugin" bson:"input_plugin"`
	Outputs []OutputPluginInfo `json:"output_plugin" bson:"output_plugin"`
}

//New Config
func NewNodeConfig() *NodeConfig {
	return &NodeConfig{
		Node:NodeInfo{Name:"unknown", TransferQueue:TransferQueueInfo{BufferSize:1000}},
		Admin:AdminInfo{Address:""},
		Registry:RegistryInfo{Address:""},
		Inputs:[]InputPluginInfo{},
		Outputs:[]OutputPluginInfo{},
	}
}
