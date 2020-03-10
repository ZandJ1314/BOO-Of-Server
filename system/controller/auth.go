package controller

import (
	"boo/system/model"
	"boo/system/servers"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type AuthLoginBody struct {
	Code 		string `json:"code"`
	UserInfo 	servers.ResUserInfo `json:"userInfo"`
}

func LoginByWein(c *gin.Context) {
	var auth AuthLoginBody
	body,_ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &auth)
	if err != nil {
		servers.ReportFormat(c,false,fmt.Sprintf("登录失败:%v",err),gin.H{})
		return
	}
	clientIp := c.ClientIP()
	userInfo,loginErr := servers.Login(auth.Code,&auth.UserInfo)
	if loginErr != nil {
		servers.ReportFormat(c,false,fmt.Sprintf("登录失败:%v",err),gin.H{})
		return
	}
	//先进行查询，如果查询不到就进行注册
	var user model.User
	initErr,userInit := model.FindDataByOpenID(userInfo.OpenID)
	if initErr != nil || userInfo == nil {
		//需要插入新的数据
		insertErr,newUser := user.Insert(userInfo,clientIp)
		if insertErr != nil {
			servers.ReportFormat(c,false,fmt.Sprintf("登录失败:%v",err),gin.H{})
			return
		}
		user = *newUser
	}else {
		user = *userInit
	}
	//登录成功了，需要返回数据
	servers.ReportFormat(c,true,fmt.Sprintf("登录成功:%v",user.ID),gin.H{"userId":user.UserId})
}
