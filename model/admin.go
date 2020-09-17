package model

import (
	"github.com/ccqstark/gdufsclub/middleware"
	"strings"
)

//所有还未审核的社团
func QueryAllNotPass() ([]Club, bool) {

	var club []Club
	if result := db.Where("pass=?", 0).Find(&club); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return club, false
	}

	return club, true
}

//所有通过的社团
func QueryAllPass() ([]Club, bool) {

	var club []Club
	if result := db.Where("pass=?", 1).Find(&club); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return club, false
	}

	return club, true
}

//审核通过
func AuditOneClub(clubID int, status int) bool {

	var club Club
	if result := db.Model(&club).Where("club_id=?", clubID).Update("pass", status); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}

//自定义字段
func QueryAllCustomField() ([]string, bool) {

	var field []string
	var style []Style
	if result := db.Select("style_extra").Find(&style); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return []string{}, false
	}
	//所有字段放在一个数组里
	for _, v := range style {
		fieldSegment := v.StyleExtra
		fieldSlice := strings.Split(fieldSegment, ",")
		for _, fs := range fieldSlice {
			field = append(field, fs)
		}
	}

	return field, true
}
