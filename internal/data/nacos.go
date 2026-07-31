package data

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v3"
	"github.com/go-kratos/kratos/v3/registry"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

// RemoteConfig 从 Nacos 拉取的远程业务配置，匹配 yaml 层级结构。
var RemoteConfig struct {
	Data struct {
		Database struct {
			User     string
			Password string
			Host     string
			Port     int
			Database string
		}
		Redis struct {
			Addr         string
			Password     string
			DB           int          `mapstructure:"db"`
			ReadTimeout  string       `mapstructure:"read_timeout"`
			WriteTimeout string       `mapstructure:"write_timeout"`
		}
		Elasticsearch struct {
			Addr string
		}
		AmapKey string `mapstructure:"amap_key"`
		SMS     struct {
			URL      string `mapstructure:"url"`
			Account  string `mapstructure:"account"`
			Password string `mapstructure:"password"`
		} `mapstructure:"sms"`
	}
}

// NacosRegistrar 服务注册客户端。
var NacosRegistrar registry.Registrar

// newNacos 从 Nacos 配置中心拉取远程配置并创建服务注册客户端。
func newNacos(addr string, port uint64, namespaceID, username, password, logDir, cacheDir, logLevel, dataID, group string) {
	sc := []constant.ServerConfig{
		{
			IpAddr: addr,
			Port:   port,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         namespaceID,
		Username:            username,
		Password:            password,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              logDir,
		CacheDir:            cacheDir,
		LogLevel:            logLevel,
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
		DataId: dataID,
		Group:  group,
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
	NacosRegistrar = nacos.New(namingClient, nacos.WithGroup(group))
	fmt.Println("nacos 链接成功")
}

// parseDuration 解析时间字符串为 Duration，失败返回默认值。
func parseDuration(s string, def time.Duration) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return def
	}
	return d
}
