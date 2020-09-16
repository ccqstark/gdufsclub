package model

import (
	"fmt"
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

//判断公告存在与否
func IsNoticeExist(clubID int, progress int, pass int) bool {

	var notice Notice
	// 检查错误是否为 RecordNotFound
	if db.Where("club_id=? and progress=? and pass=?", clubID, progress, pass).Take(&notice).RecordNotFound() {
		return false
	}

	return true
}

//社团获取公告
func ClubQueryNotice(clubID int, progress int) ([]Notice, bool) {

	var notice []Notice
	if result := db.Where("club_id=? and progress=?", clubID, progress).Find(&notice); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []Notice{}, false
	}

	return notice, true
}

//用户查看公告
func QueryNotice(clubID int, progress int, pass int) (Notice, bool) {

	var notice Notice
	if result := db.Where("club_id=? and progress=? and pass=? and publish=?", clubID, progress, pass, 1).Take(&notice); result.Error != nil {
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

//更新公告
func UpdateNotice(notice *Notice) bool {

	if result := db.Save(&notice); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}

//公告统一发布
func MakeNoticePublished(clubID int, progress int) bool {

	sql1 := fmt.Sprintf("UPDATE notice SET publish=1 WHERE club_id=%d and progress=%d;", clubID, progress)
	if result := db.Exec(sql1); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	//面试轮数+1
	sql2 := fmt.Sprintf("UPDATE process SET progress=%d,result=0 WHERE club_id=%d and progress=%d;", progress+1, clubID, progress)
	if result := db.Exec(sql2); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}
