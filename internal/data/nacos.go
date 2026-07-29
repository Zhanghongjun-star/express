package data

import (
	"fmt"
	"strings"

	"shunfeng-miniprogram/internal/conf"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v3"
	"github.com/go-kratos/kratos/v3/registry"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

// RemoteConfig 从 Nacos 拉取的远程业务配置。
var RemoteConfig conf.Bootstrap

// NacosRegistrar 服务注册客户端。
var NacosRegistrar registry.Registrar

// newNacos 初始化 Nacos：先拉取远程配置，再创建服务注册客户端。
func newNacos(r *conf.Registry) {
	if r == nil || r.Nacos == nil {
		return
	}
	nc := r.Nacos

	sc := []constant.ServerConfig{
		{
			IpAddr: nc.Addr,
			Port:   nc.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         nc.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              nc.LogDir,
		CacheDir:            nc.CacheDir,
		LogLevel:            nc.LogLevel,
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("nacos 配置客户端创建失败: %v", err))
		return
	}
	fmt.Println("nacos 配置客户端创建成功")

	configYaml, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nc.DataId,
		Group:  nc.Group,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("nacos 获取配置失败: %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("nacos 配置获取成功 %s", configYaml))

	viper.Reset()
	viper.SetConfigType("yaml")
	if err = viper.ReadConfig(strings.NewReader(configYaml)); err != nil {
		fmt.Println(fmt.Sprintf("nacos 配置读取失败: %v", err))
		return
	}
	if err = viper.Unmarshal(&RemoteConfig); err != nil {
		fmt.Println(fmt.Sprintf("nacos 配置解析失败: %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("nacos 配置解析成功 %+v", RemoteConfig))

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("nacos 服务注册客户端创建失败: %v", err))
		return
	}
	NacosRegistrar = nacos.New(namingClient, nacos.WithGroup(nc.Group))
	fmt.Println("nacos 链接成功")
}
