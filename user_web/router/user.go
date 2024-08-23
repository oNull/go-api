package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"py-mxshop-api/user_web/api"
	"py-mxshop-api/user_web/middlewares"
)

func InitUserRouter(Route *gin.RouterGroup) {
	// UserRouter := Route.Group("user").Use(middlewares.JWTAuth()) // 一组
	UserRouter := Route.Group("user")
	zap.S().Info("配置用户相关的URL")
	{
		//UserRouter.GET("list", middlewares.JWTAuth(), api.GetUserList) // 单个
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
	}

}
