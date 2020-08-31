package router

import (
	"github.com/ccqstark/gdufsclub/controller"
	"github.com/gin-gonic/gin"
)

func LoadRouter() *gin.Engine {

	r := gin.Default()

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
		v1User:=v1Group.Group("/user")
		{
			v1User.GET("/first",controller.PickFirstUser)
		}

		//club
		v1Club := v1Group.Group("/club")
		{
			v1Club.POST("",controller.SettleNewClub)
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
