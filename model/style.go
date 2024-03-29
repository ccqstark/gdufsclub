package model

import "github.com/ccqstark/gdufsclub/middleware"

type Style struct {
	StyleID         int    `gorm:"primary_key"`
	ClubID          int    `gorm:"column:club_id" json:"club_id"`
	ClubName        string `gorm:"column:club_name" json:"club_name"`
	StyleDepartment uint8  `gorm:"column:style_department" json:"style_department"`
	StyleVolunteer  uint8  `gorm:"column:style_volunteer" json:"style_volunteer"`
	StyleName       uint8  `gorm:"column:style_name" json:"style_name"`
	StyleSex        uint8  `gorm:"column:style_sex" json:"style_sex"`
	StyleCollege    uint8  `gorm:"column:style_college" json:"style_college"`
	StyleClass      uint8  `gorm:"column:style_class" json:"club_class"`
	StylePhone      uint8  `gorm:"column:style_phone" json:"style_phone"`
	StyleEmail      uint   `gorm:"column:style_email" json:"style_email"`
	StyleWechat     uint8  `gorm:"column:style_wechat" json:"style_wechat"`
	StyleImage      uint8  `gorm:"column:style_image" json:"style_image"`
	StyleHobby      uint8  `gorm:"column:style_hobby" json:"style_hobby"`
	StyleAdvantage  uint8  `gorm:"column:style_advantage" json:"style_advantage"`
	StyleSelf       uint8  `gorm:"column:style_self" json:"style_self"`
	StyleReason     uint8  `gorm:"column:style_reason" json:"style_reason"`
	StyleExtra      string `gorm:"column:style_extra" json:"style_extra"`
}

//插入新的表样式
func InsertNewStyle(style *Style) (int, bool) {

	//插入记录
	style.StyleClass = 1
	if result := db.Create(&style); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return 0, false
	}

	//获取刚刚插入的记录的id
	var _id []int
	db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &_id)
	id := _id[0]

	return id, true
}

//判断样式表是否存在
func IsStyleExist(clubID int) bool {

	var style Style
	// 检查错误是否为 RecordNotFound
	if db.Where("club_id=?", clubID).Take(&style).RecordNotFound() {
		return false
	}

	return true
}

//获取表样式
func QueryStyle(clubID int) (Style, bool) {

	var style Style
	if result := db.Where("club_id=?", clubID).Take(&style); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return Style{}, false
	}

	return style, true
}

//更新样式表
func UpdateStyle(style *Style) bool {

	style.StyleClass = 1
	if result := db.Save(&style); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}

//查询社团的部门列表
func GetDepartmentList(clubID int) string {
	var club Club
	if result := db.Where("club_id=?", clubID).Take(&club); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return ""
	}

	return club.Department
}
