package main

import (
	"github.com/ccqstark/gdufsclub/router"
	"github.com/ccqstark/gdufsclub/util"
)

func main() {

	//加载并启动路由
	r := router.LoadRouter()

	r.Run(":" + util.Cfg.AppPort) //:::8060
}

//49.234.82.226
