package controller

import (
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllNotPass(c *gin.Context){

	if notClub,ok := model.QueryAllNotPass();ok==true{
		c.IndentedJSON(200, notClub)
	}else {
		c.JSON(http.StatusOK,gin.H{
			"code":400,
			"msg":"查询失败",
		})
	}
}
