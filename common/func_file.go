package common

import (
	"os"
	"time"
)

//日志方法
func Log(log string) error {
	log = "\r\n" + time.Now().Format("2006-01-02 15:04:05") + "|" + log
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	_, err = f.WriteString(log)
	return err
}
