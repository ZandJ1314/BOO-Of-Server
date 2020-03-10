package main

import (
	"boo/lib/conf"
	"boo/lib/init/initLog"
	"boo/lib/init/initRouter"
	"boo/lib/init/initSql"
	"boo/lib/init/initTable"
	"boo/lib/init/initTimer"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main(){
	initLog.InitLog() //对日志进行初始化
	//对数据库进行初始化
	db := initSql.InitMysql(conf.GinAdminconfig.MysqlAdmin)
	//对数据库表进行注册
	initTable.RegistTable(db)
	//程序关闭前关闭数据库连接
	defer initSql.DEFAULTDB.Close()
	//timer的初始化
	initTimer.InitTimer()
	//注册路由
	Router := initRouter.InitRouter()
	Run(Router)
}

func Run(Router *gin.Engine) {
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", conf.GinAdminconfig.System.Addr),
		Handler:        Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	time.Sleep(10 * time.Microsecond)
	fmt.Printf("Addr %v",s.Addr)
	err := s.ListenAndServe()
	if err != nil {
		fmt.Errorf("listen server err %v",err)
	}
}