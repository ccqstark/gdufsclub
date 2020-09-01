package model

import (
	"errors"
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/util"
)

var (
	//ErrRecordNotFound record not found error
	ErrRecordNotFound = errors.New("record not found")
)

type Club struct {
	ClubID        int    `gorm:"primary_key"`
	ClubName      string `gorm:"column:club_name" json:"club_name"`
	ClubEmail     string `gorm:"column:club_email" json:"club_email"`
	ClubWechat    string `gorm:"column:club_wechat" json:"club_wechat"`
	ClubPhone     string `gorm:"column:club_phone" json:"club_phone"`
	ClubAccount   string `gorm:"column:club_account" json:"club_account"`
	ClubPassword  string `gorm:"column:club_password" json:"club_password"`
	TotalProgress int    `gorm:"column:total_progress" json:"total_progress"`
	Logo          string `gorm:"column:logo"`
	Pass          int    `gorm:"column:pass"`
}

func InsertNewClub(club *Club) (int, bool) {
	//md5加密
	club.ClubPassword = util.Md5SaltCrypt(club.ClubPassword)
	//插入记录
	if result := db.Create(&club); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return 0, false
	}

	//获取刚刚插入的记录的id
	var _id []int
	db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &_id)
	id := _id[0]

	//方法判断插入成功返回false
	if !db.NewRecord(&club) {
		return id, true
	} else {
		return 0, false
	}
}

//判断账户名重复与否
func IsAccountRepeat(accountStr string) bool {

	var club Club
	// 检查错误是否为 RecordNotFound
	err := db.Where("club_account=?", accountStr).Take(&club).Error
	if errors.Is(err, ErrRecordNotFound) {
		return false
	} else if err != nil {
		middleware.Log.Error(err.Error())
		return true
	}

	return true
}

//更新logo地址
func UpdateLogo(id int, path string) bool {

	var club Club
	if result := db.Model(&club).Where("club_id=?", id).Update("logo", path); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}
