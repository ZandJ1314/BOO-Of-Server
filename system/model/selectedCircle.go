package model

import (
	"boo/lib/init/initSql"
	"boo/lib/tools"
	"github.com/jinzhu/gorm"
	"time"
)

//该表是圈子和用户之间的关联表
type SelectedCircle struct {
	gorm.Model
	ID 			int64		`json:"id"` //选圈的唯一id
	UserID 		int64		`json:"user_id"` //用户id
	JoinTime	int64		`json:"join_time"` //用户加入圈子时间
	CircleID 	int64		`json:"circle_id"` //裙子id
}

func (s *SelectedCircle) Delete(id int64,db *gorm.DB) (err error) {
	err = db.Delete(&SelectedCircle{},"id = ",id).Error
	return err
}

func InsertSelectedCircle(userID,circleID int64,db *gorm.DB) (err error){
	var circle SelectedCircle
	now := time.Now().Unix()
	circle.ID = tools.GetUUid()
	circle.UserID = userID
	circle.JoinTime = now
	circle.CircleID = circleID
	err = db.Create(&circle).Error
	return err
}

func GetSelectedCircleByUserID(userId int64) (error,[]*SelectedCircle){
	selectedCircle := make([]*SelectedCircle,0)
	err := initSql.DEFAULTDB.Where("user_id = ?",userId).Find(&selectedCircle).Error
	return err,selectedCircle
}

func GetSelectedCircleByUserIdAndCircleId(userId,circleId int64) (error,*SelectedCircle) {
	var selectCircle SelectedCircle
	err := initSql.DEFAULTDB.Where("user_id = ? AND circle_id = ?",userId,circleId).First(&selectCircle).Error
	return err,&selectCircle
}

func GetSelectedUserByCircleId(circleID int64) (error,[]*SelectedCircle) {
	selectedCircle := make([]*SelectedCircle,0)
	err := initSql.DEFAULTDB.Where("circle_id = ?",circleID).Find(&selectedCircle).Error
	return err,selectedCircle
}

func GetSelectedUserListByCircleList(circleList []*SelectedCircle) []int64 {
	users := make([]int64,0)
	for _,conf := range circleList {
		users = append(users,conf.UserID)
	}
	return users
}
