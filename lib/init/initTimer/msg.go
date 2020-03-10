package initTimer

const (
	//定义timer的宏
	ACTION_HEAT_BEATE  = 0 //心跳检测
	ACTION_PARTY_START_TIMER = 1 //饭局开始的时间
)

type Msg struct {
	Action int16
	Data   interface{}
}

type GSTimer struct {
	TimerId  int64  //每个timer都有自己的唯一id
	Expired  int64 // expired time, timestamp(ms)
	Msg      Msg
}
