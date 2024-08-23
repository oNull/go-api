package brands

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"py-mxshop-api/goods_web/forms"
	"py-mxshop-api/goods_web/global"
	"py-mxshop-api/goods_web/global/helpers"
	"py-mxshop-api/goods_web/proto"
	"strconv"
)

func BrandList(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.GoodsSrvClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pages:       int32(pnInt),
		PagePerNums: int32(pSizeInt),
	})

	if err != nil {
		helpers.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	//result := make([]interface{}, 0)
	//reMap := make(map[string]interface{})
	//reMap["total"] = rsp.Total
	//for _, value := range rsp.Data[pnInt : pnInt*pSizeInt+pSizeInt] {
	//	reMap := make(map[string]interface{})
	//	reMap["id"] = value.Id
	//	reMap["name"] = value.Name
	//	reMap["logo"] = value.Logo
	//	result = append(result, reMap)
	//}
	//
	//reMap["data"] = result
	helpers.OkWithData(ctx, rsp)
	return
}

func NewBrand(ctx *gin.Context) {
	brandForm := forms.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		helpers.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateBrand(context.Background(), &proto.BrandRequest{
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		helpers.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["logo"] = rsp.Logo

	helpers.OkWithData(ctx, request)
	return
}

func UpdateBrand(ctx *gin.Context) {
	brandForm := forms.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		helpers.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		helpers.FailWithMsg(ctx, "参数错误")
		return
	}

	_, err = global.GoodsSrvClient.UpdateBrand(context.Background(), &proto.BrandRequest{
		Id:   int32(i),
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		helpers.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	helpers.Ok(ctx)
	return
}

func DeleteBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteBrand(context.Background(), &proto.BrandRequest{Id: int32(i)})
	if err != nil {
		helpers.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	helpers.Ok(ctx)
	return
}
