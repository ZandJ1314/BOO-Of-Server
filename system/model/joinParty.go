package model

import (
	"boo/lib/init/initSql"
	"boo/lib/tools"
	"github.com/jinzhu/gorm"
)

type JoinPartyInfo struct {
	gorm.Model
	ID 				int64	`json:"id"`
	UserId 			int64	`json:"user_id"`
	PartyId 		int64	`json:"party_id"`
	JoinTime		int64	`json:"join_time"`
	StartTime		int64	`json:"start_time"` //饭局开始的时间
	IsLate			bool	`json:"is_late"` //是否迟到
	IsFaiedAppoint	bool	`json:"is_faied_appoint"` //是否爽约
}

func GetJoinPartyInfoList(partyId int64) (error,[]*JoinPartyInfo) {
	userJoins := make([]*JoinPartyInfo,0)
	err := initSql.DEFAULTDB.Where("party_id = ?",partyId).Find(&userJoins).Error
	return err,userJoins
}

func InsertJoinPartyInfo(userId,partyId,joinTime,starTime int64,db *gorm.DB)(err error) {
	var u JoinPartyInfo
	u.ID = tools.GetUUid()
	u.UserId = userId
	u.PartyId = partyId
	u.JoinTime = joinTime
	u.StartTime = starTime
	u.IsLate = false
	u.IsFaiedAppoint = false
	err = db.Create(&u).Error
	return err
}

func DeleteJoinParty(partyId int64,db *gorm.DB) (err error) {
	err = db.Delete(&JoinPartyInfo{},"party_id = ",partyId).Error
	return err
}
