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
)

//获得还未审核的
func GetAllNotPass(c *gin.Context) {

	if notClub, ok := model.QueryAllNotPass(); ok == true {
		c.IndentedJSON(200, notClub)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "查询失败",
		})
	}
}

func GetAllPass(c *gin.Context) {

	if enterClub, ok := model.QueryAllPass(); ok == true {
		c.IndentedJSON(http.StatusOK, enterClub)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "查询失败",
		})
	}
}

//审核通过一个社团
func AuditOne(c *gin.Context) {

	clubIDStr := c.Param("club_id")
	clubID, err1 := strconv.Atoi(clubIDStr)
	if err1 != nil {
		middleware.Log.Error(err1.Error())
	}

	statusStr := c.Param("status")
	status, err2 := strconv.Atoi(statusStr)
	if err2 != nil {
		middleware.Log.Error(err2.Error())
	}

	if model.AuditOneClub(clubID, status) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "操作成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "操作失败",
		})
	}
}

//上传轮播图
func UploadAD(c *gin.Context) {

	imageConf := util.Cfg.Image
	file, err := c.FormFile("ad")
	if err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	}

	adIDStr := c.Param("ad_id")

	//判断文件类型是否为图片
	fileExt := strings.ToLower(path.Ext(file.Filename))
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "上传失败!只允许png, jpg, gif, jpeg文件",
		})
		return
	}

	//文件名
	fileName := adIDStr
	fileDir := fmt.Sprintf("%s/", imageConf.ADPath)

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

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "上传成功!",
		"result": gin.H{
			"path": fileNameExt,
		},
	})
}

func GetAllCustomField(c *gin.Context) {

	if field, ok := model.QueryAllCustomField(); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code":  200,
			"field": field,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "获取失败",
		})
	}

}
