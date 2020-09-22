package controller

import (
	"encoding/json"
	"fmt"
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	code2sessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	appID           = "wx573a58a5e1d18401"
	appSecret       = "558b4f2eb9c8c060a507cbdde6924f70"
)

func Demo(c *gin.Context) {

	clubIDStr := c.Param("club_id")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil {
		middleware.Log.Error(err.Error())
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		middleware.Log.Error(err.Error())
	}

	session := sessions.Default(c)
	session.Set("user_id", userID)
	session.Set("club_id", clubID)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"club_id": clubID,
		"code":    200,
	})
}

//用code获取openid登录
func UserLogin(c *gin.Context) {

	var login model.Login
	if err := c.ShouldBind(&login); err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生某种错误了呢",
		})
		return
	}
	//获取openid
	openid := getOpenID(login.Code)
	//获取openid失败
	if openid == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "授权失败",
		})
		return
	}

	//middleware.Log.Info(openid)
	//清一下session
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	if userID, ok := model.AuthUser(openid); ok == true {
		//设置登录状态session
		session.Set("user_id", userID)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "登录成功",
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "授权失败",
		})
	}
}

//获取openid
func getOpenID(code string) string {

	openid, err := sendWxAuthAPI(code)
	if err != nil {
		return ""
	}

	return openid
}

//请求微信官方接口
func sendWxAuthAPI(code string) (string, error) {
	//拼接请求url,并请求获得response
	url := fmt.Sprintf(code2sessionURL, appID, appSecret, code)
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", err
	}
	//解析返回的json
	var wxMap map[string]string
	err = json.NewDecoder(response.Body).Decode(&wxMap)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	return wxMap["openid"], nil
}
