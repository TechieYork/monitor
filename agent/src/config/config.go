package config

import (
	"github.com/spf13/viper"
)

const Version = "0.0.1"

//Transfer queue information
type TransferQueueInfo struct{
	BufferSize int `mapstructure:"buffer_size"`
}

//Node information
type NodeInfo struct {
	Name string `mapstructure:"name"`
	IP string `mapstructure:"ip"`
	TransferQueue TransferQueueInfo `mapstructure:"transfer_queue"`
}

//Admin information
type AdminInfo struct {
	Ip string `mapstructure:"ip"`
	Port int `mapstructure:"port"`
}

//Input plugin information
type InputPluginInfo struct {
	Name string `mapstructure:"plugin_name"`
	Path string `mapstructure:"plugin_path"`
	Duration int `mapstructure:"duration"`
	Active bool `mapstructure:"active"`
	PluginConfig map[string]string `mapstructure:"config"`
}

//Output plugin information
type OutputPluginInfo struct {
	Name string `mapstructure:"plugin_name"`
	Path string `mapstructure:"plugin_path"`
	Active bool `mapstructure:"active"`
	Inputs map[string]bool `mapstructure:"inputs"`
	PluginConfig map[string]string `mapstructure:"config"`
}

//Config sturcture
type Config struct {
	Node NodeInfo `mapstructure:"node"`
	Admin AdminInfo `mapstructure:"admin"`
	Inputs []InputPluginInfo `mapstructure:"input_plugin"`
	Outputs []OutputPluginInfo `mapstructure:"output_plugin"`
}

//Global config
var globalConfig *Config

//New Config
func NewConfig() *Config {
	return &Config{
		Node:NodeInfo{Name:"unknown", TransferQueue:TransferQueueInfo{BufferSize:1000}},
		Admin:AdminInfo{Ip:"127.0.0.1", Port:9000},
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
