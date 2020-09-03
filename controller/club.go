package controller

import (
	"fmt"
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/ccqstark/gdufsclub/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func SettleNewClub(c *gin.Context) {
	var club model.Club
	if err := c.ShouldBind(&club); err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生某种错误了呢",
		})
		return
	}

	//用户名是否重复
	if model.IsAccountRepeat(club.ClubAccount) {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "此用户名已被注册",
		})
		return
	}

	//插入数据
	if clubID, ok := model.InsertNewClub(&club); ok == true {
		//club_id存入session
		session := sessions.Default(c)
		session.Set("club_id", clubID)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"msg":     "申请提交成功，请耐心等待后台审核",
			"club_id": clubID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "注册失败，请重试",
		})
	}
}

//上传logo
func UploadClubLogo(c *gin.Context) {

	imageConf := util.Cfg.Image
	file, err := c.FormFile("logo")
	if err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	}

	//判断文件类型是否为图片
	fileExt := strings.ToLower(path.Ext(file.Filename))
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "上传失败!只允许png, jpg, gif, jpeg文件",
		})
		return
	}

	//生成随机码文件名
	fileName := util.Md5SaltCrypt(fmt.Sprintf("%s%s", file.Filename, time.Now().String()))
	fileDir := fmt.Sprintf("%s/", imageConf.LogoPath)

	//判断文件夹是否存在
	isExist := util.IsExists(fileDir)
	if !isExist {
		os.Mkdir(fileDir, os.ModePerm)
	}

	//保存至服务器指定目录
	filepath := fmt.Sprintf("%s%s%s", fileDir, fileName, fileExt)
	c.SaveUploadedFile(file, filepath)

	//写入数据库
	session := sessions.Default(c)
	clubID := session.Get("club_id")
	session.Save()
	if clubID == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "未登录或找不到对应社团",
		})
		return
	}

	//插入数据库
	if ok := model.UpdateLogo(clubID.(int), filepath); ok == false {
		//数据库出错
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "上传成功!",
		"result": gin.H{
			"path": filepath,
		},
	})
}
