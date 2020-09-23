package controller

import (
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

//用户查看自己的面试进程
func GetProcess(c *gin.Context) {

	//用openid获取用户id
	openid := c.Query("openid")
	var userID int
	var ok bool
	if userID, ok = model.GetUserIDByOpenid(openid); ok == false {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "查找不到用户",
		})
		return
	}

	if progress, ok := model.QueryProcess(userID); ok == true {
		c.IndentedJSON(http.StatusOK, progress)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "当前无面试",
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

//批量通过面试者
func PassBatch(c *gin.Context) {

	var batchUser model.BatchUser
	if err := c.ShouldBind(&batchUser); err != nil {
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

	if ok := model.PassBatchInterviewee(batchUser.Interviewee, clubID.(int), batchUser.Progress); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "操作成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请选择要通过的人选",
		})
	}
}
