package api

import (
	"net/http"
	"strconv"
	"time"

	"../models"
	"github.com/gin-gonic/gin"
)

func InitData(c *gin.Context) {

	pswd := c.Query("pswd")
	if pswd != "123456" {
		c.JSON(http.StatusCreated, gin.H{"code": "1", "msg": "密码不正确"})
		return
	}
	var dataObj = new(models.Data)
	dataObj, err := dataObj.GetData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取txt数据错误", "data": nil})
		return
	}

	if dataObj.Awards == nil {
		dataObj.Awards = []models.Award{}
	} else {
		for i := range dataObj.Awards {
			dataObj.Awards[i].ID = i + 1
		}
	}

	if dataObj.Actions == nil {
		dataObj.Actions = []models.DrawedAction{}
	} else {
		for i := range dataObj.Actions {
			dataObj.Actions[i].Status = "ToDo"
			dataObj.Actions[i].ID = i + 1
		}
	}
	if dataObj.Users == nil {
		dataObj.Users = []models.User{}
	} else {
		for i := range dataObj.Users {
			dataObj.Users[i].ID = i + 1
			dataObj.Users[i].IsDrawed = false
		}
	}
	dataObj.DrawedRecords = []models.DrawedRecord{}
	dataObj.BackMoneyRecords = []models.BackMoneyRecord{}

	dataObj.Count.PoolMoney = 0
	dataObj.Count.AllLuckyCount = 0
	dataObj.Count.AllPeopleCount = 0
	dataObj.Count.LuckyStaffCount = 0
	dataObj.Count.LuckyLeaderCount = 0
	dataObj.Count.NoLuckyLeaderCount = 0
	dataObj.Count.NoLuckyStaffCount = 0
	for _, item := range dataObj.Users {
		dataObj.Count.AllPeopleCount++
		if item.Level == 1 {
			if item.IsDrawed {
				dataObj.Count.LuckyStaffCount++
			} else {
				dataObj.Count.NoLuckyStaffCount++
			}
		} else if item.Level == 2 {
			if item.IsDrawed {
				dataObj.Count.LuckyLeaderCount++
			} else {
				dataObj.Count.NoLuckyLeaderCount++
			}
		}
	}
	setDataErr := new(models.Data).SetData()
	if setDataErr == nil {
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": setDataErr.Error()})
		return
	}
}

func AddPoolMoney(c *gin.Context) {
	money, err := strconv.Atoi(c.Query("money"))
	if money < 1 || err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "1", "msg": "加奖金额不正确"})
		return
	}
	memo := c.Query("memo")
	var dataObj = new(models.Data)
	dataObj, err = dataObj.GetData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取txt数据错误", "data": nil})
		return
	}
	dataObj.BackMoneyRecords = append(dataObj.BackMoneyRecords, models.BackMoneyRecord{
		BackTime: time.Now().Format("2006-01-02 15:04:05"),
		Memo:     "临时加奖:" + memo,
		Money:    money,
		UserID:   -1,
		UserName: "",
	})

	setDataErr := new(models.Data).SetData()
	if setDataErr == nil {
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": setDataErr.Error()})
		return
	}
}
