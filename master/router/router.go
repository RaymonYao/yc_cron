package router

import (
	"github.com/gin-gonic/gin"
	. "yc_cron/api"
	"yc_cron/router/middleware"
)

func InitRouters() (router *gin.Engine) {
	var (
		routerGroup *gin.RouterGroup
	)
	router = gin.Default()
	router.Use(middleware.Cors()) //增加跨域请求头
	routerGroup = router.Group("")
	authRouter := routerGroup.Group("auth")
	{
		authRouter.POST("login", Login)
		authRouter.GET("captcha", Captcha)
		authRouter.GET("captcha/:captchaId", CaptchaImg)
	}
	sysRouter := routerGroup.Group("sys").Use(middleware.JWTAuth())
	{
		sysRouter.GET("getSysInfo", GetSysInfo)
		sysRouter.POST("modifyPwd", ModifyPwd)
		sysRouter.POST("getUserList", GetUserList)
		sysRouter.POST("saveUser", SaveUser)
		sysRouter.POST("setUserStatus", SetUserStatus)
	}
	groupRouter := routerGroup.Group("group").Use(middleware.JWTAuth())
	{
		groupRouter.POST("getGroupList", GetGroupList)
		groupRouter.POST("saveGroup", SaveGroup)
		groupRouter.POST("delGroup", DelGroup)
	}
	taskRouter := routerGroup.Group("task").Use(middleware.JWTAuth())
	{
		taskRouter.POST("getTaskList", GetTaskList)
		taskRouter.POST("saveTask", SaveTask)
		taskRouter.POST("delTask", DelTask)
	}
	return
}