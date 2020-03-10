package router

import (
	"boo/system/controller"
	"github.com/gin-gonic/gin"
)

func InitIndexRouter(Router *gin.RouterGroup) {
	IndexRouter := Router.Group("index")
	{
		IndexRouter.GET("index",controller.Index)
	}
}
