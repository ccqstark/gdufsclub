package controller

import (
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetNotice(c *gin.Context) {

	progressStr := c.Param("/progress")
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

	//公告是否存在
	if !model.IsNoticeExist(clubID.(int), progress) {
		//不存在
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "没有公告",
		})
		return
	}
	//存在
	if notice, ok := model.QueryNotice(clubID.(int), progress); ok == true {
		//设置session:获取到的notice的id、name
		session.Set("notice_id", notice.NoticeID)
		session.Set("club_Name",notice.ClubName)
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"notice": gin.H{
				"club_name": notice.ClubName,
				"progress":  notice.Progress,
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
			"code":        200,
			"msg":         "公告发布成功",
			"template_id": noticeID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "公告发布失败，请重试",
		})
	}
}

//修改公告
func ModifyNotice(c *gin.Context){

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
	clubName := session.Get("clubName")
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