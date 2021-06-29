package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"yc_cron/config"
	"yc_cron/router"
)

func startServer() {
	var (
		server    *http.Server
		muxRouter *gin.Engine
		addr      string
		err       error
	)
	if config.GConfig.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	muxRouter = router.InitRouters() //初始化路由
	addr = fmt.Sprintf(":%d", config.GConfig.Port)
	server = &http.Server{
		Addr:           addr,
		Handler:        muxRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	time.Sleep(10 * time.Microsecond)
	fmt.Println("服务启动, 监听端口 " + addr)
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	config.InitConfig()
	startServer()
}
