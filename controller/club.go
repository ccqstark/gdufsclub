package controller

import (
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SettleNewClub(c *gin.Context) {
	var club model.Club
	if err := c.ShouldBind(&club); err != nil {
		c.String(http.StatusBadRequest, "%v", err)
	} else {
		//用户名是否重复
		if model.IsAccountRepeat(club.ClubAccount) {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "此用户名已被注册",
			})
		} else {
			//插入数据
			if clubId, ok := model.InsertNewClub(&club); ok == true {
				//club_id存入session
				session := sessions.Default(c)
				session.Set("club_id", clubId)
				session.Save()

				c.JSON(http.StatusOK, gin.H{
					"code":    1,
					"msg":     "申请提交成功，请耐心等待后台审核",
					"club_id": clubId,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 0,
					"msg":  "注册失败，请重试",
				})
			}
		}
	}
}
