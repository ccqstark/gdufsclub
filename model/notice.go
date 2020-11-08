package model

import (
	"fmt"
	"github.com/ccqstark/gdufsclub/middleware"
)

type Notice struct {
	NoticeID   int    `gorm:"primary_key"`
	ClubID     int    `gorm:"club_id"`
	ClubName   string `gorm:"club_name"`
	Department string `gorm:"department"`
	Progress   int    `gorm:"progress" json:"progress"`
	Pass       int    `gorm:"pass" json:"pass"`
	Content    string `gorm:"content" json:"content"`
	Publish    int    `gorm:"publish"`
}

type TwoNotice struct {
	Progress       int    `json:"progress"`
	SuccessContent string `json:"success_content"`
	FailureContent string `json:"failure_content"`
	Department     string `json:"department"`
}

//判断公告存在与否
func IsNoticeExist(clubID int, progress int, pass int, department string) bool {

	var notice Notice
	// 检查错误是否为 RecordNotFound
	if db.Where("club_id=? and progress=? and pass=? and department=?", clubID, progress, pass, department).Take(&notice).RecordNotFound() {
		return false
	}

	//内容为空也算没有公告
	if notice.Content == "" {
		return false
	}

	return true
}

//社团获取公告
func ClubQueryNotice(clubID int, progress int, department string) ([]Notice, bool) {

	var notice []Notice
	if result := db.Where("club_id=? and progress=? and department=?", clubID, progress, department).Find(&notice); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []Notice{}, false
	}

	return notice, true
}

//用户查看公告
func QueryNotice(clubID int, progress int, pass int, department string) (Notice, bool) {

	var notice Notice
	if result := db.Where("club_id=? and progress=? and pass=? and publish=1 and department=?", clubID, progress, pass, department).Take(&notice); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return Notice{}, false
	}

	return notice, true
}

//插入新的公告
func InsertNewNotice(twoNotice *TwoNotice, clubID int, clubName string, department string) bool {
	//成功的公告
	notice1 := Notice{
		ClubID:     clubID,
		ClubName:   clubName,
		Department: department,
		Progress:   twoNotice.Progress,
		Pass:       1,
		Content:    twoNotice.SuccessContent,
	}
	//失败的公告
	notice2 := Notice{
		ClubID:     clubID,
		ClubName:   clubName,
		Department: department,
		Progress:   twoNotice.Progress,
		Pass:       2,
		Content:    twoNotice.FailureContent,
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
	if result := db.Model(&notice).Where("club_id=? and progress=? and pass=1 and department=?", clubID, twoNotice.Progress, twoNotice.Department).
		Update("content", twoNotice.SuccessContent); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	if result := db.Model(&notice).Where("club_id=? and progress=? and pass=2 and department=?", clubID, twoNotice.Progress, twoNotice.Department).
		Update("content", twoNotice.FailureContent); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}

//公告统一发布
func MakeNoticePublished(clubID int, progress int, department string) bool {

	sql1 := fmt.Sprintf("UPDATE notice SET publish=1 WHERE club_id=%d and progress=%d and department='%s';", clubID, progress, department)
	if result := db.Exec(sql1); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	//通过的人面试轮数+1
	sql2 := fmt.Sprintf("UPDATE process SET progress=%d, result=0 WHERE club_id=%d and progress=%d and result=1 and department='%s';", progress+1, clubID, progress, department)
	if result := db.Exec(sql2); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	//没审核的人直接算不过
	sql3 := fmt.Sprintf("UPDATE process SET result=2 WHERE club_id=%d and progress=%d and result=0 and department='%s';", clubID, progress, department)
	if result := db.Exec(sql3); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}

//查看公告是否发布
func CheckIfNoticePublish(clubID int, nowProgress int, department string) (bool, bool) {

	var notice Notice
	if result := db.Select("publish").Where("club_id=? and progress=? and department=?", clubID, nowProgress, department).Take(&notice); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false, false
	}

	if notice.Publish == 1 {
		return true, true
	} else {
		return false, true
	}
}
