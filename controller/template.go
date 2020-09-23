package controller

import (
	"fmt"
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/ccqstark/gdufsclub/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

//模板二合一
func TplTwoInOne(c *gin.Context) {

	//二合一表信息提取
	var template model.Template
	template.TemplateName = c.PostForm("template_name")
	template.TemplateSex = c.PostForm("template_sex")
	template.TemplateClass = c.PostForm("template_class")
	template.TemplatePhone = c.PostForm("template_phone")
	template.TemplateEmail = c.PostForm("template_email")
	template.TemplateWechat = c.PostForm("template_wechat")
	template.TemplateHobby = c.PostForm("template_hobby")
	template.TemplateAdvantage = c.PostForm("template_advantage")
	template.TemplateSelf = c.PostForm("template_self")

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

	template.UserID = userID

	//profile图片
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
	fileDir := fmt.Sprintf("%s/", imageConf.ProfilePath)

	//判断文件夹是否存在
	isExist := util.IsExists(fileDir)
	if !isExist {
		os.Mkdir(fileDir, os.ModePerm)
	}

	//保存至服务器指定目录
	filepath := fmt.Sprintf("%s%s%s", fileDir, fileName, fileExt)
	fileNameExt := fmt.Sprintf("%s%s", fileName, fileExt)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	}

	//写入数据库
	template.TemplateImage = fileNameExt
	if _, ok := model.InsertNewTemplate(&template); ok == true {

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "提交成功!",
			"result": gin.H{
				"path": fileNameExt,
			},
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "模板提交失败，请重试",
		})
		return
	}
}

//获得模板
func GetTemplate(c *gin.Context) {

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

	//模板是否存在
	if !model.IsTemplateExist(userID) {
		//不存在
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "同学你还没有创建过模板噢",
		})
		return
	}
	//存在
	if tpl, ok := model.QueryTemplate(userID); ok == true {

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"template": gin.H{
				"template_id": tpl.TemplateID,
				"name":        tpl.TemplateName,
				"class":       tpl.TemplateClass,
				"sex":         tpl.TemplateSex,
				"wechat":      tpl.TemplateWechat,
				"email":       tpl.TemplateEmail,
				"phone":       tpl.TemplatePhone,
				"hobby":       tpl.TemplateHobby,
				"self":        tpl.TemplateSelf,
				"advantage":   tpl.TemplateAdvantage,
				"image":       tpl.TemplateImage,
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

	tpl.UserID = userID
	if tplID, ok := model.InsertNewTemplate(&tpl); ok == true {

		c.JSON(http.StatusOK, gin.H{
			"code":        200,
			"msg":         "模板创建成功",
			"template_id": tplID,
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

	tplIDStr := c.Query("tpl_id")
	tplID, err := strconv.Atoi(tplIDStr)
	if err != nil {
		middleware.Log.Error(err.Error())
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
	fileDir := fmt.Sprintf("%s/", imageConf.ProfilePath)

	//判断文件夹是否存在
	isExist := util.IsExists(fileDir)
	if !isExist {
		os.Mkdir(fileDir, os.ModePerm)
	}

	//保存至服务器指定目录
	filepath := fmt.Sprintf("%s%s%s", fileDir, fileName, fileExt)
	fileNameExt := fmt.Sprintf("%s%s", fileName, fileExt)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	}

	//写入数据库
	if ok := model.UpdateTplProfile(tplID, fileNameExt); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "上传成功!",
			"result": gin.H{
				"path": fileNameExt,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "上传失败",
		})
	}
}

//修改模板
func ModifyTemplate(c *gin.Context) {

	var tpl model.Template
	if err := c.ShouldBind(&tpl); err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生某种错误了呢",
		})
		return
	}

	tplIDStr := c.Query("tpl_id")
	tplID, err := strconv.Atoi(tplIDStr)
	if err != nil {
		middleware.Log.Error(err.Error())
	}

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

	tpl.TemplateID = tplID
	tpl.UserID = userID
	if ok := model.UpdateTemplateInfo(&tpl); ok == true {
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
