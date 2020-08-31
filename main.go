package main

import (
	"github.com/ccqstark/gdufsclub/router"
	"github.com/ccqstark/gdufsclub/util"
)

func main() {
	//加载配置
	cfg, err :=util.LoadConfig("./config/conf.json")
	if err != nil {
		panic(err.Error())
	}

	//加载并启动路由
	r := router.LoadRouter()


	r.Run(":"+cfg.AppPort)  //:::8060

}


