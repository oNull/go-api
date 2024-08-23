package goods

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"py-mxshop-api/goods_web/forms"
	"py-mxshop-api/goods_web/global"
	"py-mxshop-api/goods_web/global/helpers"
	"py-mxshop-api/goods_web/proto"
	"strconv"
	"strings"
)

func List(ctx *gin.Context) {
	request := &proto.GoodsFilterRequest{}
	priceMax := ctx.DefaultQuery("price_max", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	priceMin := ctx.DefaultQuery("price_min", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	isHot := ctx.DefaultQuery("is_hot", "0")
	if isHot == "1" {
		request.IsHot = true
	}

	request.KeyWords = ctx.DefaultQuery("q", "")

	r, err := global.GoodsSrvClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		KeyWords: request.KeyWords,
		IsHot:    request.IsHot,
		PriceMax: request.PriceMax,
		PriceMin: request.PriceMin,
	})

	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		helpers.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	helpers.OkWithData(ctx, r)
	return
}

func BitchList(ctx *gin.Context) {
	ids := ctx.DefaultQuery("ids", "0")

	if ids == "0" {
		helpers.FailWithMsg(ctx, "参数错误")
		return
	}

	strSlice := strings.Split(ids, ",")
	idsInt := make([]int32, len(strSlice))

	for i, s := range strSlice {
		value, err := strconv.Atoi(s)
		if err != nil {
			helpers.FailWithMsg(ctx, "参数解析错误")
			return
		}
		idsInt[i] = int32(value)
	}
	//

	r, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: idsInt,
	})

	if err != nil {
		zap.S().Errorw("[BitchList] 查询 【商品列表】失败")
		helpers.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	helpers.OkWithData(ctx, r)
	return
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		helpers.FailWithMsg(ctx, "参数错误")
		return
	}

	r, err := global.GoodsSrvClient.GetGoodsDetail(context.WithValue(context.Background(), "ginContext", ctx), &proto.GoodInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		helpers.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	rsp := map[string]interface{}{
		"id":          r.Id,
		"name":        r.Name,
		"goods_brief": r.GoodsBrief,
		"desc":        r.GoodsDesc,
		"ship_free":   r.ShipFree,
		"images":      r.Images,
		"desc_images": r.DescImages,
		"front_image": r.GoodsFrontImage,
		"shop_price":  r.ShopPrice,
		"ctegory": map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
			"logo": r.Brand.Logo,
		},
		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"on_sale": r.OnSale,
	}

	helpers.OkWithData(ctx, rsp)
	return
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		helpers.FailWithMsg(ctx, "参数错误")
		return
	}

	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{
		Id: int32(i),
	})
	if err != nil {
		helpers.HandleGrpcErrorToHttp(err, ctx)
	}
	helpers.Ok(ctx)
	return
}

func Update(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		helpers.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.WithValue(context.Background(), "ginContext", ctx), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		Brand:           goodsForm.Brand,
	}); err != nil {
		helpers.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	helpers.Ok(ctx)
	return

}
