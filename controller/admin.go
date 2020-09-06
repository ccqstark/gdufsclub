package controller

import (
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//获得还未审核的
func GetAllNotPass(c *gin.Context) {

	if notClub, ok := model.QueryAllNotPass(); ok == true {
		c.IndentedJSON(200, notClub)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "查询失败",
		})
	}
}

func GetAllPass(c *gin.Context) {

	if enterClub, ok := model.QueryAllPass(); ok == true {
		c.IndentedJSON(200, enterClub)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "查询失败",
		})
	}
}

//审核通过一个社团
func AuditOne(c *gin.Context) {

	clubIDStr := c.Param("club_id")
	clubID, err1 := strconv.Atoi(clubIDStr)
	if err1 != nil {
		middleware.Log.Error(err1.Error())
	}

	statusStr := c.Param("status")
	status, err2 := strconv.Atoi(statusStr)
	if err2 != nil {
		middleware.Log.Error(err2.Error())
	}

	if model.AuditOneClub(clubID, status) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "操作成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "操作失败",
		})
	}
}
