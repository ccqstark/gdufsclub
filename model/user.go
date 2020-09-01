package model

import (
	"github.com/ccqstark/gdufsclub/dao"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	db = dao.GetDB()

}

type User struct {
	UserID int    `gorm:"column:user_id"`
	OpenID string `gorm:"column:open_id"`
}

func GetFirstUser() User {
	u := User{}
	db.First(&u)

	return u
}
