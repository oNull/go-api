package global

import (
	ut "github.com/go-playground/universal-translator"
	"py-mxshop-api/user_web/config"
	"py-mxshop-api/user_web/proto"
)

var (
	Trans         ut.Translator
	ServerConfig  = &config.SeverConfig{}
	UserSrvClient proto.UserClient

	NacosConfig = &config.NacosConfig{}
)
