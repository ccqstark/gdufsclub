package util

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	AppName  string       `json:"app_name"`
	AppMode  string       `json:"app_mode"`
	AppHost  string       `json:"app_host"`
	AppPort  string       `json:"app_port"`
	Database DatabaseConf `json:"database"`
}

type DatabaseConf struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"db_name"`
	Timeout  string `json:"timeout"`
}

var _cfg *Config = nil

//加载配置
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&_cfg); err != nil {
		return nil, err
	}

	return _cfg, nil
}
