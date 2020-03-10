package conf

const (
	WeiXinAppId_= "zhangfan" //小程序appid
	WeiXinSecret = "zhangfan good" //小程序密钥
)

type Config struct {
	MysqlAdmin   MysqlAdmin   `json:"mysqlAdmin"`
	Qiniu        Qiniu        `json:"qiniu"`
	CasbinConfig CasbinConfig `json:"casbinConfig"`
	RedisAdmin   RedisAdmin   `json:"redisAdmin"`
	System       System       `json:"system"`
	JWT          JWT          `json:"jwt"`
	Weixin		 Weixin		  `json:"weixin"`
}

type System struct { // 系统配置
	UseMultipoint bool   `json:"useMultipoint"`
	Env           string `json:"env"`
	Addr          int    `json:"addr"`
}

type Weixin struct { //微信的配置
	appid		string `json:"appid"`
	secret		string `json:"secret"`
}


type JWT struct { // jwt签名
	SigningKey string `json:"signingKey"`
}

type CasbinConfig struct { //casbin配置
	ModelPath string `json:"modelPath"` // casbin model地址配置
}

type MysqlAdmin struct { // mysql admin 数据库配置
	Username string `json:"username"`
	Password string `json:"password"`
	Path     string `json:"path"`
	Dbname   string `json:"dbname"`
	Config   string `json:"config"`
}

type RedisAdmin struct { // Redis admin 数据库配置
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}
type Qiniu struct { // 七牛 密钥配置
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

var GinAdminconfig Config

func (wei *Weixin)GetWeiXinAppid() string {
	return GinAdminconfig.Weixin.appid
}

func (wei *Weixin)SetWeiXinAppid(appid string) {
	GinAdminconfig.Weixin.appid = appid
}

func (wei *Weixin)GetWeiXinSecret() string {
	return GinAdminconfig.Weixin.secret
}

func (wei *Weixin)SetWeiXinSecret(secret string) {
	GinAdminconfig.Weixin.secret = secret
}

func (q *MysqlAdmin) SetMysqlPwd(pwd string) {
	GinAdminconfig.MysqlAdmin.Password = pwd
}

func (q *MysqlAdmin) SetMysqlUserName(name string) {
	GinAdminconfig.MysqlAdmin.Username = name
}

func (q *MysqlAdmin) SetMysqlPath(path string) {
	GinAdminconfig.MysqlAdmin.Path = path
}

func (q *MysqlAdmin) SetMysqlDbname(Dbname string) {
	GinAdminconfig.MysqlAdmin.Dbname = Dbname
}



func init() {
	//对weixin进行初始化
	GinAdminconfig.Weixin.SetWeiXinAppid(WeiXinAppId_)
	GinAdminconfig.Weixin.SetWeiXinSecret(WeiXinSecret)
	//对数据库进行初始化
}
