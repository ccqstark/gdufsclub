package model

import (
	"github.com/ccqstark/gdufsclub/middleware"
)

func QueryAllNotPass()([]Club,bool){

	var club []Club
	if result:=db.Where("pass=?",0).Find(&club);result.Error!=nil {
		middleware.Log.Error(result.Error.Error())
		return club,false
	}

	return club,true
}