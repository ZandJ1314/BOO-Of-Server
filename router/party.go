package router

import (
	"boo/system/controller"
	"github.com/gin-gonic/gin"
)

func InitPartyRouter(Router *gin.RouterGroup){
	PartyRouter := Router.Group("party")
	{
		PartyRouter.GET("detail",controller.PartyDetail)
		PartyRouter.POST("create",controller.CreateParty)
		PartyRouter.POST("join",controller.JoinParty)
		PartyRouter.POST("delete",controller.DeleteParty)
	}
}
