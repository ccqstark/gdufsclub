package controller

import (
	"github.com/ccqstark/gdufsclub/model"
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
