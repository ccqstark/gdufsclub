package model

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
	Pass          string `gorm:"column:pass"`
}

func InsertNewClub(club *Club) bool {
	db.Create(&club)

	//判断插入成功返回false
	if !db.NewRecord(&club) {
		return true
	}
	return false

}
