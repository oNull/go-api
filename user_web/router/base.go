package router

import (
	"github.com/gin-gonic/gin"
	"py-mxshop-api/user_web/api"
)

func InitBaseRouter(Route gin.IRouter) {
	BaseRouter := Route.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
	}
}
