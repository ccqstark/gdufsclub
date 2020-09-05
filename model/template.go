package model

import (
	"github.com/ccqstark/gdufsclub/middleware"
)

type Template struct {
	TemplateID        int    `gorm:"primary_key"`
	UserID            int    `gorm:"column:user_id"`
	TemplateName      string `gorm:"column:template_name" json:"template_name"`
	TemplateClass     string `gorm:"column:template_class" json:"template_class"`
	TemplateSex       string `gorm:"column:template_sex" json:"template_sex"`
	TemplateWechat    string `gorm:"column:template_wechat" json:"template_wechat"`
	TemplateEmail     string `gorm:"column:template_email" json:"template_email"`
	TemplatePhone     string `gorm:"column:template_phone" json:"template_phone"`
	TemplateHobby     string `gorm:"column:template_hobby" json:"template_hobby"`
	TemplateSelf      string `gorm:"column:template_self" json:"template_self"`
	TemplateAdvantage string `gorm:"column:template_advantage" json:"template_advantage"`
	TemplateImage     string `gorm:"column:template_image"`
}

//判断用户是否创建了模板
func IsTemplateExist(userID int) bool {

	var tpl Template
	// 检查错误是否为 RecordNotFound
	if db.Where("user_id=?", userID).Take(&tpl).RecordNotFound() {
		return false
	}

	return true
}

//获取模板
func QueryTemplate(userID int) (Template, bool) {

	var tpl Template
	if result := db.Where("user_id=?", userID).Take(&tpl); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return Template{}, false
	}

	return tpl, true
}

//插入新模板
func InsertNewTemplate(tpl *Template) (int, bool) {

	//插入记录
	if result := db.Create(&tpl); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return 0, false
	}

	//获取刚刚插入的记录的id
	var _id []int
	db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &_id)
	id := _id[0]

	return id, true
}

//更新模板头像
func UpdateTplProfile(id int, path string) bool {

	var tpl Template
	if result := db.Model(&tpl).Where("template_id=?", id).Update("template_image", path); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}

func UpdateTemplateInfo(tpl *Template) bool {

	if result := db.Save(&tpl); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}
