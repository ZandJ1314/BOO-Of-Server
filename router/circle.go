package router

import (
	"boo/system/controller"
	"github.com/gin-gonic/gin"
)


func InitCircleRouter(Router *gin.RouterGroup) {
	CircleRouter := Router.Group("circle")
	{
		CircleRouter.POST("create",controller.CreateCircle)
		CircleRouter.POST("join",controller.JoinCircle)
		CircleRouter.POST("quit",controller.QuitCircle)
		CircleRouter.POST("delete",controller.DeleteCircle)
	}
}

