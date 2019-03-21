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
	data.Actions = []models.DrawedAction{}
	data.Actions[0] = models.DrawedAction{
		AwardID:     1,
		BackMoney:   100,
		ID:          1,
		Memo:        "memo",
		PeopleCount: 2,
		Status:      "ToDo",
	}
	data.Actions[1] = models.DrawedAction{
		AwardID:     2,
		BackMoney:   200,
		ID:          2,
		Memo:        "memo",
		PeopleCount: 3,
		Status:      "ToDo",
	}
	data.Actions[2] = models.DrawedAction{
		AwardID:     3,
		BackMoney:   300,
		ID:          3,
		Memo:        "memo",
		PeopleCount: 4,
		Status:      "ToDo",
	}
	data.Awards = []models.Award{}
	data.Awards[0] = models.Award{
		Description: "description_1",
		ID:          1,
		Money:       100,
		Name:        "award_name_1",
		PicName:     "pic_1.png",
	}
	data.Awards[1] = models.Award{
		Description: "description_2",
		ID:          2,
		Money:       200,
		Name:        "award_name_2",
		PicName:     "pic_2.png",
	}
	data.Awards[2] = models.Award{
		Description: "description_3",
		ID:          3,
		Money:       100,
		Name:        "award_name_3",
		PicName:     "pic_3.png",
	}
	data.Awards[3] = models.Award{
		Description: "description_4",
		ID:          4,
		Money:       400,
		Name:        "award_name_4",
		PicName:     "pic_4.png",
	}
	data.Awards[4] = models.Award{
		Description: "description_5",
		ID:          5,
		Money:       500,
		Name:        "award_name_5",
		PicName:     "pic_5.png",
	}
	data.BackMoneyRecords = []models.BackMoneyRecord{}
	data.BackMoneyRecords[0] = models.BackMoneyRecord{}
	data.Count = models.Count{
		AllLuckyCount:      10,
		AllPeopleCount:     0,
		LuckyLeaderCount:   0,
		LuckyStaffCount:    0,
		NoLuckyStaffCount:  0,
		NoLuckyLeaderCount: 0,
		PoolMoney:          0,
	}
	data.Version = common.GetUUID()
	data.Users = []models.User{}
	data.Users[0] = models.User{
		ID:       1,
		IsDrawed: false,
		Level:    1,
		Name:     "user_1",
	}
	data.Users[1] = models.User{
		ID:       2,
		IsDrawed: false,
		Level:    1,
		Name:     "user_2",
	}
	data.Users[2] = models.User{
		ID:       3,
		IsDrawed: false,
		Level:    1,
		Name:     "user_3",
	}
	data.Users[3] = models.User{
		ID:       4,
		IsDrawed: false,
		Level:    1,
		Name:     "user_4",
	}
	data.Users[4] = models.User{
		ID:       5,
		IsDrawed: false,
		Level:    1,
		Name:     "user_5",
	}
	data.Users[5] = models.User{
		ID:       6,
		IsDrawed: false,
		Level:    1,
		Name:     "user_6",
	}
	data.Users[6] = models.User{
		ID:       7,
		IsDrawed: false,
		Level:    1,
		Name:     "user_7",
	}
	data.Users[7] = models.User{
		ID:       8,
		IsDrawed: false,
		Level:    1,
		Name:     "user_8",
	}
	data.Users[8] = models.User{
		ID:       9,
		IsDrawed: false,
		Level:    1,
		Name:     "user_9",
	}
	data.Users[9] = models.User{
		ID:       10,
		IsDrawed: false,
		Level:    2,
		Name:     "user_10",
	}
	data.Users[10] = models.User{
		ID:       11,
		IsDrawed: false,
		Level:    2,
		Name:     "user_11",
	}
	data.Users[11] = models.User{
		ID:       12,
		IsDrawed: false,
		Level:    2,
		Name:     "user_12",
	}
	data.Users[12] = models.User{
		ID:       13,
		IsDrawed: false,
		Level:    2,
		Name:     "user_13",
	}
	data.Users[13] = models.User{
		ID:       14,
		IsDrawed: false,
		Level:    2,
		Name:     "user_14",
	}
	data.Users[14] = models.User{
		ID:       15,
		IsDrawed: false,
		Level:    2,
		Name:     "user_15",
	}
	data.Users[15] = models.User{
		ID:       16,
		IsDrawed: false,
		Level:    1,
		Name:     "user_16",
	}
	data.Users[16] = models.User{
		ID:       17,
		IsDrawed: false,
		Level:    2,
		Name:     "user_17",
	}
	data.Users[17] = models.User{
		ID:       18,
		IsDrawed: false,
		Level:    2,
		Name:     "user_18",
	}
	data.Users[18] = models.User{
		ID:       19,
		IsDrawed: false,
		Level:    2,
		Name:     "user_19",
	}
	data.Users[19] = models.User{
		ID:       20,
		IsDrawed: false,
		Level:    2,
		Name:     "user_20",
	}
}
