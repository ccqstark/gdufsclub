package controller

import (
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//用户查看自己的面试进程
func GetProcess(c *gin.Context) {

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

	if progress, ok := model.QueryProcess(userID); ok == true {
		c.IndentedJSON(http.StatusOK, progress)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "当前无面试",
		})
	}
}

//通过或者不通过一个人的面试
func OperateOne(c *gin.Context) {

	var processUser model.ProcessUser
	if err := c.ShouldBind(&processUser); err != nil {
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

	if ok := model.OperateOnePerson(clubID.(int), processUser.UserID, processUser.Pass, processUser.Department); ok == true {
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

//批量通过面试者
func PassBatch(c *gin.Context) {

	var batchUser model.BatchUser
	if err := c.ShouldBind(&batchUser); err != nil {
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

	if ok := model.PassBatchInterviewee(batchUser.Interviewee, clubID.(int), batchUser.Progress, batchUser.Department); ok == true {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "操作成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请选择要通过的人选",
		})
	}
}

//用户查看自己被录用的面试进程
func GetOfferProcess(c *gin.Context) {

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

	if offerList, ok := model.QueryOfferProcess(userID); ok == true {
		c.IndentedJSON(http.StatusOK, offerList)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "获取失败",
		})
	}
}

//接收或不接受Offer
func ReceiveOffer(c *gin.Context) {

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

	processIDStr := c.Query("process_id")
	processID, err := strconv.Atoi(processIDStr)
	if err != nil {
		middleware.Log.Error(err.Error())
	}

	passStr := c.Query("pass")
	pass, err := strconv.Atoi(passStr)
	if err != nil {
		middleware.Log.Error(err.Error())
	}

	//时间判断
	stopTimeStr := "2020-11-18 23:59:59" //截止时间
	layout := "2006-01-02 15:04:05"
	stopTime, _ := time.ParseInLocation(layout, stopTimeStr, time.Local)
	currentTime := time.Now() //当前时间
	//过了截止时间
	if currentTime.Unix() > stopTime.Unix() {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "时间已经截止",
		})
		return
	}

	if pass == 1 {
		//大于2个判断
		offerNum := 0
		var offerClubID []int
		if offerList, ok := model.QueryOfferProcess(userID); ok == true {
			for _, v := range offerList {
				if v.Result == 1 { //拿到offer
					flag := 0
					for _, cid := range offerClubID { //可以加同社团多个部门，但是不能2个以上社团
						if v.ClubID == cid {
							flag = 1
							break
						}
					}
					if flag == 0 {
						offerNum++ //offer数量
						offerClubID = append(offerClubID, v.ClubID)
					}
				}
			}
		}

		//现在要接收的是否为之前接收过offer的社团
		oldClubOffer := 0
		nowClubID := model.QueryClubIDByProcessID(processID)
		for _, v := range offerClubID {
			if v == nowClubID {
				oldClubOffer = 1
				break
			}
		}

		//已经2个了
		if offerNum == 2 && oldClubOffer == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "最多只能选择2个社团",
			})
			return
		} else {
			//接收
			if ok := model.ReceiveOffer(processID, pass); ok == true {
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
	} else {
		//不接收
		if ok := model.ReceiveOffer(processID, pass); ok == true {
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
}
