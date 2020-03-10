package model

import (
	"boo/lib/init/initSql"
	"boo/system/controller"
	"github.com/jinzhu/gorm"
	"time"
)

type UserCircle struct {
	gorm.Model
	CircleID 		int64 	`json:"circle_id"` //每个用户的圈子的唯一id
	CreateTime		int64 	`json:"create_time"` //圈子创建的时间
	CircleName		string 	`json:"circle_name"` //圈子的名称
	UpdateTime		int64 	`json:"update_time"`  //圈子的更新时间
	CreateUserId	int64	`json:"create_user_id"` //创建圈子的用户id
	CreateUserName	string	`json:"create_user_name"` //创建圈子的用户名称
	CircleNum		int16	`json:"circle_num"` //圈子的成员数量
	AppointNums		int16	`json:"appoint_nums"` //局子的数量
}

func (c *UserCircle) Delete(circleId int64,db *gorm.DB) (err error){
	err = db.Delete(&UserCircle{},"circle_id = ",circleId).Error
	return err
}

func GetUserCircleInfo(CircleId int64) (error,*UserCircle) {
	var userCircle UserCircle
	err := initSql.DEFAULTDB.Where("circle_id = ?",CircleId).First(&userCircle).Error
	if err != nil {
		return err,nil
	}
	return nil,&userCircle
}

func InsertCircle(circle *controller.CircleInfo,db *gorm.DB) (err error){
	var u UserCircle
	now := time.Now().Unix()
	u.CreateTime = now
	u.CircleID = circle.CircleID
	u.CircleName = circle.CreateUserName
	u.UpdateTime = now
	u.CreateUserId = circle.CreateUserId
	u.CircleNum = 1
	u.AppointNums = 0
	err = db.Create(&u).Error
	return err
}

func (c *UserCircle) UpdateCircleMembers(circleID int64,isAdd bool,db *gorm.DB) (err error) {
	//圈子成员数量加1
	now := time.Now().Unix()
	num := c.CircleNum
	if isAdd {
		num ++
	}else{
		num--
	}
	circle := UserCircle{
		UpdateTime:now,
		CircleNum:num,
	}
	err = db.Where("circle_id = ",circleID).First(&c).Updates(&circle).Error
	return err
}

func (c *UserCircle) UpdateCircleAppointNums(circleID int64,isAdd bool,db *gorm.DB) (err error) {
	//局子数量加1
	now := time.Now().Unix()
	num := c.AppointNums+1
	if isAdd {
		num ++
	}else{
		num--
	}
	circle := UserCircle{
		UpdateTime:now,
		AppointNums:num,
	}
	err = db.Where("circle_id = ",circleID).First(&c).Updates(&circle).Error
	return err
}
