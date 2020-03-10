package initTimer

import (
	"boo/lib/tools"
	"fmt"
	"time"
)

//这里是timer的初始化

var (
	TIMER_MAP   map[int64]*GSTimer
)


func InitTimer() {
	now := time.Now().Unix()
	if TIMER_MAP == nil {
		TIMER_MAP = make(map[int64]*GSTimer)
	}else {
		for timerId,ts := range TIMER_MAP {
			CreateTimer(timerId,ts.Expired-now,ts.Msg)
		}
	}
	//分配一个协程去跑timer
	go starTimerWork()
}

func starTimerWork(){
	//新建一个定时器，一般2秒检测一次
	tick := time.NewTicker(2*time.Second)
	for {
		select {
		//在此处等待channel中的信号，2秒循环检测一次
			case <-tick.C:
				//处理回调
				for _,ts := range TIMER_MAP {
					timerHandler(ts)
				}
		}
	}
}

func timerHandler(ts *GSTimer) {
	switch ts.Msg.Action {
	case ACTION_HEAT_BEATE:
		fmt.Printf("heat beate %v",ts.TimerId)
	default:
		fmt.Printf("timer test err %v",ts.TimerId)
	}
}

//创建一个timer，delay的单位为秒
func CreateTimer(oldTimerId int64,delay int64,msg Msg) {
	if delay <= 0 {
		delay = 0
	}
	expireTime := time.Now().Unix()+delay
	timerId := tools.GetUUid()
	gsTime := &GSTimer{}
	if oldTimerId != 0 {
		timerId = oldTimerId
		gsTime = TIMER_MAP[timerId]
	}else {
		gsTime = &GSTimer{
			TimerId: timerId,
			Expired: expireTime,
			Msg:     msg,
		}
	}
	TIMER_MAP[timerId] = gsTime
}
