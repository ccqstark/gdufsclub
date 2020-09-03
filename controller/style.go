package controller

import (
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//社团获得自己的表样式
func GetStyle(c *gin.Context) {

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

	//报名表是否存在
	if !model.IsStyleExist(clubID.(int)) {
		//不存在
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "还未创建过报名表样式噢",
		})
		return
	}
	//存在
	if style, ok := model.QueryStyle(clubID.(int)); ok == true {
		//设置session:获取到的style的id
		session.Set("style_id", style.StyleID)
		session.Set("club_name", style.ClubName)
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"template": gin.H{
				"club_name":       style.ClubName,
				"style_name":      style.StyleName,
				"style_sex":       style.StyleSex,
				"style_class":     style.StyleClass,
				"style_phone":     style.StylePhone,
				"style_email":     style.StyleEmail,
				"style_wechat":    style.StyleWechat,
				"style_image":     style.StyleImage,
				"style_hobby":     style.StyleHobby,
				"style_advantage": style.StyleAdvantage,
				"style_self":      style.StyleSelf,
				"style_reason":    style.StyleReason,
				"style_extra":     style.StyleExtra,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "样式获取失败",
		})
	}
}

//用户获得表样式，无登录限制
func GetUserStyle(c *gin.Context) {

	clubIDStr := c.Param("club_id")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil {
		middleware.Log.Error(err.Error())
	}

	//报名表是否存在
	if !model.IsStyleExist(clubID) {
		//不存在
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "改社团还未开放报名表噢",
		})
		return
	}
	//存在
	if style, ok := model.QueryStyle(clubID); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"template": gin.H{
				"club_name":       style.ClubName,
				"style_name":      style.StyleName,
				"style_sex":       style.StyleSex,
				"style_class":     style.StyleClass,
				"style_phone":     style.StylePhone,
				"style_email":     style.StyleEmail,
				"style_wechat":    style.StyleWechat,
				"style_image":     style.StyleImage,
				"style_hobby":     style.StyleHobby,
				"style_advantage": style.StyleAdvantage,
				"style_self":      style.StyleSelf,
				"style_reason":    style.StyleReason,
				"style_extra":     style.StyleExtra,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "报名表获取失败",
		})
	}
}

//创建新的表样式
func MakeNewStyle(c *gin.Context) {

	var style model.Style
	if err := c.ShouldBind(&style); err != nil {
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
		style.ClubName = clubName
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "社团信息存在异常",
		})
		return
	}

	style.ClubID = clubID.(int)
	if styleID, ok := model.InsertNewStyle(&style); ok == true {
		session.Set("style_id", styleID)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"msg":      "报名表样式创建成功",
			"style_id": styleID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "样式已存在或创建失败，请重试",
		})
	}
}

//修改表样式
func ModifyStyle(c *gin.Context) {

	var style model.Style
	if err := c.ShouldBind(&style); err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生某种错误了呢",
		})
		return
	}

	session := sessions.Default(c)
	styleID := session.Get("style_id")
	clubID := session.Get("club_id")
	clubName := session.Get("club_name")
	session.Save()

	if clubID == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "暂未登录",
		})
		return
	}

	if styleID == nil || clubName == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "找不到此报名表",
		})
		return
	}

	style.StyleID = styleID.(int)
	style.ClubID = clubID.(int)
	style.ClubName = clubName.(string)
	if ok := model.UpdateStyle(&style); ok == true {
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
