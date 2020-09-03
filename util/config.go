package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	AppName  string       `json:"app_name"`
	AppMode  string       `json:"app_mode"`
	AppHost  string       `json:"app_host"`
	AppPort  string       `json:"app_port"`
	Database DatabaseConf `json:"database"`
	Redis    RedisConf    `json:"redis"`
	Logger   LoggerConf   `json:"logger"`
	Image    ImageConf    `json:"image"`
}

type DatabaseConf struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"db_name"`
	Timeout  string `json:"timeout"`
}

type RedisConf struct {
	IdleConnection int    `json:"idle_connection"`
	Protocol       string `json:"protocol"`
	HostPort       string `json:"host_port"`
	Password       string `json:"password"`
	Key            string `json:"key"`
}

type LoggerConf struct {
	LogFilePath string `json:"log_file_path"`
	LogFileName string `json:"log_file_name"`
}

type ImageConf struct {
	LogoPath    string `json:"logo_path"`
	ProfilePath string `json:"profile_path"`
}

//全局
var Cfg *Config = nil

//加载全局配置，先于main函数执行
func init() {
	file, err := os.Open("./config/conf.json")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&Cfg); err != nil {
		fmt.Println(err)
	}
}
