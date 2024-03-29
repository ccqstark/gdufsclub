package middleware

import (
	"fmt"
	"github.com/ccqstark/gdufsclub/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

//全局logger
var Log *logrus.Logger

func init(){

	//加载全局配置
	loggerConf := util.Cfg.Logger
	logFilePath := loggerConf.LogFilePath
	logFileName := loggerConf.LogFileName

	//日志文件
	fileName := path.Join(logFilePath, logFileName)

	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//实例化
	Log = logrus.New()

	//设置输出
	Log.Out = src

	//设置日志级别
	Log.SetLevel(logrus.TraceLevel)

	//设置日志格式
	Log.SetFormatter(&logrus.TextFormatter{
		//设置时间格式
		TimestampFormat: "2006-01-02 15:04:05",
	})
}


// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		Log.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
