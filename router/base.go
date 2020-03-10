package router

import (
	"boo/system/controller"
	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.POST("login",controller.LoginByWein)
	}
}
