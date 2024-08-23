package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"py-mxshop-api/user_web/global"
	"py-mxshop-api/user_web/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=imooc", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 [用户服务] 失败", "msg", err.Error())
		return
	}

	// 调用接口
	userSrvClient := proto.NewUserClient(userConn)

	// 可以从负载均衡做 也可以下次看看连接池操作 这些连接
	global.UserSrvClient = userSrvClient
}

func InitSrvConn2() {
	// 从注册中心获取服务列表
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	userSrvHost := ""
	userSrvPort := 0
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == '%s'", global.ServerConfig.UserSrvInfo.Name))
	for _, service := range data {
		userSrvHost = service.Address
		userSrvPort = service.Port
		break
	}

	if userSrvHost == "" {
		zap.S().Errorw("[InitSrvConn] 连接 [用户服务] 失败")
		return
	}

	// 拨号链接用户GRPC服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[InitSrvConn] 连接 [用户服务2] 失败", "msg", err.Error())
		return
	}
	// 调用接口
	userSrvClient := proto.NewUserClient(userConn)

	// 可以从负载均衡做 也可以下次看看连接池操作 这些连接
	global.UserSrvClient = userSrvClient
}
