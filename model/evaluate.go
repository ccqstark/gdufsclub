package model

import "github.com/ccqstark/gdufsclub/middleware"

type Evaluate struct {
	EvaluateID int    `gorm:"primary_key"`
	UserID     int    `gorm:"column:user_id" json:"user_id"`
	ClubID     int    `gorm:"column:club_id"`
	Department string `gorm:"column:department"`
	Progress   int    `gorm:"column:progress" json:"progress"`
	Content    string `gorm:"column:content" json:"content"`
}

//获取评价
func QueryEvaluate(clubID int, userID int, progress int, department string) (Evaluate, bool) {

	var evaluate Evaluate
	if result := db.Where("club_id=? and user_id=? and progress=? and department=?", clubID, userID, progress, department).Take(&evaluate); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return Evaluate{}, false
	}

	return evaluate, true
}

//创建评价
func CreateEvaluate(evaluate Evaluate) (int, bool) {

	if result := db.Create(&evaluate); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return 0, false
	}

	//获取刚刚插入的记录的id
	var _id []int
	db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &_id)
	id := _id[0]

	return id, true
}

//修改评价
func UpdateEvaluate(evaluate *Evaluate) bool {

	if result := db.Model(&evaluate).Update("content", evaluate.Content); result.Error != nil {
		middleware.Log.Error(result.Error.Error())
		return false
	}

	return true
}
