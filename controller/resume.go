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
	"strconv"
	"strings"
	"time"
)

//用户获取自己填的报名表
func GetResume(c *gin.Context) {

	clubIDStr := c.Param("club_id")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil {
		middleware.Log.Error(err.Error())
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

	if resume, ok := model.QueryResume(userID.(int), clubID); ok == true {
		//设置session:获取到的resume的id
		session.Set("resume_id", resume.ResumeID)
		session.Set("resume_club_id", resume.ClubID)
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"resume": gin.H{
				"name":      resume.Name,
				"sex":       resume.Sex,
				"class":     resume.Class,
				"phone":     resume.Phone,
				"email":     resume.Email,
				"wechat":    resume.Wechat,
				"hobby":     resume.Hobby,
				"advantage": resume.Advantage,
				"self":      resume.Self,
				"reason":    resume.Reason,
				"image":     resume.Image,
				"extra":     resume.Extra,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "还未向此社团提交过报名表",
		})
	}
}

//填写新的报名简历
func FillNewResume(c *gin.Context) {

	var resume model.Resume
	if err := c.ShouldBind(&resume); err != nil {
		middleware.Log.Error(err.Error())
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生某种错误了呢",
		})
		return
	}

	session := sessions.Default(c)
	submitterID := session.Get("user_id")
	session.Save()
	if submitterID == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "暂未登录",
		})
		return
	}

	resume.SubmitterID = submitterID.(int)
	if resumeID, ok := model.InsertNewResume(&resume); ok == true {
		session.Set("resume_id", resumeID)
		session.Save()

		//创建面试进程
		if okk := model.CreateProcess(submitterID.(int), resume.ClubID); okk == false {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "无法创建面试流程",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "报名表提交成功",
			"resume_id": resumeID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "报名表提交失败，请重试",
		})
	}
}

//上传报名表照片
func UploadResumeProfile(c *gin.Context) {

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
	session := sessions.Default(c)
	resumeID := session.Get("resume_id")
	session.Save()
	if resumeID == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "找不到当前报名表",
		})
		return
	}

	if ok := model.UpdateResumeProfile(resumeID.(int), fileNameExt); ok == true {
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

//修改已填简历
func ModifyResume(c *gin.Context) {

	var resume model.Resume
	if err := c.ShouldBind(&resume); err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生某种错误了呢",
		})
		return
	}

	session := sessions.Default(c)
	resumeID := session.Get("resume_id")
	userID := session.Get("user_id")
	clubID := session.Get("resume_club_id")
	session.Save()

	if userID == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "暂未登录",
		})
		return
	}

	if resumeID == nil || clubID == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "找不到此报名表",
		})
		return
	}

	resume.ResumeID = resumeID.(int)
	resume.SubmitterID = userID.(int)
	resume.ClubID = clubID.(int)
	if ok := model.UpdateResumeInfo(&resume); ok == true {
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

//社团获取用户的报名表
func ClubGetResume(c *gin.Context) {

	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
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

	if resume, ok := model.QueryResume(userID, clubID.(int)); ok == true {
		if result, oks := model.QueryInterviewResult(userID, clubID.(int)); oks == true {
			//获取报名表和面试进程状态
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"basic": gin.H{
					"name":      resume.Name,
					"sex":       resume.Sex,
					"class":     resume.Class,
					"phone":     resume.Phone,
					"wechat":    resume.Wechat,
					"image":     resume.Image,
				},
				"other":gin.H{
					"reason":    resume.Reason,
					"self":      resume.Self,
					"hobby":     resume.Hobby,
					"advantage": resume.Advantage,
					"email":     resume.Email,
					"extra":     resume.Extra,
				},
				"result": result,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "查询不到面试状态",
			})
		}

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "还未向此社团提交过报名表",
		})
	}
}
