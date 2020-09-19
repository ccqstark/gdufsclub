package model

import (
	"fmt"
	"github.com/ccqstark/gdufsclub/middleware"
	"strings"
)

type Process struct {
	ProcessID     int    `gorm:"primary_key"`
	UserID        int    `gorm:"user_id"`
	ClubID        int    `gorm:"club_id"`
	Logo          string `gorm:"logo" json:"logo"`
	ClubName      string `gorm:"club_name" json:"club_name"`
	TotalProgress int    `gorm:"total_progress" json:"total_progress"`
	Progress      int    `gorm:"progress" json:"progress"`
	Result        int    `gorm:"result" json:"result"`
}

type ProcessUser struct {
	UserID int `json:"user_id"`
	Pass   int `json:"pass"`
}

type BatchUser struct {
	Interviewee []string `json:"interviewee"`
	Progress    int      `json:"progress"`
}

//获取用户所有的面试进程
func QueryProcess(userID int) ([]Process, bool) {

	var process []Process
	if result := db.Where("user_id = ?", userID).Find(&process); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []Process{}, false
	}

	return process, true

}

//提交报名表时就创建面试进程
func CreateProcess(userID int, clubID int) bool {

	var process Process
	if club, ok := QueryClubInfo(clubID); ok == true {
		process.ClubName = club.ClubName
		process.TotalProgress = club.TotalProgress
		process.Logo = club.Logo
	} else {
		return false
	}

	process.UserID = userID
	process.ClubID = clubID
	process.Progress = 1
	process.Result = 0

	if result := db.Create(&process); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}

//查询面试结果
func QueryInterviewResult(userID int, clubID int) (int, bool) {

	var process Process
	if result := db.Where("user_id = ? and club_id=?", userID, clubID).Take(&process); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return 0, false
	}

	return process.Result, true
}

//对一人的面试结果进行操作
func OperateOnePerson(clubID int, userID int, pass int) bool {

	var process Process
	if result := db.Model(&process).Where("club_id=? and user_id=?", clubID, userID).Update("result", pass);
		result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}

//批量通过面试者
func PassBatchInterviewee(batch []string, clubID int, progress int) bool {

	batchStr := strings.Join(batch, ",")
	sql := fmt.Sprintf("UPDATE process SET result=1 WHERE club_id=%d and progress=%d and user_id IN (%s);", clubID, progress, batchStr)

	if result := db.Exec(sql); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}
