package router

import (
	"github.com/ccqstark/gdufsclub/controller"
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func LoadRouter() *gin.Engine {

	r := gin.Default()

	//session:创建基于redis的存储引擎,添加密钥，并使用中间件
	redisConf := util.Cfg.Redis
	store, _ := redis.NewStore(redisConf.IdleConnection, redisConf.Protocol, redisConf.HostPort, redisConf.Password, []byte(redisConf.Key))
	r.Use(sessions.Sessions("mysession", store))

	//日志中间件
	r.Use(middleware.LoggerToFile())

	//使用跨域中间件
	r.Use(middleware.Cors())


	//v1路由组
	v1Group := r.Group("/v1")
	{
		//嵌套分模块
		//index
		v1Index := v1Group.Group("/index")
		{
			v1Index.GET("/ping", controller.PingPong)

		}

		//user
		v1User := v1Group.Group("/user")
		{
			v1User.GET("/first", controller.PickFirstUser)
			v1User.POST("", controller.Demo)
		}

		//club
		v1Club := v1Group.Group("/club")
		{
			v1Club.POST("/info", controller.SettleNewClub)
			v1Club.POST("/logo", controller.UploadClubLogo)
		}

		//template
		v1Template := v1Group.Group("/template")
		{
			v1Template.GET("", controller.GetTemplate)
			v1Template.POST("/info", controller.CreateNewTemplate)
			v1Template.POST("/profile", controller.UploadTplProfile)
			v1Template.PUT("/info", controller.ModifyTemplate)
		}

		//resume
		v1Resume := v1Group.Group("/resume")
		{
			v1Resume.GET("/:club_id", controller.GetResume)
			v1Resume.POST("/info", controller.FillNewResume)
			v1Resume.POST("/profile", controller.UploadResumeProfile)
			v1Resume.PUT("/info", controller.ModifyResume)
		}

		//style
		v1Style := v1Group.Group("/style")
		{
			v1Style.GET("", controller.GetStyle)
			v1Style.GET("/user_style/:club_id", controller.GetUserStyle)
			v1Style.POST("", controller.MakeNewStyle)
			v1Style.PUT("", controller.ModifyStyle)
		}

		//notice
		v1Notice := v1Group.Group("/notice")
		{
			v1Notice.GET("", controller.GetNotice)
			v1Notice.GET("/user_notice",controller.GetUserNotice)
			v1Notice.POST("", controller.PostNewNotice)
			v1Notice.PUT("", controller.ModifyNotice)
		}

		//process
		v1Process := v1Group.Group("/process")
		{
			v1Process.GET("/:club_id",controller.GetProcess)
			v1Process.PUT("/result",controller.OperateOne)

		}

		//admin
		v1Admin := v1Group.Group("/admin")
		{
			v1Admin.GET("/not",controller.GetAllNotPass)
			//v1Admin.PUT(":/club_id",controller.)
		}


	}

	//v2路由组
	v2Group := r.Group("/v2")
	{
		v2Group.GET("/", func(c *gin.Context) {
			c.String(200, "v2Group")
		})
	}

	return r
}
