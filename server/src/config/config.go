package config

import "github.com/spf13/viper"

const Version = "0.0.1"

//Mongodb information
type MongodbInfo struct {
	Address string `mapstructure:"addr"`
	DBName string `mapstructure:"db_name"`
}

//Influxdb information
type InfluxdbInfo struct {
	Address string `mapstructure:"addr"`
	DBName string `mapstructure:"db_name"`
}

//Nsq information
type NsqInfo struct {
	Address string `mapstructure:"addr"`
	TopicName string `mapstructure:"topic"`
}

//Server information
type ServerInfo struct {
	Address string `mapstructure:"addr"`
}

//Config sturcture
type Config struct {
	Server ServerInfo `mapstructure:"server"`
	Mongodb MongodbInfo `mapstructure:"mongodb"`
	Nsq NsqInfo `mapstructure:"nsq"`
	Influxdb InfluxdbInfo `mapstructure:"influxdb"`
}

//Global config
var globalConfig *Config

//New Config
func NewConfig() *Config {
	return &Config{
		Server:ServerInfo{Address:""},
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
