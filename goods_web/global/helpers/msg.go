package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"py-mxshop-api/goods_web/global"
	"strings"
)

func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	// 将grpc的code转换成http状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				FailWithMsg(c, e.Message())
				break
			case codes.Internal:
				FailWithMsg(c, "内部错误")
				break
			case codes.InvalidArgument:
				FailWithMsg(c, "参数错误")
				break
			case codes.Unavailable:
				FailWithMsg(c, "用户服务不可用")
				break
			default:
				FailWithMsg(c, "其他错误")
				break
			}
			return
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	// 获取validator.ValidationErrors类型的errors
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		FailWithStructAll(c, err.Error(), 401)
		return
	}
	// validator.ValidationErrors类型错误则进行翻译
	FailWithStructAll(c, RemoveTopStruct(errs.Translate(global.Trans)), 401)
}

type baseResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  interface{} `json:"msg"`
}

const (
	ERROR   = 7
	FAIL    = 400
	SUCCESS = 200
)

// Result 通用返回格式
func Result(code int, data interface{}, msg interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, baseResponse{
		code,
		data,
		msg,
	})
}

// Ok 操作成功返回
func Ok(c *gin.Context) {
	message := "Successful"
	Result(SUCCESS, nil, message, c)
}

// OkWithMsg 操作成功返回带自定义文字
func OkWithMsg(c *gin.Context, message string) {
	Result(SUCCESS, nil, message, c)
}

// OkWithData 操作成功返回带自定义数据
func OkWithData(c *gin.Context, data any) {
	message := "Successful"
	Result(SUCCESS, data, message, c)
}

// Fail 操作失败返回
func Fail(c *gin.Context) {
	message := "Operation failure"
	Result(ERROR, nil, message, c)
}

// FailWithMsg 操作失败返回指定文字
func FailWithMsg(c *gin.Context, message string) {
	Result(FAIL, nil, message, c)
}

// FailWithCode 操作失败返回指定CODE
func FailWithCode(c *gin.Context, code int) {
	message := "Operation failure"
	Result(code, nil, message, c)
}

// FailWithAll 操作失败返回指定CODE 指定文字
func FailWithAll(c *gin.Context, message string, code int) {
	Result(code, nil, message, c)
}

// FailWithStructAll 操作失败返回指定CODE 指定文字
func FailWithStructAll(c *gin.Context, message interface{}, code int) {
	Result(code, nil, message, c)
}
