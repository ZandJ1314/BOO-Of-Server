package initRouter


import (
	"boo/router"
	"github.com/gin-gonic/gin"
)

//初始化总路由
func InitRouter() *gin.Engine {
	var Router = gin.Default()
	//Router.Use(middleware.LoadTls())  // 打开就能玩https了
	//Router.Use(middleware.Logger()) // 如果不需要日志 请关闭这里
	//Router.Use(middleware.Cors())   // 跨域
	//Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ApiGroup := Router.Group("api") // 方便统一添加路由组前缀 多服务器上线使用
	router.InitBaseRouter(ApiGroup) //用户登录路由
	router.InitIndexRouter(ApiGroup) //首页
	router.InitCircleRouter(ApiGroup)//用户圈子相关
	router.InitPartyRouter(ApiGroup) //用户聚会相关
	//Router.Use(middleware.Logger())
	//router.InitUserRouter(ApiGroup)                  // 注册用户路由
	//router.InitBaseRouter(ApiGroup)                  // 注册基础功能路由 不做鉴权
	//router.InitMenuRouter(ApiGroup)                  // 注册menu路由
	//router.InitAuthorityRouter(ApiGroup)             // 注册角色路由
	//router.InitApiRouter(ApiGroup)                   // 注册功能api路由
	//router.InitFileUploadAndDownloadRouter(ApiGroup) // 文件上传下载功能路由
	//router.InitWorkflowRouter(ApiGroup)              // 工作流相关路由
	//router.InitCasbinRouter(ApiGroup)                // 权限相关路由
	//router.InitJwtRouter(ApiGroup)                   // jwt相关路由
	//router.InitSystemRouter(ApiGroup)                // system相关路由
	//router.InitCustomerRouter(ApiGroup)             // 客户路由
	return Router
}

