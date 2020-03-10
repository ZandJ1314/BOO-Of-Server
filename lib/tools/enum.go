package tools

//约局的枚举
type ApponitEnum struct {
	Success 	int8
	Failed		int8
	Wait 		int8
}

type PunishOfLateEnum struct {
	Drink		int8
	Telephone	int8
}

type FailedAppointEnum struct {
	Money		int8
	TelePhone	int8
}

var (
	EnumAppointInfo = ApponitEnum{Success:0,Failed:1,Wait:2}
	EnumPunishOfLateEnum = PunishOfLateEnum{Drink:1,Telephone:2}
	EnumFailedAppointEnum = FailedAppointEnum{Money:1,TelePhone:2}
)
