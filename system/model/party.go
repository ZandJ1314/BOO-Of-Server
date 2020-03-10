package model

import (
	"boo/lib/init/initSql"
	"boo/lib/tools"
	"boo/system/controller"
	"github.com/jinzhu/gorm"
	"time"
)

//局
type Appointment struct {
	gorm.Model
	ID 				int64  		`json:"id"`
	Topic			string		`json:"topic"` //局的主题
	CreateTime		int64		`json:"create_time"` //创建时间
	StartTime		int64		`json:"start_time"` //饭局开始的时间
	CreateUserID	int64		`json:"create_user_id"` //创建人
	UpdateTime		int64		`json:"update_time"`
	AppointPlace	string		`json:"appoint_place"`  //约会地点
	Status			int8		`json:"status"` //是否约会成功{0:成功，1:等待，2:失败}
	UserNums		int16		`json:"user_nums"`  //约会成员数量
	CircleID		int64		`json:"circle_id"`  //圈子的id
	PunishOfLate	int64		`json:"punish_of_late"` //迟到惩罚
	FailedAppoint 	int64		`json:"failed_appoint"` //爽约惩罚
	ThumbsNums		int64		`json:"thumbs_nums"` //点赞数
}

func GetUserAppointInfoByCreateUserID(userId int64) (error,*Appointment){
	var userAppoint Appointment
	err := initSql.DEFAULTDB.Where("create_user_id = ?",userId).First(&userAppoint).Error
	if err != nil {
		return err,nil
	}
	return nil,&userAppoint
}

func GetUserAppointInfoByID(partyId int64) (error,*Appointment) {
	var userAppoint Appointment
	err := initSql.DEFAULTDB.Where("id = ?",partyId).First(&userAppoint).Error
	if err != nil {
		return err,nil
	}
	return nil,&userAppoint
}

func GetUserAppointInfo(CircleId int64) (error,[]*Appointment) {
	userAppoints := make([]*Appointment,0)
	err := initSql.DEFAULTDB.Where("circle_id = ?",CircleId).Find(&userAppoints).Error
	return err,userAppoints
}

func InsertPartyInfo (party *controller.PartyInfo,db *gorm.DB) (err error) {
	var u Appointment
	u.CreateTime = party.CreateTime
	u.StartTime = party.StartTime
	u.UpdateTime = party.CreateTime
	u.ID = tools.GetUUid()
	u.Topic = party.Topic
	u.CreateUserID = party.CreateUserID
	u.AppointPlace = party.AppointPlace
	u.Status = tools.EnumAppointInfo.Wait
	u.UserNums = 1
	u.CircleID = party.CircleID
	u.PunishOfLate = party.PunishOfLate
	u.FailedAppoint = party.FailedAppoint
	u.ThumbsNums = 0
	err = db.Create(&u).Error
	return err
}

func (s *Appointment) Delete(partyId int64,db *gorm.DB) (err error) {
	err = db.Delete(&Appointment{},"id = ",partyId).Error
	return err
}

func (s *Appointment) UpdateAppointUserNums(partyId int64,db *gorm.DB) (err error) {
	now := time.Now().Unix()
	num := s.UserNums+1
	party := Appointment{
		UpdateTime:    now,
		UserNums:      num,
	}
	err = db.Where("id = ",partyId).First(&s).Updates(&party).Error
	return err
}