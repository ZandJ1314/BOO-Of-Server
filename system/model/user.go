package model

import (
	"boo/lib/init/initSql"
	"boo/lib/tools"
	"boo/system/servers"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserId			int64 `json:"user_id"`
	Avatar        	string `json:"avatar"`
	Birthday      	int    `json:"birthday"`
	Gender        	int    `json:"gender"`
	LastLoginIp   	string `json:"last_login_ip"`
	LastLoginTime 	int64  `json:"last_login_time"`
	Mobile        	string `json:"mobile"`
	Nickname      	string `json:"nickname"`
	Password      	string `json:"password" gorm:"default:'-'"`
	RegisterIp    	string `json:"register_ip"`
	RegisterTime  	int64  `json:"register_time"`
	UserLevelId   	int    `json:"user_level_id"`
	Username      	string `json:"username"`
	WeiXinOpenid  	string `json:"wei_xin_openid"`
}

//注册接口model方法
func (u *User) Insert(userInfo *servers.WXUserInfo,ip string) (err error, userInter *User) {
	u.UserId = tools.GetUUid()
	u.Avatar = userInfo.AvatarUrl
	u.Birthday = 1
	u.LastLoginIp = ip
	err = initSql.DEFAULTDB.Create(&u).Error
	return err,u
}


func FindDataByOpenID(openId string) (err error ,userInfo *User) {
	var user User
	err = initSql.DEFAULTDB.Where("wei_xin_openid = ?",openId).First(&user).Error
	if err != nil {
		return err,nil
	}
	return nil,&user
}

func FindUserByUserID(userId int64) (err error,userInfo *User) {
	var user User
	err = initSql.DEFAULTDB.Where("user_id = ?",userId).First(&user).Error
	if err != nil {
		return err,nil
	}
	return nil,&user
}
