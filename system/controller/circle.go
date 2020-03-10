package controller

import (
	"boo/lib/init/initSql"
	"boo/lib/tools"
	"boo/system/model"
	"boo/system/servers"
	"github.com/gin-gonic/gin"
)

//创建圈子和加入圈子和退出圈子三种：成功或者失败

type CircleInfo struct {
	CircleID		int64 	`json:"circle_id"` //圈子id
	CircleName 		string 	`json:"circle_name"`//圈子名称
	CreateUserId 	int64	`json:"create_user_id"`//创建人uid
	CreateUserName	string 	`json:"create_user_name"` //创建人name
}

//创建圈子
func CreateCircle(c *gin.Context) {
	topic := c.Query("circle_name")
	userId := tools.GetInt64(c.Query("userId"))
	db := initSql.DEFAULTDB.Begin()
	if db == nil {
		servers.ReportFormat(c,false,"数据库打开错误，创建圈子失败",gin.H{})
		return
	}
	defer db.Rollback()
	err,user := model.FindUserByUserID(userId)
	if err != nil {
		servers.ReportFormat(c,false,"用户不存在，创建圈子失败",gin.H{})
		return
	}
	circle := &CircleInfo{
		CircleID:       tools.GetUUid(),
		CircleName:     topic,
		CreateUserId:   userId,
		CreateUserName: user.Nickname,
	}
	err = model.InsertCircle(circle,db)
	if err != nil {
		servers.ReportFormat(c,false,"插入数据失败，创建圈子失败",gin.H{})
		return
	}
	err = model.InsertSelectedCircle(userId,circle.CircleID,db)
	if err != nil {
		servers.ReportFormat(c,false,"插入数据失败，创建圈子失败",gin.H{})
		return
	}
	//进行事务提交
	db.Commit()
	servers.ReportFormat(c,true,"创建圈子成功",gin.H{})
}

//func QueryCircleTopic(c *gin.Context) {
//	//同一个用户
//	topic := c.Query("circle_name")
//	userId := tools.GetInt64(c.Query("userId"))
//}

//加入圈子
func JoinCircle(c *gin.Context) {
	userId := tools.GetInt64(c.Query("userId"))
	circleId := tools.GetInt64(c.Query("circleId"))
	db := initSql.DEFAULTDB.Begin()
	if db == nil {
		servers.ReportFormat(c,false,"数据库打开错误，加入圈子失败",gin.H{})
		return
	}
	defer db.Rollback()
	err := model.InsertSelectedCircle(userId,circleId,db)
	if err != nil {
		servers.ReportFormat(c,false,"加入圈子失败",gin.H{})
		return
	}
	//更新圈子中的成员数量
	circleErr,circle := model.GetUserCircleInfo(circleId)
	if circleErr != nil {
		servers.ReportFormat(c,false,"加入圈子失败",gin.H{})
		return
	}
	err = circle.UpdateCircleMembers(circleId,true,db)
	if err != nil {
		servers.ReportFormat(c,false,"加入圈子失败",gin.H{})
		return
	}
	db.Commit()
	servers.ReportFormat(c,true,"加入圈子成功",gin.H{})
}

//退出圈子
func QuitCircle(c *gin.Context) {
	userId := tools.GetInt64(c.Query("userId"))
	circleId := tools.GetInt64(c.Query("circleId"))
	db := initSql.DEFAULTDB.Begin()
	if db == nil {
		servers.ReportFormat(c,false,"数据库打开错误，退出圈子失败",gin.H{})
		return
	}
	defer db.Rollback()
	err,circle := model.GetUserCircleInfo(circleId)
	if err != nil {
		servers.ReportFormat(c,false,"退出圈子失败",gin.H{})
		return
	}
	if circle.CreateUserId == userId {
		servers.ReportFormat(c,false,"创建人不能退出自己的圈子",gin.H{})
		return
	}
	err = circle.UpdateCircleMembers(circleId,false,db)
	if err != nil {
		servers.ReportFormat(c,false,"退出圈子失败",gin.H{})
		return
	}
	circleErr,selectedCircle := model.GetSelectedCircleByUserIdAndCircleId(userId,circleId)
	if circleErr != nil {
		servers.ReportFormat(c,false,"退出圈子失败",gin.H{})
		return
	}
	err = selectedCircle.Delete(selectedCircle.ID,db)
	if err != nil {
		servers.ReportFormat(c,false,"退出圈子失败",gin.H{})
		return
	}
	db.Commit()
	servers.ReportFormat(c,true,"退出圈子成功",gin.H{})
}

//删除圈子
func DeleteCircle(c *gin.Context) {
	userId := tools.GetInt64(c.Query("userId"))
	circleId := tools.GetInt64(c.Query("circleId"))
	db := initSql.DEFAULTDB.Begin()
	if db == nil {
		servers.ReportFormat(c,false,"数据库打开错误，删除圈子失败",gin.H{})
		return
	}
	defer db.Rollback()
	err,circle := model.GetUserCircleInfo(circleId)
	if err != nil {
		servers.ReportFormat(c,false,"删除圈子失败",gin.H{})
		return
	}
	if circle.CreateUserId != userId {
		servers.ReportFormat(c,false,"只有创建人可以删除自己的圈子",gin.H{})
		return
	}
	circleErr,selectedCircle := model.GetSelectedCircleByUserIdAndCircleId(userId,circleId)
	if circleErr != nil {
		servers.ReportFormat(c,false,"删除圈子失败",gin.H{})
		return
	}
	err = selectedCircle.Delete(selectedCircle.ID,db)
	if err != nil {
		servers.ReportFormat(c,false,"删除圈子失败",gin.H{})
		return
	}
	err = circle.Delete(circleId,db)
	if err != nil {
		servers.ReportFormat(c,false,"删除圈子失败",gin.H{})
		return
	}
	db.Commit()
	servers.ReportFormat(c,true,"删除圈子成功",gin.H{})
}