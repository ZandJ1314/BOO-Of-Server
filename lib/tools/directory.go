package tools

import (
	"os"
	"strconv"
	"time"
)

var baseUid = int64(1)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetUUid() int64 {
	if baseUid == 1 {
		return 1
	}
	baseUid++
	return baseUid
}


func GetInt64(str string) int64 {
	number, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return number
	}
	return 0
}

func GetDaysToNowByTime(t int64) int32 {
	now := time.Now().Unix() //获取当前时间的时间戳，秒为单位
	tValue := now-t
	//1d=24h*60m*60s
	days := int32(tValue/(24*60*60))
	return days
}