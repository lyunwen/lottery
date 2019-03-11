package common

import (
	"os"
	"time"
)

//日志方法
func Log(log string) {
	log = time.Now().Format("2006-01-02 15:04:05") + "|" + log
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	f.WriteString(log)
}
