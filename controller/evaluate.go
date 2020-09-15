package controller

import (
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//获取评论
func GetAEvaluate(c *gin.Context) {

	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		middleware.Log.Error(err.Error())
	}

	progressStr := c.Query("progress")
	progress, err := strconv.Atoi(progressStr)
	if err != nil {
		middleware.Log.Error(err.Error())
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

	//查询评论
	if evaluate, ok := model.QueryEvaluate(clubID.(int), userID, progress); ok == true {
		session.Set("evaluate_id", evaluate.EvaluateID)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"content": evaluate.Content,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    401,
			"content": "",
		})
	}
}

//创建评论
func NewAEvaluate(c *gin.Context) {

	var evaluate model.Evaluate
	if err := c.ShouldBind(&evaluate); err != nil {
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

	evaluate.ClubID = clubID.(int)
	if id, ok := model.CreateEvaluate(evaluate); ok == true {
		session.Set("evaluate_id", id)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "评价成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "评价失败",
		})
	}
}

//修改评论
func ModifyEvaluate(c *gin.Context) {

	var evaluate model.Evaluate
	if err := c.ShouldBind(&evaluate); err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生某种错误了呢",
		})
		return
	}

	session := sessions.Default(c)
	evaluateID := session.Get("evaluate_id")
	session.Save()
	if evaluateID == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "获取评论失败",
		})
		return
	}

	evaluate.EvaluateID = evaluateID.(int)

	if ok := model.UpdateEvaluate(&evaluate); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "修改成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "修改失败",
		})
	}
}
