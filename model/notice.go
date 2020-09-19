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
	Publish  int    `gorm:"publish"`
}

type TwoNotice struct {
	Progress       int    `json:"progress"`
	SuccessContent string `json:"success_content"`
	FailureContent string `json:"failure_content"`
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
	if result := db.Where("club_id=? and progress=? and pass=? and publish=1", clubID, progress, pass).Take(&notice); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return Notice{}, false
	}

	return notice, true
}

//插入新的公告
func InsertNewNotice(twoNotice *TwoNotice, clubID int, clubName string) bool {
	//成功的公告
	notice1 := Notice{
		ClubID:   clubID,
		ClubName: clubName,
		Progress: twoNotice.Progress,
		Pass:     1,
		Content:  twoNotice.SuccessContent,
	}
	//失败的公告
	notice2 := Notice{
		ClubID:   clubID,
		ClubName: clubName,
		Progress: twoNotice.Progress,
		Pass:     2,
		Content:  twoNotice.FailureContent,
	}

	//插入记录
	if result := db.Create(&notice1); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	if result := db.Create(&notice2); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}

//更新公告
func UpdateNotice(twoNotice *TwoNotice, clubID int) bool {

	var notice Notice
	if result := db.Model(&notice).Where("club_id=? and progress=? and pass=1", clubID, twoNotice.Progress).
		Update("content", twoNotice.SuccessContent); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	if result := db.Model(&notice).Where("club_id=? and progress=? and pass=2", clubID, twoNotice.Progress).
		Update("content", twoNotice.FailureContent); result.Error != nil {
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

//查看公告是否发布
func CheckIfNoticePublish(clubID int, nowProgress int) (bool, bool) {

	var notice Notice
	if result := db.Select("publish").Where("club_id=? and progress=?", clubID, nowProgress).Take(&notice); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false, false
	}

	if notice.Publish == 1 {
		return true, true
	} else {
		return false, true
	}
}
