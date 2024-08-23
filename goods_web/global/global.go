package global

import (
	ut "github.com/go-playground/universal-translator"
	"py-mxshop-api/goods_web/config"
	"py-mxshop-api/goods_web/proto"
)

var (
	Trans          ut.Translator
	ServerConfig   = &config.SeverConfig{}
	GoodsSrvClient proto.GoodsClient

	NacosConfig = &config.NacosConfig{}
)
