package initialize

import (
	"github.com/gin-gonic/gin"
	"py-mxshop-api/goods_web/middlewares"
	router2 "py-mxshop-api/goods_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/u/v1")
	router2.InitGoodsRouter(ApiGroup)
	router2.InitBrandRouter(ApiGroup)
	return Router
}
