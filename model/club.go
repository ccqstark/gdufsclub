package model

import (
	"fmt"
	"github.com/ccqstark/gdufsclub/middleware"
	"github.com/ccqstark/gdufsclub/util"
	"strings"
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
	Department    string `gorm:"column:department"`
}

type UserList struct {
	ResumeID   int    `gorm:"column:resume_id"`
	UserID     int    `gorm:"column:submitter_id"`
	ClubID     int    `gorm:"column:club_id"`
	Name       string `gorm:"column:name"`
	Sex        string `gorm:"column:sex"`
	Class      string `gorm:"column:class"`
	Phone      string `gorm:"column:phone"`
	Wechat     string `gorm:"column:wechat"`
	Result     int    `gorm:"column:result"`
	Department string `gorm:"column:department"`
}

type ClubAccount struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type ClubModInfo struct {
	ClubName      string `json:"club_name"`
	ClubEmail     string `json:"club_email"`
	ClubWechat    string `json:"club_wechat"`
	ClubPhone     string `json:"club_phone"`
	TotalProgress int    `json:"total_progress"`
	Department    string `gorm:"column:department"`
}

type NewField struct {
	StudentNumber  string
	PoliticsStatus string
	Campus         string
	TakeOffice     string
	Remarks        string
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

	//club表更新
	sql := fmt.Sprintf("UPDATE club SET logo='%s' WHERE club_id=%d", path, id)

	if result := db.Exec(sql); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	//process表更新
	sql = fmt.Sprintf("UPDATE process SET logo='%s' WHERE club_id=%d", path, id)

	if result := db.Exec(sql); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}

//用社团id查找社团名
func QueryClubName(id int) (string, bool) {

	var club Club
	if result := db.Where("club_id=?", id).Take(&club); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return "", false
	}

	return club.ClubName, true
}

//用社团ID获取社团总信息
func QueryClubInfo(id int) (Club, bool) {

	var club Club
	if result := db.Where("club_id=?", id).Take(&club); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return Club{}, false
	}

	return club, true
}

//查找所有社团信息
func QueryAllClubInfo() ([]Club, bool) {

	var club []Club
	if result := db.Select("club_id,club_name,club_email,club_phone,club_wechat,total_progress,logo,department").
		Where("pass=?", 1).Find(&club); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []Club{}, false
	}

	return club, true
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

//获取当前面所有通过者ID
func QueryPasser(clubID int, progress int) []int {

	var process []Process
	var process2 []Process
	if result := db.Where("club_ID=? and progress=? and result=?", clubID, progress, 1).
		Select("resume_id").Find(&process); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []int{}
	}

	if result := db.Where("club_ID=? and progress>?", clubID, progress).
		Select("resume_id").Find(&process2); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []int{}
	}
	//用于存储所有通过者的id
	var passerID []int

	for _, pro := range process {
		passerID = append(passerID, pro.ResumeID)
	}

	for _, pro := range process2 {
		passerID = append(passerID, pro.ResumeID)
	}

	return passerID
}

//获取经历过当前面所有人的ID
func QueryAllPerson(clubID int, progress int) []int {

	var process []Process
	if result := db.Where("club_ID=? and progress>=?", clubID, progress).
		Select("resume_id").Find(&process); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []int{}
	}

	//用于存储所有人的id
	var personID []int
	for _, pro := range process {
		personID = append(personID, pro.ResumeID)
	}

	return personID
}

//通过ID数组批量获取用户提交的报名表上的信息
func GainInfoByArray(clubID int, ResumeID []int) ([]Resume, bool) {

	var resumeArr []Resume
	if result := db.Where("club_id=? and resume_id IN (?)", clubID, ResumeID).
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

//修改社团信息
func UpdateClubInfo(info ClubModInfo, clubID int) bool {

	var club Club
	if result := db.Model(&club).Where("club_id=?", clubID).
		Updates(map[string]interface{}{
			"club_name":      info.ClubName,
			"club_phone":     info.ClubPhone,
			"club_email":     info.ClubEmail,
			"club_wechat":    info.ClubWechat,
			"total_progress": info.TotalProgress,
			"department":     info.Department,
		}); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	db.Model(Process{}).Where("club_id=?", clubID).Updates(map[string]interface{}{
		"club_name":      info.ClubName,
		"total_progress": info.TotalProgress,
	})
	db.Model(Style{}).Where("club_id=?", clubID).Update("club_name", info.ClubName)
	db.Model(Notice{}).Where("club_id=?", clubID).Update("club_name", info.ClubName)

	return true
}

//生成这一轮的用户的基本信息列表
func QueryUserListBrief(clubID int, progress int) ([]UserList, bool) {

	//基本信息: 姓名，性别，班级，手机号，微信号
	var userList []UserList
	var resumeIDArr []int
	type Parr struct {
		ProcessID  int
		ResumeID   int
		UserID     int
		Result     int
		Department string
	}
	//先获取对应的用户的id
	sql := "SELECT process_id,resume_id,user_id,result,department FROM process WHERE club_id=? and progress=?"
	var parr []Parr
	if result := db.Raw(sql, clubID, progress).Scan(&parr); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []UserList{}, false
	}

	//转为id数组
	for i := range parr {
		resumeIDArr = append(resumeIDArr, parr[i].ResumeID)
	}

	//用用户id查询信息
	sql = "SELECT resume_id, submitter_id, club_id, name, sex, class, phone, wechat, department FROM resume WHERE club_id=? and resume_id IN (?)"
	if result := db.Raw(sql, clubID, resumeIDArr).Scan(&userList); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []UserList{}, false
	}

	//补充面试结果信息,部门
	for i1 := range userList {
		for i2 := range parr {
			if userList[i1].ResumeID == parr[i2].ResumeID {
				userList[i1].Result = parr[i2].Result
			}
		}
	}

	return userList, true
}

//政治面貌,所在校区,任职,备注
func ParseField(str string) NewField {

	strArr := strings.Split(str, "@@")

	keyArr := strings.Split(strArr[0], ",")
	valueArr := strings.Split(strArr[1], "**")

	newField := NewField{}

	for i := range keyArr {
		switch keyArr[i] {
		case "学号":
			newField.StudentNumber = valueArr[i]
		case "政治面貌":
			newField.PoliticsStatus = valueArr[i]
		case "所在校区":
			newField.Campus = valueArr[i]
		}
	}

	return newField
}
