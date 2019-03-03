package api

import (
	"net/http"
	"strconv"
	"time"

	"../common"
	"../models"
	"github.com/gin-gonic/gin"
)

func NDraw(c *gin.Context) {
	drawer := c.Query("drawer")
	memo := c.Query("memo")
	leaderCount, err := strconv.Atoi(c.Query("leaderCount"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "领导人数人数错误", "data": nil})
		return
	}
	staffCount, err := strconv.Atoi(c.Query("staffCount"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "员工人数错误", "data": nil})
		return
	}
	dataObj := common.GetData()
	actionID := c.Query("actionID")
	var awardID int
	var allPeopleCount int
	var backMoney int
	for i, item := range dataObj.Actions {
		if item.Status == "ToDo" {
			if strconv.Itoa(item.ID) == actionID {
				dataObj.Actions[i].Status = "Done"
				awardID = item.AwardID
				allPeopleCount = item.PeopleCount
				backMoney = item.BackMoney
				break
			} else {
				c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "当前奖品状态不正确", "data": nil})
				return
			}
		} else if i+1 == len(dataObj.Actions) {
			c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "奖品已经抽完", "data": nil})
			return
		}
	}
	allPeopleCount = allPeopleCount - staffCount - leaderCount;
	if allPeopleCount < 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "员工人数与领导人数之和 超过设置总人数", "data": nil})
		return
	}
	var returnData []models.DrawedRecord

	for i, item := range dataObj.Users {
		if leaderCount <= 0 {
			break
		}
		if item.Level == 2 && item.IsDrawed == false {
			record := models.DrawedRecord{
				AwardID:        awardID,
				AwardCount:     1,
				Drawer:         drawer,
				DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
				LuckyUserID:    item.ID,
				LuckyUserLevel: item.Level,
				LuckyUserName:  item.Name,
				Memo:           memo,
			}
			dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
			returnData = append(returnData, record)
			dataObj.Users[i].IsDrawed = true
			leaderCount--
		}
	}
	if leaderCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "未中奖领导人数不够", "data": nil})
		return
	}

	for i, item := range dataObj.Users {
		if staffCount <= 0 {
			break
		}
		if item.Level == 1 && item.IsDrawed == false {
			record := models.DrawedRecord{
				AwardID:        awardID,
				AwardCount:     1,
				Drawer:         drawer,
				DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
				LuckyUserID:    item.ID,
				LuckyUserLevel: item.Level,
				LuckyUserName:  item.Name,
				Memo:           memo,
			}
			dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
			returnData = append(returnData, record)
			dataObj.Users[i].IsDrawed = true
			staffCount--
		}
	}
	if staffCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "未中奖员工人数不够", "data": nil})
		return
	}

	for i, item := range dataObj.Users {
		if allPeopleCount <= 0 {
			break
		}
		if item.IsDrawed == false {
			record := models.DrawedRecord{
				AwardID:        awardID,
				AwardCount:     1,
				Drawer:         drawer,
				DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
				LuckyUserID:    item.ID,
				LuckyUserLevel: item.Level,
				LuckyUserName:  item.Name,
				Memo:           memo,
			}
			dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
			returnData = append(returnData, record)
			dataObj.Users[i].IsDrawed = true
			allPeopleCount--
		}
	}
	if allPeopleCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "未中奖人数不够", "data": nil})
		return
	}

	for _, item := range returnData {
		if item.LuckyUserLevel == 2 {
			backMoneyRecord := models.BackMoneyRecord{
				BackTime: time.Now().Format("2006-01-02 15:04:05"),
				Memo:     "奖品返奖",
				Money:    backMoney,
				UserID:   item.LuckyUserID,
				UserName: item.LuckyUserName,
			}
			dataObj.BackMoneyRecords = append(dataObj.BackMoneyRecords, backMoneyRecord)
		}
	}

	setDataErr:=common.SetData(dataObj)
	if setDataErr==nil{
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": returnData})
		return
	}else {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": "SetData error", "data": returnData})
		return
	}
}

func ExDraw(c *gin.Context) {
	drawer := c.Query("drawer")
	memo := c.Query("memo")

	awardID, err := strconv.Atoi(c.Query("awardID"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "奖品ID错误", "data": nil})
		return
	}
	awardCount, err := strconv.Atoi(c.Query("awardCount"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "奖品数量错误", "data": nil})
		return
	}
	backMoney, err := strconv.Atoi(c.Query("backMoney"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "返奖金额错误", "data": nil})
		return
	}

	mixPeopleCount, err := strconv.Atoi(c.Query("mixPeopleCount"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "混合人数错误", "data": nil})
		return
	}
	leaderCount, err := strconv.Atoi(c.Query("leaderCount"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "领导人数人数错误", "data": nil})
		return
	}
	staffCount, err := strconv.Atoi(c.Query("staffCount"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "员工人数错误", "data": nil})
		return
	}
	dataObj := common.GetData()
	for i, item := range dataObj.Awards {
		if item.ID == awardID {
			break
			if (i == len(dataObj.Awards)) {
				c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "奖品类型不在字典内", "data": nil})
				return
			}
		}
	}

	var returnData []models.DrawedRecord
	for i, item := range dataObj.Users {
		if mixPeopleCount <= 0 {
			break
		}
		if item.IsDrawed == false {
			record := models.DrawedRecord{
				AwardID:        awardID,
				AwardCount:     awardCount,
				Drawer:         drawer,
				DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
				LuckyUserID:    item.ID,
				LuckyUserLevel: item.Level,
				LuckyUserName:  item.Name,
				Memo:           memo,
			}
			dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
			returnData = append(returnData, record)
			dataObj.Users[i].IsDrawed = true
			mixPeopleCount--
		}
	}
	if mixPeopleCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "未中奖人数不够", "data": nil})
		return
	}

	for i, item := range dataObj.Users {
		if leaderCount <= 0 {
			break
		}
		if item.IsDrawed == false && item.Level == 2 {
			record := models.DrawedRecord{
				AwardID:        awardID,
				AwardCount:     awardCount,
				Drawer:         drawer,
				DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
				LuckyUserID:    item.ID,
				LuckyUserLevel: item.Level,
				LuckyUserName:  item.Name,
				Memo:           memo,
			}
			dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
			returnData = append(returnData, record)
			dataObj.Users[i].IsDrawed = true
			leaderCount--
		}
	}
	if leaderCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "领导未中奖人数不够", "data": nil})
		return
	}

	for i, item := range dataObj.Users {
		if staffCount <= 0 {
			break
		}
		if item.IsDrawed == false && item.Level == 1 {
			record := models.DrawedRecord{
				AwardID:        awardID,
				AwardCount:     awardCount,
				Drawer:         drawer,
				DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
				LuckyUserID:    item.ID,
				LuckyUserLevel: item.Level,
				LuckyUserName:  item.Name,
				Memo:           memo,
			}
			dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
			returnData = append(returnData, record)
			dataObj.Users[i].IsDrawed = true
			staffCount--
		}
	}
	if staffCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "员工未中奖人数不够", "data": nil})
		return
	}
	for _, item := range returnData {
		if item.LuckyUserLevel == 2 {
			backMoneyRecord := models.BackMoneyRecord{
				BackTime: time.Now().Format("2006-01-02 15:04:05"),
				Memo:     "奖品返奖",
				Money:    backMoney,
				UserID:   item.LuckyUserID,
				UserName: item.LuckyUserName,
			}
			dataObj.BackMoneyRecords = append(dataObj.BackMoneyRecords, backMoneyRecord)
		}
	}
	setDataErr:=common.SetData(dataObj)
	if setDataErr==nil{
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": returnData})
		return
	}else {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": "SetData error", "data": returnData})
		return
	}
}

func PoolDraw(c *gin.Context) {
	drawer := c.Query("drawer")
	memo := c.Query("memo")
	backRatio, err := strconv.Atoi(c.Query("backRatio"))
	if err != nil && backRatio >= 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "返奖系数错误", "data": nil})
		return
	}
	awardCount, err := strconv.Atoi(c.Query("awardCount"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "奖品数量错误", "data": nil})
		return
	}

	mixPeopleCount, err := strconv.Atoi(c.Query("mixPeopleCount"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "混合人数错误", "data": nil})
		return
	}
	leaderCount, err := strconv.Atoi(c.Query("leaderCount"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "领导人数人数错误", "data": nil})
		return
	}
	staffCount, err := strconv.Atoi(c.Query("staffCount"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "员工人数错误", "data": nil})
		return
	}
	dataObj := common.GetData()
	var returnData []models.DrawedRecord
	for i, item := range dataObj.Users {
		if mixPeopleCount <= 0 {
			break
		}
		if item.IsDrawed == false {
			record := models.DrawedRecord{
				AwardID:        0, //返奖池奖品类型
				AwardCount:     awardCount,
				Drawer:         drawer,
				DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
				LuckyUserID:    item.ID,
				LuckyUserLevel: item.Level,
				LuckyUserName:  item.Name,
				Memo:           memo,
			}
			dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
			returnData = append(returnData, record)
			dataObj.Users[i].IsDrawed = true
			mixPeopleCount--
		}
	}
	if mixPeopleCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "未中奖人数不够", "data": nil})
		return
	}

	for i, item := range dataObj.Users {
		if leaderCount <= 0 {
			break
		}
		if item.IsDrawed == false && item.Level == 2 {
			record := models.DrawedRecord{
				AwardID:        0, //返奖池奖品类型
				AwardCount:     awardCount,
				Drawer:         drawer,
				DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
				LuckyUserID:    item.ID,
				LuckyUserLevel: item.Level,
				LuckyUserName:  item.Name,
				Memo:           memo,
			}
			dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
			returnData = append(returnData, record)
			dataObj.Users[i].IsDrawed = true
			leaderCount--
		}
	}
	if leaderCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "领导未中奖人数不够", "data": nil})
		return
	}

	for i, item := range dataObj.Users {
		if staffCount <= 0 {
			break
		}
		if item.IsDrawed == false && item.Level == 1 {
			record := models.DrawedRecord{
				AwardID:        0, //返奖池奖品类型
				AwardCount:     awardCount,
				Drawer:         drawer,
				DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
				LuckyUserID:    item.ID,
				LuckyUserLevel: item.Level,
				LuckyUserName:  item.Name,
				Memo:           memo,
			}
			dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
			returnData = append(returnData, record)
			dataObj.Users[i].IsDrawed = true
			staffCount--
		}
	}
	if staffCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "员工未中奖人数不够", "data": nil})
		return
	}
	if backRatio > 0 {
		for _, item := range returnData {
			if item.LuckyUserLevel == 2 {
				backMoneyRecord := models.BackMoneyRecord{
					BackTime: time.Now().Format("2006-01-02 15:04:05"),
					Memo:     "奖品返奖",
					Money:    awardCount * backRatio,
					UserID:   item.LuckyUserID,
					UserName: item.LuckyUserName,
				}
				dataObj.BackMoneyRecords = append(dataObj.BackMoneyRecords, backMoneyRecord)
			}
		}
	}
	setDataErr := common.SetData(dataObj)
	if setDataErr==nil{
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": returnData})
		return
	}else {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": "SetData error", "data": returnData})
		return
	}
}
