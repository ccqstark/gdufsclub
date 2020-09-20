package model

import (
	"github.com/ccqstark/gdufsclub/dao"
	"github.com/ccqstark/gdufsclub/middleware"
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

type Login struct {
	Code string `json:"code"`
}

//用户授权登录
func AuthUser(openid string) (int, bool) {

	var user User
	user.OpenID = openid

	if db.Where("open_id=?", openid).Take(&user).RecordNotFound() {
		user.OpenID = openid
		if result := db.Create(&user); result.Error != nil {
			middleware.Log.Error(result.Error.Error())
			return 0, false
		}

		//获取刚刚插入的记录的id
		var _id []int
		db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &_id)
		id := _id[0]

		return id, true
	}

	user.OpenID = openid
	db.Where("open_id=?", openid).Take(&user)

	return user.UserID, true
}
