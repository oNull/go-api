package api

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"py-mxshop-api/user_web/forms"
	"py-mxshop-api/user_web/global"
	"py-mxshop-api/user_web/global/helpers"
	"py-mxshop-api/user_web/global/reponse"
	"py-mxshop-api/user_web/middlewares"
	"py-mxshop-api/user_web/models"
	"py-mxshop-api/user_web/proto"
	"strconv"
	"time"
)

func GetUserList(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "0")
	pSizeInt, _ := strconv.Atoi(pSize)

	list, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})

	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 [用户列表] 失败", "msg", err.Error())
		helpers.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户id：%d", currentUser.ID)
	//zap.S().Debug("获取用户列表页")
	result := make([]interface{}, 0)
	//result := make(map[string]interface{})

	for _, value := range list.Data {
		//data := make(map[string]interface{})

		user := reponse.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: time.Unix(int64(value.BirthDay), 0).Format("2006-01-02 01:02:03"),
			//Birthday: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			//Birthday: value.BirthDay,
			Gender: value.Gender,
			Mobile: value.Mobile,
		}
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthday"] = value.BirthDay
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile
		result = append(result, user)
	}

	helpers.OkWithData(ctx, result)

}

func PassWordLogin(ctx *gin.Context) {
	// 表单验证
	passwordLoginForm := forms.PassWordLoginForm{}

	if err := ctx.ShouldBindJSON(&passwordLoginForm); err != nil {
		//如何返回错误信息
		helpers.HandleValidatorError(ctx, err)
		return
	}

	// 验证验证码
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		helpers.FailWithMsg(ctx, "验证码错误")
		return
	}

	//登录逻辑 查询手机号是否存在，查询密码
	rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: passwordLoginForm.Mobile})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				helpers.FailWithMsg(ctx, "用户不存在")
				break
			default:
				helpers.FailWithMsg(ctx, "登录失败")
			}
			return
		}
	}

	// 查询到了用户 检查密码
	passRsp, passErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
		Password:          passwordLoginForm.PassWord,
		EncryptedPassword: rsp.Password,
	})

	if passErr != nil {
		helpers.FailWithMsg(ctx, "密码错误")
		return
	}

	if passRsp.Success {
		// 生成jwt token
		j := middlewares.NewJWT()
		claims := models.CustomClaims{
			ID:          uint(rsp.Id),
			NickName:    rsp.NickName,
			AuthorityId: uint(rsp.Role),
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),
				ExpiresAt: time.Now().Unix() + 60*24*30,
				Issuer:    "py-imooc",
			},
		}

		token, err := j.CreateToken(claims)

		if err != nil {
			helpers.FailWithMsg(ctx, "生成Token失败")
			return
		}

		// 登录成功过
		helpers.OkWithData(ctx, gin.H{
			"id":         rsp.Id,
			"nick_name":  rsp.NickName,
			"token":      token,
			"expired_at": (time.Now().Unix() + 60*24*30) * 1000,
		})
		return
	} else {
		// 登录成功过
		helpers.FailWithMsg(ctx, "登录失败")
		return
	}

}

func Register(ctx *gin.Context) {
	registerFrom := forms.RegisterForm{}
	if err := ctx.ShouldBindJSON(&registerFrom); err != nil {
		helpers.HandleValidatorError(ctx, err)
		return
	}

	helpers.InitRedis(0)

	value, err := helpers.RedisGet(context.Background(), registerFrom.Mobile)

	if err != nil || value != registerFrom.Code {
		helpers.FailWithMsg(ctx, "验证码错误")
		return
	}

	users, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerFrom.Mobile,
		Mobile:   registerFrom.Mobile,
		Password: registerFrom.PassWord,
	})

	if err != nil {
		zap.S().Errorw("[Register] 链接 [注册用户] 失败", "msg", err.Error())
		helpers.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 生成jwt token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(users.Id),
		NickName:    users.NickName,
		AuthorityId: uint(users.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*24*30,
			Issuer:    "py-imooc",
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		helpers.FailWithMsg(ctx, "生成Token失败")
		return
	}

	// 登录成功过
	helpers.OkWithData(ctx, gin.H{
		"id":         users.Id,
		"nick_name":  users.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*24*30) * 1000,
	})
	return

}
