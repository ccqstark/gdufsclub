package controller

import (
	"github.com/ccqstark/gdufsclub/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SettleNewClub(c *gin.Context) {
	var club model.Club
	if err := c.ShouldBind(&club); err != nil {
		c.String(http.StatusBadRequest, "%v", err)
	}else {
		//插入数据
		if model.InsertNewClub(&club){
			c.JSON(http.StatusOK,gin.H{
				"code": 1,
			})
		}else {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
			})
		}
	}

}
