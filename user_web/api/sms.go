package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"net/url"
	"py-mxshop-api/user_web/forms"
	"py-mxshop-api/user_web/global"
	"py-mxshop-api/user_web/global/helpers"
	"py-mxshop-api/user_web/proto"
	"time"
)

// SendAlSms SendSms 阿里云短信发送
func SendAlSms(ctx *gin.Context) {
	smsCodeForm := forms.SmsCodeForm{}
	if err := ctx.ShouldBindJSON(&smsCodeForm); err != nil {
		//如何返回错误信息
		helpers.HandleValidatorError(ctx, err)
		return
	}
	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", "xxxx", "xxx")
	if err != nil {
		panic(err)
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = "xxx"                            //手机号
	request.QueryParams["SignName"] = "xxx"                                //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "xxx"                            //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + RandCode() + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	//  fmt.Print(response)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
}

func SendJhSms(ctx *gin.Context) {
	// 表单验证
	smsCodeForm := forms.SmsCodeForm{}

	if err := ctx.ShouldBindJSON(&smsCodeForm); err != nil {
		//如何返回错误信息
		helpers.HandleValidatorError(ctx, err)
		return
	}

	is := findUser(smsCodeForm.Mobile)

	if !is {
		helpers.FailWithMsg(ctx, "用户已经被注册")
		return
	}

	//初始化参数
	param := url.Values{}
	code := RandCode()
	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("mobile", smsCodeForm.Mobile)        //接受短信的用户手机号码
	param.Set("tpl_id", "68609")                   //您申请的短信模板ID，根据实际情况修改
	param.Set("tpl_value", "#code#="+code)         //您设置的模板变量，根据实际情况
	param.Set("key", global.ServerConfig.JuHe.Key) //应用APPKEY(应用详细页查询)

	//发送请求
	data, err := helpers.HttpGet(global.ServerConfig.JuHe.Url, param)
	if err != nil {
		zap.S().Errorw("[SendJhSms] 发送短信失败", "msg", err.Error())
		helpers.FailWithAll(ctx, "短信发送失败", 401)
		return
	}

	juhe := forms.Juho{}
	err = json.Unmarshal(data, &juhe)
	if err != nil {
		zap.S().Errorw("[SendJhSms] 发送短信Json失败", "msg", err.Error())
		helpers.FailWithAll(ctx, "短信发送失败2", 401)
		return
	}
	if juhe.ErrorCode != 0 {
		helpers.FailWithMsg(ctx, "短信发送失败")
		return
	}

	// 保存验证码到redis
	helpers.InitRedis(0)
	err = helpers.RedisSet(context.Background(), smsCodeForm.Mobile, code)

	if err != nil {
		helpers.FailWithMsg(ctx, "短信发送失败")
		return
	}

}

func findUser(mobile string) bool {
	//登录逻辑 查询手机号是否存在
	rsp, _ := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: mobile})
	if rsp != nil {
		return false
	}

	return true
}

// RandCode 产生六位数验证码
func RandCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return code
}
