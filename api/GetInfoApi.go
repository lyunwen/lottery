package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../common"
)

func GetAwards(c *gin.Context)  {
	dataObj := common.GetData()
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": dataObj.Awards})
	return
}

func GetCount (c *gin.Context) {
	dataObj := common.GetData()
	c.JSON(http.StatusCreated, gin.H{"code": "0", "count": dataObj.Count})
}

func GetNextAction(c *gin.Context) {
	dataObj := common.GetData()
	for _, item := range dataObj.Actions {
		if item.Status == "ToDo" {
			c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": item})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": "1", "msg": "奖品已抽完"})
	return
}