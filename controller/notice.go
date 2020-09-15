package controller

import (
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//获得社团自己的公告
func GetNotice(c *gin.Context) {

	progressStr := c.Query("progress")
	progress, err := strconv.Atoi(progressStr)
	if err != nil {
		middleware.Log.Error(err.Error())
	}

	passStr := c.Query("pass")
	pass, err := strconv.Atoi(passStr)
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

	//公告是否存在
	if !model.IsNoticeExist(clubID.(int), progress, pass) {
		//不存在
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "没有公告",
		})
		return
	}
	//存在
	if notice, ok := model.QueryNotice(clubID.(int), progress, pass); ok == true {
		//设置session:获取到的notice的id、name
		session.Set("notice_id", notice.NoticeID)
		session.Set("club_name", notice.ClubName)
		session.Save()
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
			"msg":  "公告获取失败",
		})
	}
}

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

//发布新公告
func PostNewNotice(c *gin.Context) {

	var notice model.Notice
	if err := c.ShouldBind(&notice); err != nil {
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

	//查找社团名
	if clubName, ok := model.QueryClubName(clubID.(int)); ok == true {
		notice.ClubName = clubName
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "社团信息存在异常",
		})
		return
	}

	notice.ClubID = clubID.(int)
	if noticeID, ok := model.InsertNewNotice(&notice); ok == true {
		session.Set("notice_id", noticeID)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "公告发布成功",
			"notice_id": noticeID,
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

	var notice model.Notice
	if err := c.ShouldBind(&notice); err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生某种错误了呢",
		})
		return
	}

	session := sessions.Default(c)
	noticeID := session.Get("notice_id")
	clubID := session.Get("club_id")
	clubName := session.Get("club_name")
	session.Save()

	if noticeID == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "公告未创建",
		})
		return
	}

	if clubID == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "暂未登录",
		})
		return
	}

	if clubName == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "社团信息出现问题了",
		})
		return
	}

	notice.NoticeID = noticeID.(int)
	notice.ClubID = clubID.(int)
	notice.ClubName = clubName.(string)
	if ok := model.UpdateNotice(&notice); ok == true {
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
