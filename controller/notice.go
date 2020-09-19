package controller

import (
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//用户获得社团的公告,根据自己通过与否查看对应的公告
func GetUserNotice(c *gin.Context) {

	clubIDStr := c.Query("club_id")
	clubID, err1 := strconv.Atoi(clubIDStr)
	if err1 != nil {
		middleware.Log.Error(err1.Error())
	}

	progressStr := c.Query("progress")
	progress, err2 := strconv.Atoi(progressStr)
	if err2 != nil {
		middleware.Log.Error(err2.Error())
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

	//查询面试结果
	var pass int
	var ok bool
	if pass, ok = model.QueryInterviewResult(userID.(int), clubID); ok == false {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "查询不到面试结果",
		})
		return
	}

	//公告是否存在
	if !model.IsNoticeExist(clubID, progress, pass) {
		//不存在
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "公告暂未发布",
		})
		return
	}


	//存在
	if notice, ok := model.QueryNotice(clubID, progress, pass); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"notice": gin.H{
				"club_name": notice.ClubName,
				"content":   notice.Content,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "公告暂未发布",
		})
	}
}

//用户获取成功的公告
func GetSuccessNotice(c *gin.Context) {

	clubIDStr := c.Query("club_id")
	clubID, err1 := strconv.Atoi(clubIDStr)
	if err1 != nil {
		middleware.Log.Error(err1.Error())
	}

	progressStr := c.Query("progress")
	progress, err2 := strconv.Atoi(progressStr)
	if err2 != nil {
		middleware.Log.Error(err2.Error())
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

	//公告是否存在
	if !model.IsNoticeExist(clubID, progress, 1) {
		//不存在
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "公告暂未发布",
		})
		return
	}

	//存在
	if notice, ok := model.QueryNotice(clubID, progress, 1); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"notice": gin.H{
				"club_name": notice.ClubName,
				"content":   notice.Content,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "公告暂未发布",
		})
	}
}

//发布新公告
func PostNewNotice(c *gin.Context) {

	var twoNotice model.TwoNotice
	if err := c.ShouldBind(&twoNotice); err != nil {
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

	var clubName string
	var ok bool
	//查找社团名
	if clubName, ok = model.QueryClubName(clubID.(int)); ok == false {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "社团信息存在异常",
		})
		return
	}

	if ok := model.InsertNewNotice(&twoNotice, clubID.(int), clubName); ok == true {

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "公告发布成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "公告发布失败，请重试",
		})
	}
}

//修改公告
func ModifyNotice(c *gin.Context) {

	var twoNotice model.TwoNotice
	if err := c.ShouldBind(&twoNotice); err != nil {
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

	if ok := model.UpdateNotice(&twoNotice, clubID.(int)); ok == true {
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

//公告统一发布
func PublishNotice(c *gin.Context) {

	progressStr := c.Param("progress")
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

	if ok := model.MakeNoticePublished(clubID.(int), progress); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "操作成功",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "操作失败",
		})
		return
	}

}
