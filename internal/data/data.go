package data

import (
	"context"
	"fmt"

	"shunfeng-miniprogram/internal/conf"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	RDB      *redis.Client
	ESClient *elasticsearch.Client
)

// ProviderSet 数据层依赖注入。
var ProviderSet = wire.NewSet(NewTodoRepo, NewAddressRepo, NewUserRepo, NewOrderRepo, NewChannelRepo, NewLockerRepo, NewLockerBoxRepo, NewServicePointRepo, NewPickupRepo)

// InitData 初始化所有存储客户端。
func InitData(c *conf.Data, r *conf.Registry) {
	newGorm(c.Database)
	newRedis(c.Redis)
	newES(c.Elasticsearch)
	newNacos(r)
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
	if err = db.AutoMigrate(
		&Todo{},
		&User{},
		&Address{},
		&Order{},
		&IdentityRealNameAuth{},
		&ShippingChannel{},
		&ShippingChannelArea{},
		&Locker{},
		&LockerBox{},
		&ServicePoint{},
		&PickupTimeSlot{},
		&PickupReservation{},
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
