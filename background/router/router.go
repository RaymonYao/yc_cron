package router

import (
	"github.com/gin-gonic/gin"
	. "yc_cron/api"
)

func InitRouters() (router *gin.Engine) {
	var (
		routerGroup *gin.RouterGroup
	)
	router = gin.Default()
	//router.Use(middleware.Cors()) //增加跨域请求头
	routerGroup = router.Group("")

	authRouter := routerGroup.Group("my")
	{
		authRouter.GET("test", Test)
	}
	return
}
