package controller

import (
	"fmt"
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/ccqstark/gdufsclub/pkg/sego"
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

var Segmenter sego.Segmenter

func init() {
	// 载入词典
	Segmenter.LoadDictionary("./pkg/sego/data/dictionary.txt")
}

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

//搜索社团
func SearchClub(c *gin.Context) {

	keyWord := c.Query("key_word")

	// 分词
	text := []byte(keyWord)
	segments := Segmenter.Segment(text)

	// 处理分词结果
	cutWordStr := sego.SegmentsToString(segments, true)
	cutWordSlice := strings.Split(cutWordStr, ",")
	cutWordSlice = cutWordSlice[:len(cutWordSlice)-1]

	resultGather := model.SearchByWord(cutWordSlice)

	//去重
	var resultGatherClean []model.Club //去重后结果
	var flag bool
	for _, dirtyResult := range resultGather {
		flag = false
		for _, cleanResult := range resultGatherClean {
			if cleanResult.ClubID == dirtyResult.ClubID {
				flag = true
				break
			}
		}
		if !flag {
			resultGatherClean = append(resultGatherClean, dirtyResult)
		}
	}

	if resultGatherClean != nil {
		c.IndentedJSON(http.StatusOK, resultGatherClean)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "没有找到相关结果",
		})
	}
}


//获取一共的页数
func GetUserTotalPage(c *gin.Context){

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

	if page,ok := model.QueryUserTotalPage(clubID.(int),progress);ok==true{
		c.JSON(http.StatusOK,gin.H{
			"code":200,
			"total":page,
		})
	}else {
		c.JSON(http.StatusOK,gin.H{
			"code":400,
			"msg":"查询失败",
		})
	}
}




//获取对应面试轮数的用户列表
func GetUserListBrief(c *gin.Context){

	progressStr := c.Param("progress")
	progress, err := strconv.Atoi(progressStr)
	if err != nil {
		middleware.Log.Error(err.Error())
	}








}


//获得用户的信息
func GetUserResume(c *gin.Context) {

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


	if resume, ok := model.QueryResume(userID, clubID); ok == true {
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
