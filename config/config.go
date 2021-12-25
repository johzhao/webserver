package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Logger Logger `yaml:"logging"`
	DB     DB     `yaml:"db"`
	Server Server `yaml:"server"`
}

type Logger struct {
	Level     string `yaml:"level"`
	Filepath  string `yaml:"filepath"`
	MaxSize   int    `yaml:"max_size"`
	MaxBackup int    `yaml:"max_backup"`
	MaxAge    int    `yaml:"max_age"`
	Compress  bool   `yaml:"compress"`
}

type DB struct {
	DriverName     string `yaml:"driver_name"`
	DataSourceName string `yaml:"data_source_name"`
}

type Server struct {
	Address string `yaml:"address"`
}

func LoadConfig(filepath string) (Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	var result Config
	if err := yaml.Unmarshal(data, &result); err != nil {
		return Config{}, err
	}

	return result, nil
}
