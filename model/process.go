package model

import (
	"github.com/ccqstark/gdufsclub/middleware"
)

type Process struct {
	ProcessID int    `gorm:"primary_key"`
	UserID    int    `gorm:"user_id"`
	ClubID    int    `gorm:"club_id"`
	ClubName  string `gorm:"club_name"`
	Progress  int    `gorm:"progress" json:"progress"`
	Result    int    `gorm:"result" json:"result"`
}

type ProcessUser struct {
	UserID int `json:"user_id"`
	Pass   int `json:"pass"`
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

//提交报表时就创建面试进程
func CreateProcess(userID int, clubID int) bool {

	var process Process
	if clubName, ok := QueryClubName(clubID); ok == true {
		process.ClubName = clubName
	} else {
		return false
	}

	process.UserID = userID
	process.ClubID = clubID
	process.Progress = 0
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
	if result := db.Where("club_id=? and user_id=?", clubID, userID).Take(&process); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}
	process.Progress = process.Progress + 1
	process.Result = pass

	if result := db.Save(&process); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}
