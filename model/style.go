package model

import "github.com/ccqstark/gdufsclub/middleware"

type Style struct {
	StyleID        int    `gorm:"primary_key"`
	ClubID         int    `gorm:"column:club_id" json:"club_id"`
	ClubName       string `gorm:"column:club_name" json:"club_name"`
	StyleName      uint8  `gorm:"column:style_name" json:"style_name"`
	StyleSex       uint8  `gorm:"column:style_sex" json:"style_sex"`
	StyleClass     uint8  `gorm:"column:style_class" json:"club_class"`
	StylePhone     uint8  `gorm:"column:style_phone" json:"style_phone"`
	StyleWechat    uint8  `gorm:"column:style_wechat" json:"style_wechat"`
	StyleImage     uint8  `gorm:"column:style_image" json:"style_image"`
	StyleHobby     uint8  `gorm:"column:style_hobby" json:"style_hobby"`
	StyleAdvantage uint8  `gorm:"column:style_advantage" json:"style_advantage"`
	StyleSelf      uint8  `gorm:"column:style_self" json:"style_self"`
	StyleReason    uint8  `gorm:"column:style_reason" json:"style_reason"`
	StyleExtra     string `gorm:"column:style_extra" json:"style_extra"`
}


func InsertNewStyle(style *Style)(int,bool){

	//插入记录
	if result := db.Create(&style); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return 0, false
	}

	//获取刚刚插入的记录的id
	var _id []int
	db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &_id)
	id := _id[0]

	//方法判断插入成功返回false
	if !db.NewRecord(&style) {
		return id, true
	} else {
		return 0, false
	}
}


