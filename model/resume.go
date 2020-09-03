package model

import (
	"github.com/ccqstark/gdufsclub/middleware"
)

type Resume struct {
	ResumeID    int    `gorm:"primary_key"`
	SubmitterID int    `gorm:"column:submitter_id"`
	ClubID      int    `gorm:"column:club_uid json:"club_id""`
	Name        string `gorm:"column:name" json:"name"`
	Sex         string `gorm:"column:sex" json:"sex"`
	Class       string `gorm:"column:class" json:"class"`
	Phone       string `gorm:"column:phone" json:"phone"`
	Email       string `gorm:"column:email" json:"email"`
	Wechat      string `gorm:"column:wechat" json:"wechat"`
	Hobby       string `gorm:"column:hobby" json:"hobby"`
	Advantage   string `gorm:"column:advantage" json:"advantage"`
	Self        string `gorm:"column:self" json:"self"`
	Reason      string `gorm:"column:reason" json:"reason"`
	Image       string `gorm:"column:image" json:"image"`
	Extra       string `gorm:"column:extra" json:"extra"`
}

func InsertNewResume(resume *Resume) (int, bool) {

	//插入记录
	if result := db.Create(&resume); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return 0, false
	}

	//获取刚刚插入的记录的id
	var _id []int
	db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &_id)
	id := _id[0]

	//方法判断插入成功返回false
	if !db.NewRecord(&resume) {
		return id, true
	} else {
		return 0, false
	}
}

func UpdateResumeProfile(id int, path string) bool {

	var resume Resume
	if result := db.Model(&resume).Where("resume_id=?", id).Update("image", path); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}

func QueryResume(userID int, clubID int) (Resume, bool) {

	var resume Resume
	if result := db.Where("user_id=? and club_id=?", userID, clubID).Take(&resume); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return Resume{}, false
	}

	return resume, true
}

func UpdateResumeInfo(resume *Resume)bool{

	if result := db.Save(&resume); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}
