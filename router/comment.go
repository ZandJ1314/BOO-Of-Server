package router

import (
	"boo/system/controller"
	"github.com/gin-gonic/gin"
)

func InitCommentRouter(Router *gin.RouterGroup){
	CommentRouter := Router.Group("comment")
	{
		CommentRouter.GET("detail",controller.PartyDetail)
	}
}
