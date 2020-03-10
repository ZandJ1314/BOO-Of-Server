package controller

import (
	"boo/lib/init/initSql"
	"boo/lib/init/initTimer"
	"boo/lib/tools"
	"boo/system/model"
	"boo/system/servers"
	"github.com/gin-gonic/gin"
	"time"
)


type PartyInfo struct {
	Topic			string		`json:"topic"` //局的主题
	CreateUserID	int64		`json:"create_user_id"` //创建人
	AppointPlace	string		`json:"appoint_place"`  //约会地点
	CircleID		int64		`json:"circle_id"`  //圈子id
	CreateTime		int64		`json:"create_time"` //饭局的时间
	StartTime		int64		`json:"start_time"` //饭局开始的时间
	PunishOfLate	int64		`json:"punish_of_late"` //迟到惩罚
	FailedAppoint 	int64		`json:"failed_appoint"` //爽约惩罚
}

type PartyDetailInfo struct {
	PartyId			int64 		`json:"party_id"` //局id
	CircleTopic		string		`json:"circle_topic"` //圈子主题
	PartyTopic		string		`json:"party_topic"` //局主题
	PartySite		string		`json:"party_site"` //局地点
	StarTime		int64		`json:"star_time"` //开始时间
	ThumbsNums		int64		`json:"thumbs_nums"`//点赞数
	UserList		[]*User 	`json:"user_list"`
	CommentList 	[]*Comment 	`json:"comment_list"` //评价
}

type Comment struct {
	UserId 		int64	`json:"user_id"`
	NickName	string	`json:"nick_name"`
	Message 	string  `json:"message"` //留言
	ImageUrl	string	`json:"image_url"` //图片
	CreateTime	int64	`json:"create_time"` //发表时间
}

type User struct {
	JoinTime		int64		`json:"join_time"` //加入时间
	UserId 			int64		`json:"user_id"`
	NickName		string		`json:"nick_name"`
	AvatarUrl		string		`json:"avatar_url"`
}

func CreateParty(c *gin.Context){
	now := time.Now().Unix()
	topic := c.Query("topic")
	site := c.Query("site")
	userId := tools.GetInt64(c.Query("userId"))
	circleId := tools.GetInt64(c.Query("circleId"))
	createTime := tools.GetInt64(c.Query("create_time"))
	startTime := tools.GetInt64(c.Query("start_time"))
	punishOfLate := tools.GetInt64(c.Query("punish_of_late"))
	failedAppoint := tools.GetInt64(c.Query("failed_appoint"))
	if userId == 0 || circleId == 0 || createTime == 0 {
		servers.ReportFormat(c,false,"创建局失败",gin.H{})
		return
	}
	db := initSql.DEFAULTDB.Begin()
	if db == nil {
		servers.ReportFormat(c,false,"数据库打开错误，创建局失败",gin.H{})
		return
	}
	defer db.Rollback()
	party := &PartyInfo{
		Topic:        topic,
		CreateUserID: userId,
		AppointPlace: site,
		CircleID:     circleId,
		CreateTime:createTime,
		PunishOfLate:punishOfLate,
		FailedAppoint:failedAppoint,
		StartTime:startTime,

	}
	err := model.InsertPartyInfo(party,db)
	if err != nil {
		servers.ReportFormat(c,false,"创建局失败",gin.H{})
		return
	}
	partyErr,newParty := model.GetUserAppointInfoByCreateUserID(userId)
	if partyErr != nil {
		servers.ReportFormat(c,false,"创建局失败",gin.H{})
		return
	}
	err = model.InsertJoinPartyInfo(userId,newParty.ID,createTime,startTime,db)
	if err != nil {
		servers.ReportFormat(c,false,"创建局失败",gin.H{})
		return
	}
	//更新圈子数据
	circleErr,circle := model.GetUserCircleInfo(circleId)
	if circleErr != nil {
		servers.ReportFormat(c,false,"创建局失败",gin.H{})
		return
	}
	err = circle.UpdateCircleAppointNums(circleId,true,db)
	if err != nil {
		servers.ReportFormat(c,false,"创建局失败",gin.H{})
		return
	}
	db.Commit()
	//对约局开始时间创建timer,到点要发推送和检查是否有人迟到
	initTimer.CreateTimer(0,startTime-now,initTimer.Msg{
		Action: initTimer.ACTION_PARTY_START_TIMER,
		Data:   newParty.ID,
	})
	servers.ReportFormat(c,true,"创建局成功",gin.H{})
}

//加入局
func JoinParty(c *gin.Context) {
	userId := tools.GetInt64(c.Query("userId"))
	partyId := tools.GetInt64(c.Query("partyId"))
	if userId == 0 || partyId == 0 {
		servers.ReportFormat(c,false,"加入局失败",gin.H{})
		return
	}
	db := initSql.DEFAULTDB.Begin()
	if db == nil {
		servers.ReportFormat(c,false,"数据库打开错误，加入局失败",gin.H{})
		return
	}
	defer db.Rollback()
	//更新局数量
	err,party := model.GetUserAppointInfoByID(partyId)
	if err != nil {
		servers.ReportFormat(c,false,"加入局失败",gin.H{})
		return
	}
	err = party.UpdateAppointUserNums(partyId,db)
	if err != nil {
		servers.ReportFormat(c,false,"加入局失败",gin.H{})
		return
	}
	err = model.InsertJoinPartyInfo(userId,partyId,time.Now().Unix(),party.StartTime,db)
	if err != nil {
		servers.ReportFormat(c,false,"加入局失败",gin.H{})
		return
	}
	db.Commit()
	servers.ReportFormat(c,true,"加入局成功",gin.H{})
}


//删除局
func DeleteParty(c *gin.Context) {
	userId := tools.GetInt64(c.Query("userId"))
	partyId := tools.GetInt64(c.Query("partyId"))
	if userId == 0 || partyId == 0 {
		servers.ReportFormat(c,false,"删除局失败",gin.H{})
		return
	}
	db := initSql.DEFAULTDB.Begin()
	if db == nil {
		servers.ReportFormat(c,false,"数据库打开错误，删除局失败",gin.H{})
		return
	}
	defer db.Rollback()
	err,party := model.GetUserAppointInfoByID(partyId)
	if err != nil {
		servers.ReportFormat(c,false,"删除局失败",gin.H{})
		return
	}
	err = party.Delete(party.ID,db)
	if err != nil {
		servers.ReportFormat(c,false,"删除局失败",gin.H{})
		return
	}
	err = model.DeleteJoinParty(partyId,db)
	if err != nil {
		servers.ReportFormat(c,false,"删除局失败",gin.H{})
		return
	}
	//更新圈子数据
	circleErr,circle := model.GetUserCircleInfo(party.CircleID)
	if circleErr != nil {
		servers.ReportFormat(c,false,"删除局失败",gin.H{})
		return
	}
	err = circle.UpdateCircleAppointNums(circle.CircleID,false,db)
	if err != nil {
		servers.ReportFormat(c,false,"删除局失败",gin.H{})
		return
	}
	db.Commit()
	servers.ReportFormat(c,true,"删除局成功",gin.H{})
}

//约局详情
func PartyDetail(c *gin.Context) {
	userId := tools.GetInt64(c.Query("userId"))
	partyId := tools.GetInt64(c.Query("partyId"))
	if userId == 0 || partyId == 0 {
		servers.ReportFormat(c,false,"获取约局详情失败",gin.H{})
		return
	}
	//获取局详情
	err,party := model.GetUserAppointInfoByID(partyId)
	if err != nil {
		servers.ReportFormat(c,false,"获取约局详情失败",gin.H{})
		return
	}
	partyErr,joinPartyList := model.GetJoinPartyInfoList(partyId)
	if partyErr != nil {
		servers.ReportFormat(c,false,"获取约局详情失败",gin.H{})
		return
	}
	//获取圈子详情
	circleErr,circle := model.GetUserCircleInfo(party.CircleID)
	if circleErr != nil {
		servers.ReportFormat(c,false,"获取约局详情失败",gin.H{})
		return
	}
	partyDetail := &PartyDetailInfo{
		PartyId:partyId,
		CircleTopic:circle.CircleName,
		PartyTopic:party.Topic,
		PartySite:party.AppointPlace,
		StarTime:party.StartTime,
		UserList:getUserByJoinPartyInfo(joinPartyList),
		CommentList:getPartyCommentList(partyId),
	}
	servers.ReportFormat(c,true,"获取成功",gin.H{"detail":partyDetail})
}

func getPartyCommentList(partyId int64) []*Comment {
	comments := make([]*Comment,0)
	err,commentList :=model.GetUserPartyCommentList(partyId)
	if err != nil {
		return comments
	}
	for _,comment := range commentList {
		err,user := model.FindUserByUserID(comment.UserId)
		if err != nil {
			continue
		}
		commentInfo := &Comment{
			UserId:     user.UserId,
			NickName:   user.Nickname,
			Message:    comment.Message,
			ImageUrl:   comment.ImageUrl,
			CreateTime: comment.CreateTime,
		}
		comments = append(comments,commentInfo)
	}
	return comments
}

func getUserByJoinPartyInfo(joinList []*model.JoinPartyInfo) []*User {
	users := make([]*User,0)
	for _,joinInfo := range joinList {
		if joinInfo == nil {
			continue
		}
		err,user := model.FindUserByUserID(joinInfo.UserId)
		if err != nil {
			continue
		}
		userInfo := &User{
			JoinTime:  joinInfo.JoinTime,
			UserId:    user.UserId,
			NickName:  user.Nickname,
			AvatarUrl: user.Avatar,
		}
		users = append(users,userInfo)
	}
	return users
}
