package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"py-mxshop-api/goods_web/api/goods"
)

func InitGoodsRouter(Route *gin.RouterGroup) {
	// UserRouter := Route.Group("user").Use(middlewares.JWTAuth()) // 一组
	GoodsRouter := Route.Group("goods")
	zap.S().Info("配置商品相关的URL")
	{
		//UserRouter.GET("list", middlewares.JWTAuth(), api.GetUserList) // 单个
		GoodsRouter.GET("list", goods.List)
		GoodsRouter.GET("bitch/list", goods.BitchList)
		GoodsRouter.GET("/:id", goods.Detail)    //获取商品的详情
		GoodsRouter.DELETE("/:id", goods.Delete) //删除商品

		GoodsRouter.PUT("/:id", goods.Update)
	}

}
