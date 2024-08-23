package initialize

import (
	"github.com/gin-gonic/gin"
	"py-mxshop-api/user_web/middlewares"
	router2 "py-mxshop-api/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/u/v1")
	router2.InitUserRouter(ApiGroup)
	router2.InitBaseRouter(ApiGroup)
	router2.InitSmsRouter(ApiGroup)
	return Router
}
