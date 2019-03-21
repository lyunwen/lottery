package models

import (
	"errors"
	"math/rand"
	"time"
)

type Data struct {
	Version          string
	Description      string
	Count            Count
	Awards           []Award
	Actions          []DrawedAction
	Users            []User
	DrawedRecords    []DrawedRecord
	BackMoneyRecords []BackMoneyRecord
}

func (data *Data) NDraw(drawKey string, actionID int, staffCount int, leaderCount int, drawer string, memo string) (*Data, error) {
	var currentAction *DrawedAction
	actionLen := len(data.Actions)
	for i, _ := range data.Actions {
		if data.Actions[i].Status == "ToDo" {
			if data.Actions[i].ID == actionID {
				data.Actions[i].Status = "Done"
				currentAction = &data.Actions[i]
			} else {
				return data, errors.New("当前步骤奖品状态不正确")
			}
		} else if i+1 == actionLen {
			return data, errors.New("奖品已经抽完")
		}
	}
	lessCount := currentAction.PeopleCount - staffCount - leaderCount
	if leaderCount < 0 {
		return data, errors.New("员工人数与领导人数之和 超过设置总人数")
	}
	var indexList []int
	var indexErr error
	if leaderCount > 0 {
		thisIndex, indexErr := getLuckyUserIndex(data.Users, 2)
		if indexErr != nil {
			return data, indexErr
		}
		indexList = append(indexList, thisIndex)
	}
	if staffCount > 0 {
		thisIndex, indexErr := getLuckyUserIndex(data.Users, 1)
		if indexErr != nil {
			return data, indexErr
		}
		indexList = append(indexList, thisIndex)
	}
	if lessCount > 0 {
		thisIndex, indexErr := getLuckyUserIndex(data.Users, 0)
		if indexErr != nil {
			return data, indexErr
		}
		indexList = append(indexList, thisIndex)
	}

	for _, item := range indexList {
		data.Users[item].IsDrawed = true
		record := DrawedRecord{
			DrawKey:        drawKey,
			AwardID:        currentAction.AwardID,
			AwardCount:     1,
			Drawer:         drawer,
			DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:    data.Users[item].ID,
			LuckyUserLevel: data.Users[item].Level,
			LuckyUserName:  data.Users[item].Name,
			Memo:           memo,
		}
		data.DrawedRecords = append(data.DrawedRecords, record)
		if data.Users[item].Level == 2 {
			backMoneyRecord := BackMoneyRecord{
				BackTime: time.Now().Format("2006-01-02 15:04:05"),
				Memo:     "奖品返奖",
				Money:    currentAction.BackMoney,
				UserID:   data.Users[item].ID,
				UserName: data.Users[item].Name,
			}
			data.BackMoneyRecords = append(data.BackMoneyRecords, backMoneyRecord)
		}
	}
	return data, indexErr
}

func (data *Data) ExDraw(drawKey string, staffCount int, leaderCount int, mixCount int, awardID int, awardCount int, backMoney int, drawer string, memo string) (*Data, error) {
	var indexList []int
	var indexErr error
	if leaderCount > 0 {
		thisIndex, indexErr := getLuckyUserIndex(data.Users, 2)
		if indexErr != nil {
			return data, indexErr
		}
		indexList = append(indexList, thisIndex)
	}
	if staffCount > 0 {
		thisIndex, indexErr := getLuckyUserIndex(data.Users, 1)
		if indexErr != nil {
			return data, indexErr
		}
		indexList = append(indexList, thisIndex)
	}
	if mixCount > 0 {
		thisIndex, indexErr := getLuckyUserIndex(data.Users, 0)
		if indexErr != nil {
			return data, indexErr
		}
		indexList = append(indexList, thisIndex)
	}

	for _, item := range indexList {
		data.Users[item].IsDrawed = true
		record := DrawedRecord{
			DrawKey:        drawKey,
			AwardID:        awardID,
			AwardCount:     awardCount,
			Drawer:         drawer,
			DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:    data.Users[item].ID,
			LuckyUserLevel: data.Users[item].Level,
			LuckyUserName:  data.Users[item].Name,
			Memo:           memo,
		}
		data.DrawedRecords = append(data.DrawedRecords, record)
		if data.Users[item].Level == 2 {
			backMoneyRecord := BackMoneyRecord{
				BackTime: time.Now().Format("2006-01-02 15:04:05"),
				Memo:     "奖品返奖",
				Money:    backMoney,
				UserID:   data.Users[item].ID,
				UserName: data.Users[item].Name,
			}
			data.BackMoneyRecords = append(data.BackMoneyRecords, backMoneyRecord)
		}
	}
	// for i, item := range data.Awards {
	// 	if item.ID == awardID {
	// 		break
	// 	} else if i+1 == len(data.Awards) {
	// 		return data, errors.New("奖品类型:" + strconv.Itoa(awardID) + "不正确")
	// 	}
	// }
	return data, indexErr
}

func (data *Data) PoolDraw(drawKey string, staffCount int, leaderCount int, mixCount int, awardCount int, backRatio int, drawer string, memo string) (*Data, error) {
	var indexList []int
	var indexErr error
	if leaderCount > 0 {
		thisIndex, indexErr := getLuckyUserIndex(data.Users, 2)
		if indexErr != nil {
			return data, indexErr
		}
		indexList = append(indexList, thisIndex)
	}
	if staffCount > 0 {
		thisIndex, indexErr := getLuckyUserIndex(data.Users, 1)
		if indexErr != nil {
			return data, indexErr
		}
		indexList = append(indexList, thisIndex)
	}
	if mixCount > 0 {
		thisIndex, indexErr := getLuckyUserIndex(data.Users, 0)
		if indexErr != nil {
			return data, indexErr
		}
		indexList = append(indexList, thisIndex)
	}

	for _, item := range indexList {
		data.Users[item].IsDrawed = true
		record := DrawedRecord{
			DrawKey:        drawKey,
			AwardID:        0, //返奖池奖品类型
			AwardCount:     awardCount,
			Drawer:         drawer,
			DrawTime:       time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:    data.Users[item].ID,
			LuckyUserLevel: data.Users[item].Level,
			LuckyUserName:  data.Users[item].Name,
			Memo:           memo,
		}
		data.DrawedRecords = append(data.DrawedRecords, record)
		if data.Users[item].Level == 2 && backRatio > 0 {
			backMoneyRecord := BackMoneyRecord{
				BackTime: time.Now().Format("2006-01-02 15:04:05"),
				Memo:     "现金池返奖",
				Money:    awardCount * backRatio,
				UserID:   data.Users[item].ID,
				UserName: data.Users[item].Name,
			}
			data.BackMoneyRecords = append(data.BackMoneyRecords, backMoneyRecord)
		}
	}
	return data, indexErr
}

func getLuckyUserIndex(users []User, level int) (int, error) {
	lenUser := len(users)
	if lenUser == 0 {
		return 0, errors.New("未抽奖人数为空")
	}
	luckyIndex := rand.Intn(lenUser)
	if users[luckyIndex].IsDrawed == false && (level < 1 || level == users[luckyIndex].Level) {
		return luckyIndex, nil
	} else {
		users = append(users[:luckyIndex], users[luckyIndex+1:]...)
		return getLuckyUserIndex(users, level)
	}
}
