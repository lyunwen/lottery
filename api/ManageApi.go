package api

import (
	"net/http"
	"strconv"
	"time"

	"../common"
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
	setDataErr := dataObj.SetData()
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

	setDataErr := dataObj.SetData()
	if setDataErr == nil {
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": setDataErr.Error()})
		return
	}
}

func InitSystem(c *gin.Context) {
	pswd := c.Query("pswd")
	if pswd != "123456" {
		c.JSON(http.StatusCreated, gin.H{"code": "1", "msg": "密码不正确"})
		return
	}
	var data = new(models.Data)
	data.Description = "自动生成时间" + time.Now().Format("2006-01-02 15:04:05")
	data.Actions = []models.DrawedAction{
		{
			AwardID:     1,
			BackMoney:   100,
			ID:          1,
			Memo:        "memo",
			PeopleCount: 1,
			Status:      "ToDo",
		},
		{
			AwardID:     2,
			BackMoney:   200,
			ID:          2,
			Memo:        "memo",
			PeopleCount: 2,
			Status:      "ToDo",
		},
		{
			AwardID:     3,
			BackMoney:   300,
			ID:          3,
			Memo:        "memo",
			PeopleCount: 3,
			Status:      "ToDo",
		},
	}
	data.Awards = []models.Award{
		{
			Description: "description_1",
			ID:          1,
			Money:       100,
			Name:        "award_name_1",
			PicName:     "pic_1.png",
		},
		{
			Description: "description_2",
			ID:          2,
			Money:       200,
			Name:        "award_name_2",
			PicName:     "pic_2.png",
		},
		{
			Description: "description_3",
			ID:          3,
			Money:       300,
			Name:        "award_name_3",
			PicName:     "pic_3.png",
		},
		{
			Description: "description_4",
			ID:          4,
			Money:       400,
			Name:        "award_name_4",
			PicName:     "pic_4.png",
		},
		{
			Description: "description_5",
			ID:          5,
			Money:       500,
			Name:        "award_name_5",
			PicName:     "pic_5.png",
		},
	}
	data.BackMoneyRecords = []models.BackMoneyRecord{}
	data.DrawedRecords = []models.DrawedRecord{}
	data.Users = []models.User{
		{
			ID:       1,
			IsDrawed: false,
			Level:    1,
			Name:     "user_1",
		},
		{
			ID:       2,
			IsDrawed: false,
			Level:    1,
			Name:     "user_2",
		},
		{
			ID:       3,
			IsDrawed: false,
			Level:    1,
			Name:     "user_3",
		},
		{
			ID:       4,
			IsDrawed: false,
			Level:    1,
			Name:     "user_4",
		},
		{
			ID:       5,
			IsDrawed: false,
			Level:    1,
			Name:     "user_5",
		},
		{
			ID:       6,
			IsDrawed: false,
			Level:    1,
			Name:     "user_6",
		},
		{
			ID:       7,
			IsDrawed: false,
			Level:    1,
			Name:     "user_7",
		},
		{
			ID:       8,
			IsDrawed: false,
			Level:    1,
			Name:     "user_8",
		},
		{
			ID:       9,
			IsDrawed: false,
			Level:    1,
			Name:     "user_9",
		},
		{
			ID:       10,
			IsDrawed: false,
			Level:    2,
			Name:     "user_10",
		},
		{
			ID:       11,
			IsDrawed: false,
			Level:    2,
			Name:     "user_11",
		},
		{
			ID:       12,
			IsDrawed: false,
			Level:    2,
			Name:     "user_12",
		},
		{
			ID:       13,
			IsDrawed: false,
			Level:    2,
			Name:     "user_13",
		},
		{
			ID:       14,
			IsDrawed: false,
			Level:    2,
			Name:     "user_14",
		},
		{
			ID:       15,
			IsDrawed: false,
			Level:    2,
			Name:     "user_15",
		},
		{
			ID:       16,
			IsDrawed: false,
			Level:    1,
			Name:     "user_16",
		},
		{
			ID:       17,
			IsDrawed: false,
			Level:    2,
			Name:     "user_17",
		},
		{
			ID:       18,
			IsDrawed: false,
			Level:    2,
			Name:     "user_18",
		},
		{
			ID:       19,
			IsDrawed: false,
			Level:    2,
			Name:     "user_19",
		},
		{
			ID:       20,
			IsDrawed: false,
			Level:    2,
			Name:     "user_20",
		},
	}
	data.Version = common.GetUUID()
	err := data.ReSetData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": err.Error()})
	} else {

		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok"})
	}
}
