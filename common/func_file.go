package common

import (
	"os"
	"strings"
	"time"
)

//日志方法
func Log(log string) error {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	//_, err = f.WriteString(log)

	fd_time := time.Now().Format("2006-01-02 15:04:05")
	log = strings.Join([]string{"======", fd_time, "=====", log, "\n"}, "")
	buf := []byte(log)
	f.Write(buf)
	f.Close()

	return err
}
