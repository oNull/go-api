package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"py-mxshop-api/goods_web/global"
	"py-mxshop-api/goods_web/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=imooc", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 [商品服务] 失败", "msg", err.Error())
		return
	}

	// 调用接口
	userSrvClient := proto.NewGoodsClient(userConn)

	// 可以从负载均衡做 也可以下次看看连接池操作 这些连接
	global.GoodsSrvClient = userSrvClient
}
