package model

import (
	"github.com/ccqstark/gdufsclub/middleware"
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
	if result := db.Model(&club).Where("club_id=?", clubID).Update("pass",status); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}
