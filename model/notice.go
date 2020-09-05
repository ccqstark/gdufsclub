package model

import (
	"github.com/ccqstark/gdufsclub/middleware"
)

type Notice struct {
	NoticeID int    `gorm:"primary_key"`
	ClubID   int    `gorm:"club_id"`
	ClubName string `gorm:"club_name"`
	Progress int    `gorm:"progress" json:"progress"`
	Pass     int    `gorm:"pass" json:"pass"`
	Content  string `gorm:"content" json:"content"`
}

func IsNoticeExist(clubID int, progress int, pass int) bool {

	var notice Notice
	// 检查错误是否为 RecordNotFound
	if db.Where("club_id=? and progress=? and pass=?", clubID, progress, pass).Take(&notice).RecordNotFound() {
		return false
	}

	return true
}

func QueryNotice(clubID int, progress int, pass int) (Notice, bool) {

	var notice Notice
	if result := db.Where("club_id=? and progress=? pass=?", clubID, progress, pass).Take(&notice); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return Notice{}, false
	}

	return notice, true
}

//插入新的公告
func InsertNewNotice(notice *Notice) (int, bool) {

	//插入记录
	if result := db.Create(&notice); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return 0, false
	}

	//获取刚刚插入的记录的id
	var _id []int
	db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &_id)
	id := _id[0]

	return id, true
}

func UpdateNotice(notice *Notice) bool {

	if result := db.Save(&notice); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}
