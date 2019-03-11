package api

import (
	"errors"
	"math/rand"
	"strconv"

	"../models"
)

///////
func GetLuckyUserID(users []models.User, level int) (luckyID int, err error) {
	len := len(users)
	for i, item := range users {
		if item.IsDrawed == false && (level == 0 || item.Level == level) {
			break
		} else if i == len-1 {
			return 0, errors.New("对应level:" + strconv.Itoa(level) + "未中奖人数不足")
		}
	}
	luckyID = rand.Intn(len-1) + 1
	if (users[luckyID].IsDrawed == false) && (level == 0 || users[luckyID].Level == level) {
		return luckyID, nil
	}
	return GetLuckyUserID(users, level)
}
