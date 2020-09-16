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
	UserID int    `gorm:"column:submitter_id"`
	Name   string `gorm:"column:name"`
	Sex    string `gorm:"column:sex"`
	Class  string `gorm:"column:class"`
	Phone  string `gorm:"column:phone"`
	Wechat string `gorm:"column:wechat"`
	Result int    `gorm:"column:result"`
}

type ClubAccount struct {
	Account  string `json:"account"`
	Password string `json:"password"`
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

//当前面用户总页数
//func QueryUserTotalPage(clubID int, progress int) (int, bool) {
//
//	var process Process
//	total := 0
//	if result := db.Model(&process).Where("club_ID=? and progress=?", clubID, progress).Count(&total); result.Error != nil {
//		middleware.Log.Error(result.Error.Error())
//		return 0, false
//	}
//
//	var pageFloat float32 = float32(total / recordPerPage)
//	var pageInt float32 = float32(int(pageFloat))
//	if (pageFloat - pageInt) > 0{
//		return int(pageInt), true
//	} else {
//		return int(pageInt+1), true
//	}
//}

//生成这一轮用户通过者用户基本信息列表
func QueryUserListBrief(clubID int, progress int) ([]UserList, bool) {

	//基本信息: 姓名，性别，班级，手机号，微信号
	var userList []UserList
	//原生sql子查询，获取当前面的面试者列表
	sql := "SELECT b.submitter_id, b.name, b.sex, b.class, b.phone, " +
		"b.wechat, a.result FROM process a, resume b " +
		"WHERE a.user_id = b.submitter_id AND a.club_id = ? AND a.progress = ?;"
	if result := db.Raw(sql, clubID, progress).Scan(&userList); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []UserList{}, false
	}

	//计算页面第一条和最后一条文位置
	//startRecord := (page-1) * recordPerPage
	//endRecord := startRecord + recordPerPage
	////分页，用切片截取
	//if startRecord > len(userList){
	//	return []UserList{}, false
	//}else if endRecord > len(userList){
	//	userList = userList[startRecord:]
	//} else {
	//	userList = userList[startRecord:endRecord]
	//}

	return userList, true
}

//获取当前面所有通过者ID
func QueryPasser(clubID int, progress int) []int {

	var process []Process
	if result := db.Where("club_ID=? and progress=? and result=?", clubID, progress, 1).
		Select("user_id").Find(&process); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []int{}
	}

	var passerID []int

	for _, pro := range process {
		passerID = append(passerID, pro.UserID)
	}

	return passerID
}

//通过ID数组批量获取用户提交的报名表上的信息
func GainInfoByArray(clubID int, userID []int) ([]Resume, bool) {

	var resumeArr []Resume
	if result := db.Where("club_id=? and submitter_id IN (?)", clubID, userID).
		Find(&resumeArr); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []Resume{}, false
	}

	return resumeArr, true
}

//登录时判断密码
func JudgePassword(account string, password string) (Club, int) {
	var club Club
	if result := db.Where("club_account=?", account).Take(&club); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return Club{}, 0
	}

	//0 错误   1 可以登录   2 审核不通过
	if util.Md5SaltCrypt(password) == club.ClubPassword {
		if club.Pass == 1 {
			return club, 1
		} else {
			return club, 2
		}
	} else {
		return Club{}, 0
	}

}
