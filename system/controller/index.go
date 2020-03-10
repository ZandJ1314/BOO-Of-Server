package controller

import (
	"boo/lib/tools"
	"boo/system/model"
	"boo/system/servers"
	"fmt"
	"github.com/gin-gonic/gin"
)

//登录成功后返回给前端的数据
type Circle struct {
	Title 			string   `json:"title"`//圈子主题
	Members			[]*Users `json:"members"`//圈子成员
	AppointNums		int16    `json:"appoint_nums"`//约局次数
	Days 			int32    `json:"days"`//从第一次约到现在多少天了
	MeetInfo		[]*Meet  `json:"meet_info"`//圈子局子
}

//局
type Meet struct {
	CreateTime 		int64  `json:"create_time"`//创建时间
	Topic			string `json:"topic"`//局主题
	Site			string `json:"site"`//地点
	Status			int8   `json:"status"`//状态{0:成功，1:等待，2:失败}
}

type Users struct {
	NickName	string `json:"nick_name"`
	AvatarUrl	string `json:"avatar_url"`
}

type IndexResponseInfo struct {
	UserId 		int64     `json:"user_id"`
	CircleInfo  []*Circle `json:"circle_info"`
}


func Index(c *gin.Context) {
	userId := tools.GetInt64(c.Query("userId"))
	//通过userid从数据库取用户数据
	//返回所有玩家的圈子数量
	err,circleInfoList := model.GetSelectedCircleByUserID(userId)
	if err != nil {
		servers.ReportFormat(c,false,fmt.Sprintf("不存在该玩家:%v",err),gin.H{})
		return
	}
	if len(circleInfoList) == 0 {
		//说明该玩家没有圈子返回虚拟圈子
		servers.ReportFormat(c,true,"获取成功",gin.H{"data":nil})
		return
	}
	circleList := make([]*Circle,0)
	for _,conf := range circleInfoList {
		if conf == nil {
			continue
		}
		err,circle := model.GetUserCircleInfo(conf.CircleID)
		if err != nil {
			continue
		}
		err,appoints := model.GetUserAppointInfo(conf.CircleID)
		if err != nil {
			continue
		}
		meetInfo := make([]*Meet,0)
		for _,appointInfo := range appoints {
			meet := &Meet{
				CreateTime: appointInfo.CreateTime,
				Topic:      appointInfo.Topic,
				Site:       appointInfo.AppointPlace,
				Status:     appointInfo.Status,
			}
			meetInfo = append(meetInfo,meet)
		}
		err,selectedCircle := model.GetSelectedUserByCircleId(conf.CircleID)
		if err != nil {
			continue
		}
		userIdList := model.GetSelectedUserListByCircleList(selectedCircle)
		circleConf := &Circle{
			Title:       circle.CircleName,
			Members:     getUsersByUserList(userIdList),
			AppointNums: circle.AppointNums,
			Days:        tools.GetDaysToNowByTime(circle.CreateTime),
			MeetInfo:    meetInfo,
		}
		circleList = append(circleList,circleConf)
	}
	resp := &IndexResponseInfo{
		UserId:userId,
		CircleInfo:circleList,
	}
	servers.ReportFormat(c,true,"获取成功",gin.H{"data":resp})
}

func getUsersByUserList(users []int64) []*Users{
	userInfo := make([]*Users,0)
	for _,userID := range users {
		err,user := model.FindUserByUserID(userID)
		if err != nil {
			continue
		}
		userConf := &Users{
			NickName:  user.Nickname,
			AvatarUrl: user.Avatar,
		}
		userInfo = append(userInfo,userConf)
	}
	return userInfo
}
