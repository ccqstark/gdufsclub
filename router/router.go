package router

import (
	"github.com/ccqstark/gdufsclub/controller"
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func LoadRouter() *gin.Engine {

	r := gin.Default()

	//session:创建基于cookie的存储引擎,添加密钥，并使用中间件
	store := cookie.NewStore([]byte("wdnmd"))
	r.Use(sessions.Sessions("mysession", store))

	//日志中间件
	r.Use(middleware.LoggerToFile())

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
		}

		//club
		v1Club := v1Group.Group("/club")
		{
			v1Club.POST("/info", controller.SettleNewClub)
			v1Club.POST("/logo",controller.UploadClubLogo)
		}




		//admin
		v1Admin := v1Group.Group("/admin")
		{
			v1Admin.GET("/users", func(c *gin.Context) {
				c.String(200, "/v1/admin/users")
			})
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
