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
			v1User.POST("", controller.Demo)
			v1User.POST("/openid", controller.UserLogin)
		}

		//club
		v1Club := v1Group.Group("/club")
		{
			//社团入住
			v1Club.POST("/info", controller.SettleNewClub)
			v1Club.POST("/logo", controller.UploadClubLogo)

			//社团登录
			v1Club.POST("login", controller.ClubLogin)

			//获取所有社团
			v1Club.GET("/all", controller.GetAllClubInfo)

			//搜索
			v1Club.GET("/search", controller.SearchClub)

			//获取面试者分页列表
			v1Club.GET("/interviewee_list", controller.GetUserListBrief)

			//社团获取面试者信息
			v1Club.GET("/user_resume/:club_id/:user_id", controller.GetUserResume)

			//社团获取公告
			v1Club.GET("/notice/:progress", controller.ClubGetNotice)

			//导出excel
			v1Club.GET("/excel/:progress", controller.GetExcel)

			//获得社团信息
			v1Club.GET("/info", controller.GetOneClubInfo)
			//修改社团信息
			v1Club.PUT("/info", controller.ModifyClubInfo)
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
			v1Resume.GET("/for_user/:club_id", controller.GetResume)
			v1Resume.POST("/info", controller.FillNewResume)
			v1Resume.POST("/profile", controller.UploadResumeProfile)
			v1Resume.PUT("/info", controller.ModifyResume)
			//社团获取用户简历和面试状态
			v1Resume.GET("/for_club/:user_id", controller.ClubGetResume)
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
			v1Notice.GET("", controller.GetUserNotice)
			v1Notice.POST("", controller.PostNewNotice)
			v1Notice.PUT("", controller.ModifyNotice)
			v1Notice.PUT("/publish/:progress", controller.PublishNotice)
		}

		//process
		v1Process := v1Group.Group("/process")
		{
			v1Process.GET("", controller.GetProcess)
			v1Process.PUT("/result", controller.OperateOne)
			v1Process.PUT("/batch", controller.PassBatch)
		}

		//evaluate
		v1Evaluate := v1Group.Group("/evaluate")
		{
			v1Evaluate.GET("", controller.GetAEvaluate)
			v1Evaluate.POST("", controller.NewAEvaluate)
			v1Evaluate.PUT("", controller.ModifyEvaluate)

		}

		//admin
		v1Admin := v1Group.Group("/admin")
		{
			v1Admin.GET("/not", controller.GetAllNotPass)
			v1Admin.GET("/enter", controller.GetAllPass)
			v1Admin.PUT("/:club_id/:status", controller.AuditOne)
			v1Admin.POST("/ad/:ad_id", controller.UploadAD)
			v1Admin.GET("/field", controller.GetAllCustomField)
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
