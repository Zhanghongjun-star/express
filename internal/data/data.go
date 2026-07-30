package data

import (
	"context"
	"fmt"
	"time"

	"shunfeng-miniprogram/internal/conf"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/durationpb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	RDB      *redis.Client
	ESClient *elasticsearch.Client
)

// ProviderSet 数据层依赖注入。
var ProviderSet = wire.NewSet(
	NewTodoRepo,
	NewAddressRepo,
	NewUserRepo,
	NewOrderRepo,
	// 查快递模块仓储
	NewOrderFollowRepo,
	NewOrderLabelRepo,
	NewOrderShareRepo,
	NewUserMessageRepo,
)
var ProviderSet = wire.NewSet(NewTodoRepo, NewAddressRepo, NewUserRepo, NewOrderRepo, NewChannelRepo, NewLockerRepo, NewLockerBoxRepo, NewServicePointRepo, NewPickupRepo, NewAuthRepo, NewFreightRepo, NewAmapGeocoder)

// InitData 初始化所有存储客户端。
func InitData(r *conf.Registry) {
	if r == nil || r.Nacos == nil {
		return
	}
	nc := r.Nacos
	newNacos(nc.Addr, nc.Port, nc.NamespaceId, nc.Username, nc.Password, nc.LogDir, nc.CacheDir, nc.LogLevel, nc.DataId, nc.Group)

	d := RemoteConfig.Data
	newGorm(&conf.Data_Database{
		User:     d.Database.User,
		Password: d.Database.Password,
		Host:     d.Database.Host,
		Port:     int32(d.Database.Port),
		Database: d.Database.Database,
	})
	newRedis(&conf.Data_Redis{
		Addr:         d.Redis.Addr,
		Password:     d.Redis.Password,
		Db:           int32(d.Redis.DB),
		ReadTimeout:  durationpb.New(parseDuration(d.Redis.ReadTimeout, 200*time.Millisecond)),
		WriteTimeout: durationpb.New(parseDuration(d.Redis.WriteTimeout, 200*time.Millisecond)),
	})
	newES(&conf.Data_Elasticsearch{
		Addr: d.Elasticsearch.Addr,
	})
}

// newGorm 初始化 MySQL 连接
func newGorm(c *conf.Data_Database) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(fmt.Sprintf("mysql 链接失败: %v", err))
		return
	}
	// 自动建表/迁移：包含已有表和查快递模块新增的 4 张表
	if err = db.AutoMigrate(
		&Todo{},
		&User{},
		&Address{},
		&Order{},
		&IdentityUser{},
		&IdentityRealNameAuth{},
		// 查快递模块：关注·标签·分享·消息
		&OrderFollow{},
		&OrderLabel{},
		&OrderShare{},
		&UserMessage{},
		&ShippingChannel{},
		&ShippingChannelArea{},
		&Locker{},
		&LockerBox{},
		&ServicePoint{},
		&PickupTimeSlot{},
		&PickupReservation{},
		&IdentitySecurityLog{},
	); err != nil {
		fmt.Println(fmt.Sprintf("mysql 链接失败: %v", err))
		return
	}
	fmt.Println("mysql 链接成功")
	DB = db
}

// newRedis 初始化 Redis 连接
func newRedis(c *conf.Data_Redis) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           int(c.Db),
		ReadTimeout:  c.ReadTimeout.AsDuration(),
		WriteTimeout: c.WriteTimeout.AsDuration(),
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		fmt.Println(fmt.Sprintf("redis 链接失败: %v", err))
		return
	}
	fmt.Println("redis 链接成功")
	RDB = rdb
}

// newES 初始化 Elasticsearch 连接
func newES(c *conf.Data_Elasticsearch) {
	if c == nil || c.Addr == "" {
		return
	}
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{c.Addr},
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("elasticsearch 链接失败: %v", err))
		return
	}
	fmt.Println("es 链接成功")
	ESClient = esClient
}

// CloseData 关闭所有存储连接。
func CloseData() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
	if RDB != nil {
		RDB.Close()
	}
	fmt.Println("关闭数据资源")
}
