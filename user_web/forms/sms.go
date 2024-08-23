package forms

type SmsCodeForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号加一个验证
}

type Juho struct {
	ErrorCode int `json:"error_code"`
}
