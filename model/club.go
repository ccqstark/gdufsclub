package model

import (
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/util"
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

type UserList struct {
	UserID int
	Name   string
	Sex    string
	Class  string
	Phone  string
	Wechat string
}

//插入新的社团
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

	return id, true
}

//判断账户名重复与否
func IsAccountRepeat(accountStr string) bool {

	var club Club
	if db.Where("club_account=?", accountStr).Take(&club).RecordNotFound() {
		return false
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

//通过关键词搜索
func SearchByWord(cutWord []string) []Club {

	var clubSegment []Club
	var clubGather []Club
	for _, word := range cutWord {
		db.Select("club_id,club_name,total_progress,logo").Where("club_name like ? and pass=?", "%"+word+"%", 1).Find(&clubSegment)
		clubGather = append(clubGather, clubSegment...)
	}

	return clubGather
}

func QueryUserTotalPage(clubID int, progress int) (int, bool) {

	var process Process
	total := 0
	if result := db.Model(&process).Where("club_ID=? and progress=? and pass <> ?", clubID, progress-1, 2).Count(&total); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return 0, false
	}

	return total, true
}

//生成这一轮用户通过者用户基本信息列表
func QueryUserListBrief(clubID int, progress int,numberRows int,begin int) ([]UserList, bool) {

	//基本信息: 姓名，性别，班级，手机号，微信号
	//先获取通过了的id, 再通过id查信息, 分页
	var userList []UserList
	if result := db.Table("process").Where("club_ID=? and progress=? and pass <> ?", clubID, progress-1, 2).
		Joins("left join resume on resume.user_id=process.user_id").
		Select("user_id,name,sex,class,phone,wechat").
		Order("user_id").Limit(numberRows).Offset(begin).Find(&userList); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []UserList{}, false
	}

	return userList, true
}
