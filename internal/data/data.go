package data

import (
	"database/sql"
	"time"

	"shunfeng-miniprogram/internal/conf"

	"github.com/go-kratos/kratos/v3/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet 数据层依赖注入。
var ProviderSet = wire.NewSet(NewData, NewTodoRepo, NewAddressRepo, NewUserRepo, NewOrderRepo, NewAuthRepo)

// Data 数据层共享客户端。
type Data struct {
	// DB 是 MySQL 连接池，仓储层只通过 Data 复用它。
	DB *sql.DB
	// Redis 是 Redis 客户端，用于验证码、Token 会话和黑名单。
	Redis *redis.Client
}

// NewData 创建数据层并按配置初始化 MySQL / Redis 客户端。
func NewData(c *conf.Data) (*Data, func(), error) {
	d := &Data{}
	if c != nil && c.GetDatabase() != nil && c.GetDatabase().GetSource() != "" {
		db, err := sql.Open(c.GetDatabase().GetDriver(), c.GetDatabase().GetSource())
		if err != nil {
			return nil, nil, err
		}
		d.DB = db
	}
	if c != nil && c.GetRedis() != nil && c.GetRedis().GetAddr() != "" {
		network := c.GetRedis().GetNetwork()
		if network == "" {
			network = "tcp"
		}
		readTimeout := 200 * time.Millisecond
		if c.GetRedis().GetReadTimeout() != nil {
			readTimeout = c.GetRedis().GetReadTimeout().AsDuration()
		}
		writeTimeout := 200 * time.Millisecond
		if c.GetRedis().GetWriteTimeout() != nil {
			writeTimeout = c.GetRedis().GetWriteTimeout().AsDuration()
		}
		d.Redis = redis.NewClient(&redis.Options{
			Network:      network,
			Addr:         c.GetRedis().GetAddr(),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		})
	}
	cleanup := func() {
		log.Info("closing the data resources")
		if d.Redis != nil {
			_ = d.Redis.Close()
		}
		if d.DB != nil {
			_ = d.DB.Close()
		}
	}
	return d, cleanup, nil
}
