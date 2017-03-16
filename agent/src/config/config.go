package config

import (
	"github.com/spf13/viper"
)

const Version = "0.0.1"

//Transfer queue information
type TransferQueueInfo struct{
	BufferSize int `mapstructure:"buffer_size" json:"buffer_size"`
}

//Node information
type NodeInfo struct {
	Name string `mapstructure:"name" json:"name"`
	IP string `mapstructure:"ip" json:"ip"`
	TransferQueue TransferQueueInfo `mapstructure:"transfer_queue" json:"transfer_queue"`
}

//Input plugin information
type InputPluginInfo struct {
	Name string `mapstructure:"plugin_name" json:"plugin_name"`
	Path string `mapstructure:"plugin_path" json:"plugin_path"`
	Duration int `mapstructure:"duration" json:"duration"`
	Active bool `mapstructure:"active" json:"active"`
	PluginConfig map[string]string `mapstructure:"config" json:"config"`
}

//Output plugin information
type OutputPluginInfo struct {
	Name string `mapstructure:"plugin_name" json:"plugin_name"`
	Path string `mapstructure:"plugin_path" json:"plugin_path"`
	Active bool `mapstructure:"active" json:"active"`
	Inputs map[string]bool `mapstructure:"inputs" json:"inputs"`
	PluginConfig map[string]string `mapstructure:"config" json:"config"`
}

//Config sturcture
type Config struct {
	Node NodeInfo `mapstructure:"node" json:"node"`
	Inputs []InputPluginInfo `mapstructure:"input_plugin" json:"input_plugin"`
	Outputs []OutputPluginInfo `mapstructure:"output_plugin" json:"output_plugin"`
}

//Global config
var globalConfig *Config

//New Config
func NewConfig() *Config {
	return &Config{
		Node:NodeInfo{Name:"unknown", TransferQueue:TransferQueueInfo{BufferSize:1000}},
		Inputs:[]InputPluginInfo{},
		Outputs:[]OutputPluginInfo{},
	}
}

//Get singleton config
func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = NewConfig()
	}

	return globalConfig
}

//Init config from json file
func (config *Config) Init (path string) error {
	//Set viper setting
	viper.SetConfigType("json")
	viper.SetConfigFile(path)
	viper.AddConfigPath("../conf/")

	//Read in config
	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	//Unmarshal config
	err = viper.Unmarshal(config)

	if err != nil {
		return err
	}

	return nil
}
