package controller

import (
	"fmt"
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/ccqstark/gdufsclub/pkg/sego"
	"github.com/ccqstark/gdufsclub/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
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

//社团入住
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
	fileNameExt := fmt.Sprintf("%s%s", fileName, fileExt)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	}

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
	if ok := model.UpdateLogo(clubID.(int), fileNameExt); ok == false {
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
			"path": fileNameExt,
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

//获取对应面试轮数的用户列表
func GetUserListBrief(c *gin.Context) {

	progressStr := c.Query("progress")
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

	if userList, ok := model.QueryUserListBrief(clubID.(int), progress); ok == true {
		c.IndentedJSON(http.StatusOK, userList)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "获取失败",
		})
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
			"msg":  "找不到此人的报名表",
		})
	}
}

//导出当前面通过者excel
func GetExcel(c *gin.Context) {

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

	//获取此面所有通过者的user_id
	passerID := model.QueryPasser(clubID.(int), progress)

	//获取这些通过者的具体信息
	var infoArr []model.Resume
	var ok bool
	infoArr, ok = model.GainInfoByArray(clubID.(int), passerID)
	if ok == false {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "获取不到数据",
		})
		return
	}

	//生成Excel
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("info_list")
	if err != nil {
		fmt.Printf(err.Error())
	}

	//添加项目标题
	row := sheet.AddRow()

	nameCell := row.AddCell()
	nameCell.Value = "姓名"

	sexCell := row.AddCell()
	sexCell.Value = "性别"

	classCell := row.AddCell()
	classCell.Value = "班级"

	phoneCell := row.AddCell()
	phoneCell.Value = "手机"

	wechatCell := row.AddCell()
	wechatCell.Value = "微信号"

	//载入data
	for _, info := range infoArr {
		row := sheet.AddRow()

		nameCell := row.AddCell()
		nameCell.Value = info.Name

		sexCell := row.AddCell()
		sexCell.Value = info.Sex

		classCell := row.AddCell()
		classCell.Value = info.Class

		phoneCell := row.AddCell()
		phoneCell.Value = info.Phone

		wechatCell := row.AddCell()
		wechatCell.Value = info.Wechat
	}

	fileName := util.Md5SaltCrypt(fmt.Sprintf("%s%s", "excelqlg", time.Now().String()))
	fileDir := "./file/"

	//判断文件夹是否存在
	isExist := util.IsExists(fileDir)
	if !isExist {
		os.Mkdir(fileDir, os.ModePerm)
	}

	//保存至服务器指定的目录
	filepath := fmt.Sprintf("%s%s%s", fileDir, fileName, ".xlsx")
	fileNameExt := fmt.Sprintf("%s%s", fileName, ".xlsx")
	err = file.Save(filepath)

	if err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "导出失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "导出成功",
		"filename": fileNameExt,
	})
}

//社团登录
func ClubLogin(c *gin.Context) {

	var clubAccount model.ClubAccount
	if err := c.ShouldBind(&clubAccount); err != nil {
		middleware.Log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生某种错误了呢",
		})
		return
	}

	club, ok := model.JudgePassword(clubAccount.Account, clubAccount.Password)
	if ok == 1 {
		//记录登录状态
		session := sessions.Default(c)
		session.Set("club_id", club.ClubID)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"code":           200,
			"club_name":      club.ClubName,
			"logo":           club.Logo,
			"total_progress": club.TotalProgress,
			"msg":            "登录成功",
		})
	} else if ok == 2 {
		c.JSON(http.StatusOK, gin.H{
			"code": 402,
			"msg":  "审核未通过",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
	}
}

//社团获取对应面公告
func ClubGetNotice(c *gin.Context) {

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

	if notice, ok := model.ClubQueryNotice(clubID.(int), progress); ok == true {
		c.IndentedJSON(http.StatusOK, notice)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "获取失败",
		})
	}
}

//社团修改信息
func ModifyClubInfo(c *gin.Context) {

	var clubModInfo model.ClubModInfo
	if err := c.ShouldBind(&clubModInfo); err != nil {
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
			"msg":  "未登录或找不到对应社团",
		})
		return
	}

	if ok := model.UpdateClubInfo(clubModInfo, clubID.(int)); ok == true {
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

//获取单个社团信息
func GetOneClubInfo(c *gin.Context) {

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

	if club, ok := model.QueryClubInfo(clubID.(int)); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"club_name":      club.ClubName,
				"club_phone":     club.ClubPhone,
				"club_email":     club.ClubEmail,
				"club_wechat":    club.ClubWechat,
				"total_progress": club.TotalProgress,
				"logo":           club.Logo,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "获取社团信息失败",
		})
	}
}

//获取社团信息
func GetAllClubInfo(c *gin.Context) {

	//不用登录
	if club, ok := model.QueryAllClubInfo(); ok == true {

		c.IndentedJSON(http.StatusOK, club)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "查询失败",
		})
		return
	}
}
