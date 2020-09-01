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

//获得模板
func GetTemplate(c *gin.Context) {
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

	//模板是否存在
	if !model.IsTemplateExist(userID.(int)) {
		//不存在
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "同学你还没有创建过模板噢",
		})
		return
	}
	//存在
	if tpl, ok := model.QueryTemplate(userID.(int)); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"template": gin.H{
				"name":      tpl.TemplateName,
				"class":     tpl.TemplateClass,
				"sex":       tpl.TemplateSex,
				"wechat":    tpl.TemplateWechat,
				"email":     tpl.TemplateEmail,
				"phone":     tpl.TemplatePhone,
				"hobby":     tpl.TemplateHobby,
				"self":      tpl.TemplateSelf,
				"advantage": tpl.TemplateAdvantage,
				"image":     tpl.TemplateImage,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "模板获取失败",
		})
	}
}

//创建新模板
func CreateNewTemplate(c *gin.Context) {

	var tpl model.Template
	if err := c.ShouldBind(&tpl); err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生某种错误了呢",
		})
		return
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

	tpl.UserID = userID.(int)
	if TplID, ok := model.InsertNewTemplate(&tpl); ok == true {
		session.Set("template_id", TplID)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"code":        200,
			"msg":         "模板创建成功",
			"template_id": TplID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "模板已存在或创建失败，请重试",
		})
	}
}

//上传模板头像
func UploadTplProfile(c *gin.Context) {

	imageConf := util.Cfg.Image
	file, err := c.FormFile("profile")
	if err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "图片获取失败!",
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

	//生成不重复文件名
	fileName := util.Md5SaltCrypt(fmt.Sprintf("%s%s", file.Filename, time.Now().String()))
	fileDir := fmt.Sprintf("%s/%d%s/", imageConf.ProfilePath, time.Now().Year(), time.Now().Month().String())

	//判断文件夹是否存在
	isExist := util.IsExists(fileDir)
	if !isExist {
		os.Mkdir(fileDir, os.ModePerm)
	}

	filepath := fmt.Sprintf("%s%s%s", fileDir, fileName, fileExt)
	c.SaveUploadedFile(file, filepath)
	//写入数据库
	session := sessions.Default(c)
	TplId := session.Get("template_id")
	session.Save()
	if TplId == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "找不到当前模板",
		})
		return
	}

	if ok := model.UpdateTplProfile(TplId.(int), filepath); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "上传成功!",
			"result": gin.H{
				"path": filepath,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "上传失败",
		})
	}
}
