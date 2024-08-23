package router

import (
	"github.com/gin-gonic/gin"
	"py-mxshop-api/goods_web/api/brands"
)

func InitBrandRouter(Route *gin.RouterGroup) {
	//UserRouter := Route.Group("user").Use(middlewares.JWTAuth()) // 一组
	BrandRouter := Route.Group("brands")
	{
		BrandRouter.GET("list", brands.BrandList)      // 品牌列表页
		BrandRouter.DELETE("/:id", brands.DeleteBrand) // 删除品牌
		BrandRouter.POST("create", brands.NewBrand)    //新建品牌
		BrandRouter.PUT("/:id", brands.UpdateBrand)    //修改品牌信息
	}

	CategoryBrandRouter := Route.Group("categorybrands")
	{
		CategoryBrandRouter.GET("list", brands.CategoryBrandList)      // 类别品牌列表页
		CategoryBrandRouter.DELETE("/:id", brands.DeleteCategoryBrand) // 删除类别品牌
		CategoryBrandRouter.POST("create", brands.NewCategoryBrand)    //新建类别品牌
		CategoryBrandRouter.PUT("/:id", brands.UpdateCategoryBrand)    //修改类别品牌
		CategoryBrandRouter.GET("/:id", brands.GetCategoryBrandList)   //获取分类的品牌
	}
}
