package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"py-mxshop-api/goods_web/global"
	"py-mxshop-api/goods_web/initialize"
	"py-mxshop-api/goods_web/utils"
)

func main() {
	// 1 初始化logger
	initialize.InitLogger()

	// 2 初始化配置文件
	initialize.InitConfig()

	// 3 初始化全局Validator翻译器
	err := initialize.InitTrans("zh")
	if err != nil {
		panic(err)
	}

	// 4 初始化路由
	Router := initialize.Routers()

	// 5 初始化用户SRV-GRPC服务
	initialize.InitSrvConn()

	viper.AutomaticEnv()
	debug := viper.GetBool("DEBUG")
	// 如果是线上采用 获取端口号启动
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
			zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)
		}
	}

	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)
	err = Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
	if err != nil {
		zap.S().Panic("启动失败:", err.Error())
		return
	}
}
