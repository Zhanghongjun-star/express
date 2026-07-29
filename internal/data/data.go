package data

import (
	"shunfeng-miniprogram/internal/conf"

	"github.com/go-kratos/kratos/v3/log"
	"github.com/google/wire"
)

// ProviderSet 数据层依赖注入。
var ProviderSet = wire.NewSet(NewData, NewTodoRepo, NewAddressRepo, NewUserRepo, NewOrderRepo)

// Data 数据层共享客户端。
type Data struct {
	// TODO wrapped database client
}

// NewData 创建数据层。
func NewData(c *conf.Data) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}
