package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"py-mxshop-api/user_web/global"
)

//type MysqlConfig struct {
//	Host string `mapstructure:"host"`
//	Port int    `mapstructure:"port"`
//}
//type ServerConfig struct {
//	Name      string      `mapstructure:"name"`
//	MysqlInfo MysqlConfig `mapstructure:"mysql"`
//}

func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func InitConfig() {
	zap.S().Info("初始化配置文件...")

	data := GetEnvInfo("DEBUG")
	zap.S().Infof("获取ENV：%s", data)
	var configFileName string
	configFileNamePrefix := "config"
	if data == "true" {
		//configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFileNamePrefix)
		configFileName = fmt.Sprintf("./user_web/config/%s-debug.yaml", configFileNamePrefix)
	} else {
		configFileName = fmt.Sprintf("./user_web/config/%s-pro.yaml", configFileNamePrefix)
		//configFileName = fmt.Sprintf("./%s-debug.yaml", configFileNamePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: %v", global.NacosConfig)

	//从nacos中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "user_web/tmp/nacos/log",
		CacheDir:            "user_web/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		panic(err)
	}
	//fmt.Println(content) //字符串 - yaml
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}
	fmt.Println(&global.ServerConfig)
}

func InitConfig2() {
	// 本地版本 配置信息没放到Nacos
	data := GetEnvInfo("DEBUG")
	zap.S().Infof("获取ENV：%s", data)
	var configFileName string
	configFileNamePrefix := "config"
	if data == "true" {
		//configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFileNamePrefix)
		configFileName = fmt.Sprintf(".config/%s-debug.yaml", configFileNamePrefix)
	} else {
		configFileName = fmt.Sprintf(".config/%s-pro.yaml", configFileNamePrefix)
		//configFileName = fmt.Sprintf("./%s-debug.yaml", configFileNamePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}

	zap.S().Infof("配置文件：%s,配置信息：%v", configFileName, global.ServerConfig)

	//go func() {
	//	v.WatchConfig()
	//	v.OnConfigChange(func(e fsnotify.Event) {
	//		fmt.Println("Config file changed:", e.Name)
	//		_ = v.ReadInConfig() // 读取配置数据
	//		_ = v.Unmarshal(global.ServerConfig)
	//		fmt.Println(global.ServerConfig)
	//	})
	//}()
	//time.Sleep(time.Second * 3000)
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件产生变化：%v", e.Name)
		_ = v.ReadInConfig() // 读取配置数据
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息为：%v", global.ServerConfig)
	})

}
