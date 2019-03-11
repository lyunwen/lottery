package api

import (
	"net/http"

	"../models"
	"github.com/gin-gonic/gin"
)

func GetAwards(c *gin.Context) {
	var dataObj = new(models.Data)
	dataObj, err := dataObj.GetData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取txt数据错误", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": dataObj.Awards})
	return
}

func GetCount(c *gin.Context) {
	var dataObj = new(models.Data)
	dataObj, err := dataObj.GetData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取txt数据错误", "data": nil})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": "0", "count": dataObj.Count})
}

func GetNextAction(c *gin.Context) {
	var dataObj = new(models.Data)
	dataObj, err := dataObj.GetData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取txt数据错误", "data": nil})
		return
	}
	for _, item := range dataObj.Actions {
		if item.Status == "ToDo" {
			c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": item})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": "1", "msg": "奖品已抽完"})
	return
}
