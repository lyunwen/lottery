package models

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"

	"../common"
)

func (data *Data) GetData() (*Data, error) {
	dataName := "data.json"
	dataJsonByte, err := ioutil.ReadFile(dataName)
	if err != nil {
		return data, err
	}
	dataJsonStr := string(dataJsonByte)
	err = json.Unmarshal([]byte(dataJsonStr), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (data *Data) SetData() error {
	mutex.Lock()
	if fileLock {
		return errors.New("data.json file has locked")
	} else {
		fileLock = true
	}
	err := data.SetDataCore()
	fileLock = false
	mutex.Unlock()
	return err
}

var (
	mutex sync.Mutex
)

//非线程安全
func (data *Data) SetDataCore() error {
	preData, err := new(Data).GetData()
	if err != nil {
		return err
	}
	if preData.Version != data.Version {
		return errors.New("data version error")
	} else {
		data.Version = common.GetUUID()
	}
	data, err = data.CountData()
	if err != nil {
		return err
	}
	dataJsonByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	dataJsonStr := string(dataJsonByte)
	err = ioutil.WriteFile("data.json", []byte(dataJsonStr), 0644)
	if err != nil {
		return err
	}
	return nil
}

var fileLock bool = false

func (data *Data) CountData() (*Data, error) {
	//PoolMoney
	data.Count.PoolMoney = 0
	for _, item := range data.BackMoneyRecords {
		data.Count.PoolMoney += item.Money
	}
	for _, item := range data.DrawedRecords {
		if item.AwardID == 0 {
			data.Count.PoolMoney -= item.AwardCount
		}
	}
	if data.Count.PoolMoney < 0 {
		return data, errors.New("pool money less than 0")
	}
	//AllPeopleCount
	data.Count.AllPeopleCount = len(data.Users)

	data.Count.AllLuckyCount = 0
	for _, item := range data.Users {
		if item.IsDrawed {
			data.Count.AllLuckyCount += 1
		}
	}
	//LuckyLeaderCount
	data.Count.LuckyLeaderCount = 0
	for _, item := range data.Users {
		if item.IsDrawed && item.Level == 2 {
			data.Count.LuckyLeaderCount += 1
		}
	}

	//NoLuckyLeaderCount
	data.Count.NoLuckyLeaderCount = 0
	for _, item := range data.Users {
		if item.IsDrawed == false && item.Level == 2 {
			data.Count.NoLuckyLeaderCount += 1
		}
	}

	//LuckyStaffCount
	data.Count.LuckyStaffCount = 0
	for _, item := range data.Users {
		if item.IsDrawed == true && item.Level == 1 {
			data.Count.LuckyStaffCount += 1
		}
	}

	//NoLuckyStaffCount
	data.Count.NoLuckyStaffCount = 0
	for _, item := range data.Users {
		if item.IsDrawed == false && item.Level == 1 {
			data.Count.NoLuckyStaffCount += 1
		}
	}
	return data, nil
}
