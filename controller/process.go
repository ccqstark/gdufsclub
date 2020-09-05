package controller

import (
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//用户查看自己的面试进程
func GetProcess(c *gin.Context) {

	clubIDStr := c.Param("club_id")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil {
		middleware.Log.Error(err.Error())
	}

	session := sessions.Default(c)
	userID := session.Get("user_id")
	session.Save()
	if userID == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "暂未登录",
		})
		return
	}

	if progress, result, ok := model.QueryProcess(userID.(int), clubID); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"progress": progress,
			"result":   result,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "获取面试进程失败",
		})
	}
}

//通过或者不通过一个人的面试
func OperateOne(c *gin.Context) {

	var processUser model.ProcessUser
	if err := c.ShouldBind(&processUser); err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生某种错误了呢",
		})
		return
	}

	session := sessions.Default(c)
	clubID := session.Get("club_id")
	session.Save()
	if clubID == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "暂未登录",
		})
		return
	}

	if ok := model.OperateOnePerson(clubID.(int), processUser.UserID, processUser.Pass); ok == true {
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
