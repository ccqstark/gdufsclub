package controller

import (
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PickFirstUser(c *gin.Context) {

	u := model.GetFirstUser()
	c.JSON(http.StatusOK, gin.H{
		"user_id": u.UserID,
		"open_id": u.OpenID,
	})

}

func Demo(c *gin.Context){
	session := sessions.Default(c)
	session.Set("user_id",66)
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"code":200,
	})
}
