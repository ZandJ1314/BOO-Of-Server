package model

import (
	"boo/lib/init/initSql"
	"github.com/jinzhu/gorm"
)

//针对饭局的评价表

type Comment struct {
	gorm.Model
	ID 			int64	`json:"id"`
	UserId 		int64	`json:"user_id"`
	PartyId		int64	`json:"party_id"`
	Message 	string  `json:"message"` //留言
	ImageUrl	string	`json:"image_url"` //图片
	CreateTime	int64	`json:"create_time"` //发表时间
	UpdateTime	int64	`json:"update_time"` //更新时间
}

func GetUserPartyCommentList(partyId int64) (error,[]*Comment) {
	userComment := make([]*Comment,0)
	err := initSql.DEFAULTDB.Where("party_id = ?",partyId).Find(&userComment).Error
	return err,userComment
}
