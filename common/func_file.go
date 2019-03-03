package common

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"io/ioutil"

	"../models"
)

func GetData() models.Data {
	dataName := "data.json"
	dataJsonByte, err := ioutil.ReadFile(dataName)
	if err != nil {
		fmt.Print(err)
	}
	dataJsonStr := string(dataJsonByte)
	var dataObj models.Data
	json.Unmarshal([]byte(dataJsonStr), &dataObj)
	return dataObj
}

var (
	mutex sync.Mutex
)

func SetData(data models.Data) error {
	mutex.Lock()
	{
		if fileLock {
			return errors.New("data file has locked")
		} else {
			fileLock = false
		}
		preData := GetData()
		if preData.Version != data.Version {
			return errors.New("data version error")
		} else {
			data.Version = GetUUID()
		}
		data, err := CountData(data)
		if err != nil {
			return errors.New("CountData error")
		}
		dataJsonByte, err := json.Marshal(data)
		if err != nil {
			return errors.New("json.Marshal error")
		}
		dataJsonStr := string(dataJsonByte)
		ioutil.WriteFile("data.json", []byte(dataJsonStr), 0644)
		fileLock = true
	}
	mutex.Unlock()
	return nil
}

var fileLock bool = false

func CountData(dataObj models.Data) (models.Data, error) {
	//PoolMoney
	dataObj.Count.PoolMoney = 0
	for _, item := range dataObj.BackMoneyRecords {
		dataObj.Count.PoolMoney += item.Money
	}
	for _, item := range dataObj.DrawedRecords {
		if item.AwardID == 0 {
			dataObj.Count.PoolMoney -= item.AwardCount
		}
	}
	if dataObj.Count.PoolMoney < 0 {
		return dataObj, errors.New("PoolMoney 小于0")
	}
	//AllPeopleCount
	dataObj.Count.AllPeopleCount = len(dataObj.Users)

	dataObj.Count.AllLuckyCount = 0
	for _, item := range dataObj.Users {
		if item.IsDrawed {
			dataObj.Count.AllLuckyCount += 1
		}
	}
	//LuckyLeaderCount
	dataObj.Count.LuckyLeaderCount = 0
	for _, item := range dataObj.Users {
		if item.IsDrawed && item.Level == 2 {
			dataObj.Count.LuckyLeaderCount += 1
		}
	}

	//NoLuckyLeaderCount
	dataObj.Count.NoLuckyLeaderCount = 0
	for _, item := range dataObj.Users {
		if item.IsDrawed == false && item.Level == 2 {
			dataObj.Count.NoLuckyLeaderCount += 1
		}
	}

	//LuckyStaffCount
	dataObj.Count.LuckyStaffCount = 0
	for _, item := range dataObj.Users {
		if item.IsDrawed == true && item.Level == 1 {
			dataObj.Count.LuckyStaffCount += 1
		}
	}

	//NoLuckyStaffCount
	dataObj.Count.NoLuckyStaffCount = 0
	for _, item := range dataObj.Users {
		if item.IsDrawed == false && item.Level == 1 {
			dataObj.Count.NoLuckyStaffCount += 1
		}
	}
	return dataObj, nil
}

func Log(log string) {
	log = time.Now().Format("2006-01-02 15:04:05") + "|" + log
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	f.WriteString(log)
}

func GetUUID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
	return hex.EncodeToString(h.Sum(nil))
}
