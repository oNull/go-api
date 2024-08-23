package router

import (
	"github.com/gin-gonic/gin"
	"py-mxshop-api/user_web/api"
)

func InitSmsRouter(Route *gin.RouterGroup) {
	UserRouter := Route.Group("sms")

	{
		UserRouter.POST("jh", api.SendJhSms)
	}

}
