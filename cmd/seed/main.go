package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"shunfeng-miniprogram/internal/conf"
	"shunfeng-miniprogram/internal/data"

	"github.com/go-kratos/kratos/v3/config"
	"github.com/go-kratos/kratos/v3/config/file"
	"github.com/go-kratos/kratos/v3/log"
	"golang.org/x/crypto/bcrypt"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	logger := log.NewLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	log.SetDefault(logger)

	c := config.New(config.WithSource(file.NewSource(flagconf)))
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	data.InitData(bc.Registry)
	defer data.CloseData()

	if data.DB == nil {
		fmt.Println("DB not initialized")
		return
	}

	now := time.Now()
	hash := func(pw string) string {
		b, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		return string(b)
	}
	cipher := func(s string) string {
		return fmt.Sprintf("%x", []byte(s))
	}
	shaHash := func(s string) string {
		sum := sha256.Sum256([]byte(s))
		return hex.EncodeToString(sum[:])
	}

	fmt.Println("=== 开始写入种子数据 ===")

	data.DB.Exec("DELETE FROM sf_todo")
	data.DB.Exec("DELETE FROM sf_user_address")
	data.DB.Exec("DELETE FROM sf_order_history")
	data.DB.Exec("DELETE FROM order_shipping_channel_area")
	data.DB.Exec("DELETE FROM order_shipping_channel")
	data.DB.Exec("DELETE FROM order_locker_box")
	data.DB.Exec("DELETE FROM order_locker")
	data.DB.Exec("DELETE FROM order_service_point")
	data.DB.Exec("DELETE FROM order_pickup_time_slot")
	data.DB.Exec("DELETE FROM order_pickup_reservation")
	data.DB.Exec("DELETE FROM identity_real_name_auth")
	data.DB.Exec("DELETE FROM identity_security_log")
	data.DB.Exec("DELETE FROM identity_user")
	data.DB.Exec("DELETE FROM sf_user")

	// ── Users ──
	users := []struct {
		id       int64
		phone    string
		nickname string
		avatar   string
	}{
		{1, "13800138001", "张三", "https://example.com/avatar/zhangsan.png"},
		{2, "13800138002", "李四", "https://example.com/avatar/lisi.png"},
		{3, "13800138003", "王五", "https://example.com/avatar/wangwu.png"},
	}
	data.DB.Exec("ALTER TABLE sf_user AUTO_INCREMENT = 1")
	data.DB.Exec("ALTER TABLE identity_user AUTO_INCREMENT = 1")
	for _, u := range users {
		pw := fmt.Sprintf("pass%s", u.phone[len(u.phone)-4:])
		email := fmt.Sprintf("user%d@test.com", u.id)
		data.DB.Create(&data.IdentityUser{
			ID:            u.id,
			UserNo:        fmt.Sprintf("U%06d", u.id),
			PhoneCipher:   cipher(u.phone),
			PhoneHash:     shaHash(u.phone),
			Email:         email,
			PasswordHash:  hash(pw),
			RoleCode:      "user",
			AccountStatus: "normal",
			CreatedAt:     now,
			UpdatedAt:     now,
		})
		data.DB.Create(&data.User{
			ID:            u.id,
			Phone:         u.phone,
			NickName:      u.nickname,
			AvatarURL:     u.avatar,
			AccountStatus: 1,
			RealNameAuth:  0,
			IsEnterprise:  false,
			CreateTime:    now,
			UpdateTime:    now,
		})
		fmt.Printf("  user %d: %s / %s (pass: %s)\n", u.id, u.phone, u.nickname, pw)
	}

	// ── Addresses ──
	addresses := []data.Address{
		{UserID: 1, AddrType: 1, ReceiverName: "张三", ReceiverPhone: "13800138001", Province: "广东省", City: "深圳市", District: "南山区", DetailAddr: "科技南路18号深圳湾科技生态园12栋B座", IsDefault: true, Latitude: 22.525, Longitude: 113.943, CreateTime: now, UpdateTime: now},
		{UserID: 1, AddrType: 1, ReceiverName: "张三公司", ReceiverPhone: "13800138001", Province: "广东省", City: "深圳市", District: "福田区", DetailAddr: "深南大道1006号深圳国际创新中心A座", IsDefault: false, Latitude: 22.541, Longitude: 114.062, CreateTime: now, UpdateTime: now},
		{UserID: 1, AddrType: 2, ReceiverName: "张三朋友", ReceiverPhone: "13600136001", Province: "香港", City: "香港", District: "中西区", DetailAddr: "中环皇后大道中99号", IsDefault: false, Latitude: 22.282, Longitude: 114.158, CreateTime: now, UpdateTime: now},
		{UserID: 2, AddrType: 1, ReceiverName: "李四", ReceiverPhone: "13800138002", Province: "广东省", City: "广州市", District: "天河区", DetailAddr: "珠江新城华夏路10号富力中心", IsDefault: true, Latitude: 23.125, Longitude: 113.323, CreateTime: now, UpdateTime: now},
		{UserID: 2, AddrType: 1, ReceiverName: "李四父母", ReceiverPhone: "13800138002", Province: "湖南省", City: "长沙市", District: "岳麓区", DetailAddr: "梅溪湖路188号梅溪湖壹号", IsDefault: false, Latitude: 28.193, Longitude: 112.895, CreateTime: now, UpdateTime: now},
		{UserID: 3, AddrType: 1, ReceiverName: "王五", ReceiverPhone: "13800138003", Province: "上海市", City: "上海市", District: "浦东新区", DetailAddr: "张江高科技园区博云路2号", IsDefault: true, Latitude: 31.200, Longitude: 121.594, CreateTime: now, UpdateTime: now},
	}
	for _, a := range addresses {
		data.DB.Create(&a)
	}
	fmt.Printf("  addresses: %d\n", len(addresses))

	// ── Orders ──
	orders := []data.Order{
		{UserID: 1, OrderNo: "SF1234567890", ExpressNo: "SF1234567890123", SenderName: "张三", ReceiverName: "李四", Status: "已签收", CreateTime: now.Add(-72 * time.Hour), UpdateTime: now.Add(-24 * time.Hour)},
		{UserID: 1, OrderNo: "SF1234567891", ExpressNo: "SF1234567890124", SenderName: "张三", ReceiverName: "王五", Status: "运输中", CreateTime: now.Add(-24 * time.Hour), UpdateTime: now.Add(-12 * time.Hour)},
		{UserID: 1, OrderNo: "SF1234567892", ExpressNo: "SF1234567890125", SenderName: "张三公司", ReceiverName: "客户A", Status: "已揽收", CreateTime: now.Add(-6 * time.Hour), UpdateTime: now.Add(-3 * time.Hour)},
		{UserID: 2, OrderNo: "SF1234567893", ExpressNo: "SF1234567890126", SenderName: "李四", ReceiverName: "张三", Status: "已签收", CreateTime: now.Add(-120 * time.Hour), UpdateTime: now.Add(-96 * time.Hour)},
		{UserID: 2, OrderNo: "SF1234567894", ExpressNo: "SF1234567890127", SenderName: "李四", ReceiverName: "总部", Status: "运输中", CreateTime: now.Add(-8 * time.Hour), UpdateTime: now.Add(-2 * time.Hour)},
		{UserID: 3, OrderNo: "SF1234567895", ExpressNo: "SF1234567890128", SenderName: "王五", ReceiverName: "深圳客户", Status: "待揽收", CreateTime: now.Add(-1 * time.Hour), UpdateTime: now},
	}
	for _, o := range orders {
		data.DB.Create(&o)
	}
	fmt.Printf("  orders: %d\n", len(orders))

	// ── Todos ──
	todos := []data.Todo{
		{Title: "寄快递给李四", Content: "将合同文件寄给广州天河区李四", Completed: true, CreateTime: now.Add(-48 * time.Hour), UpdateTime: now.Add(-24 * time.Hour)},
		{Title: "预约上门取件", Content: "周五下午3点预约上门取件", Completed: false, CreateTime: now.Add(-12 * time.Hour), UpdateTime: now.Add(-12 * time.Hour)},
		{Title: "查询物流进度", Content: "查看SF1234567890124的物流状态", Completed: false, CreateTime: now.Add(-6 * time.Hour), UpdateTime: now.Add(-6 * time.Hour)},
		{Title: "修改默认地址", Content: "将默认地址改为公司地址", Completed: true, CreateTime: now.Add(-96 * time.Hour), UpdateTime: now.Add(-72 * time.Hour)},
		{Title: "申请电子发票", Content: "申请运费发票，金额286元", Completed: false, CreateTime: now.Add(-2 * time.Hour), UpdateTime: now},
		{Title: "实名认证", Content: "提交身份证实名认证", Completed: true, CreateTime: now.Add(-168 * time.Hour), UpdateTime: now.Add(-168 * time.Hour)},
		{Title: "查看优惠券", Content: "查看我的优惠券和运费折扣", Completed: false, CreateTime: now, UpdateTime: now},
	}
	for _, t := range todos {
		data.DB.Create(&t)
	}
	fmt.Printf("  todos: %d\n", len(todos))

	// ── Shipping Channels ──
	channels := []data.ShippingChannel{
		{ChannelCode: "SF_EXPRESS", ChannelName: "顺丰特快", ChannelDesc: "最快当日/次日达，航空运输", Status: 1, BaseFee: 23.00, NeedAddressGeo: 0, NeedPickupSlot: 0, NeedLockerBox: 0, NeedServicePoint: 0, SortNo: 1, CreatedAt: now, UpdatedAt: now},
		{ChannelCode: "SF_STANDARD", ChannelName: "顺丰标快", ChannelDesc: "2-3天送达，陆运为主", Status: 1, BaseFee: 18.00, NeedAddressGeo: 0, NeedPickupSlot: 1, NeedLockerBox: 1, NeedServicePoint: 0, SortNo: 2, CreatedAt: now, UpdatedAt: now},
		{ChannelCode: "SF_ECONOMY", ChannelName: "顺丰经济", ChannelDesc: "4-6天送达，经济实惠", Status: 1, BaseFee: 12.00, NeedAddressGeo: 0, NeedPickupSlot: 0, NeedLockerBox: 1, NeedServicePoint: 1, SortNo: 3, CreatedAt: now, UpdatedAt: now},
		{ChannelCode: "SF_SAME_CITY", ChannelName: "顺丰同城", ChannelDesc: "同城急送，2小时达", Status: 1, BaseFee: 15.00, NeedAddressGeo: 1, NeedPickupSlot: 0, NeedLockerBox: 0, NeedServicePoint: 0, SortNo: 4, CreatedAt: now, UpdatedAt: now},
		{ChannelCode: "SF_INTERNATIONAL", ChannelName: "顺丰国际", ChannelDesc: "国际快递服务", Status: 0, BaseFee: 300.00, NeedAddressGeo: 0, NeedPickupSlot: 0, NeedLockerBox: 0, NeedServicePoint: 0, SortNo: 5, CreatedAt: now, UpdatedAt: now},
	}
	for _, ch := range channels {
		data.DB.Create(&ch)
	}

	// ── Channel Areas ──
	areas := []data.ShippingChannelArea{
		{ChannelCode: "SF_EXPRESS", ProvinceCode: "440000", CityCode: "440300", DistrictCode: "440305", Status: 1, ExtraFee: 0, CreatedAt: now, UpdatedAt: now},
		{ChannelCode: "SF_EXPRESS", ProvinceCode: "440000", CityCode: "440100", DistrictCode: "440103", Status: 1, ExtraFee: 0, CreatedAt: now, UpdatedAt: now},
		{ChannelCode: "SF_STANDARD", ProvinceCode: "440000", CityCode: "440300", DistrictCode: "440305", Status: 1, ExtraFee: 0, CreatedAt: now, UpdatedAt: now},
		{ChannelCode: "SF_STANDARD", ProvinceCode: "440000", CityCode: "440100", DistrictCode: "440103", Status: 1, ExtraFee: 0, CreatedAt: now, UpdatedAt: now},
		{ChannelCode: "SF_STANDARD", ProvinceCode: "310000", CityCode: "310100", DistrictCode: "310115", Status: 1, ExtraFee: 2.00, CreatedAt: now, UpdatedAt: now},
	}
	for _, a := range areas {
		data.DB.Create(&a)
	}

	// ── Lockers ──
	lockers := []data.Locker{
		{ID: 1, LockerCode: "SZ_NANSHAN_001", LockerName: "深圳湾科技生态园9栋", ProvinceCode: "440000", ProvinceName: "广东省", CityCode: "440300", CityName: "深圳市", DistrictCode: "440305", DistrictName: "南山区", DetailAddress: "科技南路18号深圳湾科技生态园9栋1楼", Longitude: 113.943, Latitude: 22.525, BusinessStartTime: "1970-01-01 06:00:00", BusinessEndTime: "1970-01-01 23:00:00", Status: 1, CreatedAt: now, UpdatedAt: now},
		{ID: 2, LockerCode: "SZ_NANSHAN_002", LockerName: "深圳软件产业基地", ProvinceCode: "440000", ProvinceName: "广东省", CityCode: "440300", CityName: "深圳市", DistrictCode: "440305", DistrictName: "南山区", DetailAddress: "滨海大道3398号深圳软件产业基地5栋", Longitude: 113.938, Latitude: 22.530, BusinessStartTime: "1970-01-01 06:00:00", BusinessEndTime: "1970-01-01 23:30:00", Status: 1, CreatedAt: now, UpdatedAt: now},
		{ID: 3, LockerCode: "SZ_FUTIAN_001", LockerName: "深圳国际创新中心", ProvinceCode: "440000", ProvinceName: "广东省", CityCode: "440300", CityName: "深圳市", DistrictCode: "440304", DistrictName: "福田区", DetailAddress: "深南大道1006号深圳国际创新中心B座", Longitude: 114.062, Latitude: 22.541, BusinessStartTime: "1970-01-01 00:00:00", BusinessEndTime: "1970-01-01 23:59:00", Status: 1, CreatedAt: now, UpdatedAt: now},
	}
	for _, l := range lockers {
		data.DB.Create(&l)
	}

	// ── Locker Boxes ──
	epoch := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	boxTypes := []data.LockerBox{
		{LockerID: 1, BoxNo: "A01", BoxType: "S", BoxLength: 200, BoxWidth: 150, BoxHeight: 100, MaxWeight: 5.00, BoxFee: 0, Status: 1, ReservedAt: epoch, ExpireAt: epoch, CreatedAt: now, UpdatedAt: now},
		{LockerID: 1, BoxNo: "A02", BoxType: "M", BoxLength: 350, BoxWidth: 250, BoxHeight: 200, MaxWeight: 15.00, BoxFee: 2.00, Status: 1, ReservedAt: epoch, ExpireAt: epoch, CreatedAt: now, UpdatedAt: now},
		{LockerID: 1, BoxNo: "A03", BoxType: "L", BoxLength: 500, BoxWidth: 350, BoxHeight: 300, MaxWeight: 30.00, BoxFee: 5.00, Status: 1, ReservedAt: epoch, ExpireAt: epoch, CreatedAt: now, UpdatedAt: now},
		{LockerID: 1, BoxNo: "A04", BoxType: "XL", BoxLength: 700, BoxWidth: 500, BoxHeight: 400, MaxWeight: 50.00, BoxFee: 8.00, Status: 1, ReservedAt: epoch, ExpireAt: epoch, CreatedAt: now, UpdatedAt: now},
		{LockerID: 1, BoxNo: "B01", BoxType: "S", BoxLength: 200, BoxWidth: 150, BoxHeight: 100, MaxWeight: 5.00, BoxFee: 0, Status: 1, ReservedAt: epoch, ExpireAt: epoch, CreatedAt: now, UpdatedAt: now},
		{LockerID: 1, BoxNo: "B02", BoxType: "M", BoxLength: 350, BoxWidth: 250, BoxHeight: 200, MaxWeight: 15.00, BoxFee: 2.00, Status: 1, ReservedAt: epoch, ExpireAt: epoch, CreatedAt: now, UpdatedAt: now},
		{LockerID: 2, BoxNo: "A01", BoxType: "S", BoxLength: 200, BoxWidth: 150, BoxHeight: 100, MaxWeight: 5.00, BoxFee: 0, Status: 1, ReservedAt: epoch, ExpireAt: epoch, CreatedAt: now, UpdatedAt: now},
		{LockerID: 2, BoxNo: "A02", BoxType: "M", BoxLength: 350, BoxWidth: 250, BoxHeight: 200, MaxWeight: 15.00, BoxFee: 2.00, Status: 1, ReservedAt: epoch, ExpireAt: epoch, CreatedAt: now, UpdatedAt: now},
		{LockerID: 2, BoxNo: "A03", BoxType: "L", BoxLength: 500, BoxWidth: 350, BoxHeight: 300, MaxWeight: 30.00, BoxFee: 5.00, Status: 1, ReservedAt: epoch, ExpireAt: epoch, CreatedAt: now, UpdatedAt: now},
		{LockerID: 3, BoxNo: "A01", BoxType: "S", BoxLength: 200, BoxWidth: 150, BoxHeight: 100, MaxWeight: 5.00, BoxFee: 0, Status: 1, ReservedAt: epoch, ExpireAt: epoch, CreatedAt: now, UpdatedAt: now},
		{LockerID: 3, BoxNo: "A02", BoxType: "M", BoxLength: 350, BoxWidth: 250, BoxHeight: 200, MaxWeight: 15.00, BoxFee: 2.00, Status: 1, ReservedAt: epoch, ExpireAt: epoch, CreatedAt: now, UpdatedAt: now},
	}
	for _, b := range boxTypes {
		data.DB.Create(&b)
	}

	// ── Service Points ──
	points := []data.ServicePoint{
		{ID: 1, PointCode: "SP_SZ_NS_001", PointName: "顺丰速运南山区科技园营业点", PointType: 1, ProvinceCode: "440000", ProvinceName: "广东省", CityCode: "440300", CityName: "深圳市", DistrictCode: "440305", DistrictName: "南山区", DetailAddress: "科技中一路腾讯大厦1楼", Longitude: 113.934, Latitude: 22.540, BusinessStartTime: "1970-01-01 08:00:00", BusinessEndTime: "1970-01-01 20:00:00", ContactPhoneMask: "138****8001", Status: 1, ExtraFee: 0, CreatedAt: now, UpdatedAt: now},
		{ID: 2, PointCode: "SP_SZ_FT_001", PointName: "顺丰速运福田中心营业点", PointType: 1, ProvinceCode: "440000", ProvinceName: "广东省", CityCode: "440300", CityName: "深圳市", DistrictCode: "440304", DistrictName: "福田区", DetailAddress: "福华三路卓越世纪中心1号楼", Longitude: 114.058, Latitude: 22.535, BusinessStartTime: "1970-01-01 08:30:00", BusinessEndTime: "1970-01-01 20:00:00", ContactPhoneMask: "138****8002", Status: 1, ExtraFee: 0, CreatedAt: now, UpdatedAt: now},
		{ID: 3, PointCode: "SP_SZ_NS_002", PointName: "顺丰速运深圳大学代收点", PointType: 2, ProvinceCode: "440000", ProvinceName: "广东省", CityCode: "440300", CityName: "深圳市", DistrictCode: "440305", DistrictName: "南山区", DetailAddress: "南海大道3688号深圳大学南校区", Longitude: 113.938, Latitude: 22.518, BusinessStartTime: "1970-01-01 09:00:00", BusinessEndTime: "1970-01-01 19:00:00", ContactPhoneMask: "138****8003", Status: 1, ExtraFee: 0, CreatedAt: now, UpdatedAt: now},
	}
	for _, p := range points {
		data.DB.Create(&p)
	}

	// ── Pickup Time Slots ──
	dates := []string{time.Now().Format("2006-01-02"), time.Now().Add(24 * time.Hour).Format("2006-01-02")}
	slots := []struct{ code, start, end string }{
		{"MORNING", "1970-01-01 09:00:00", "1970-01-01 12:00:00"},
		{"AFTERNOON", "1970-01-01 12:00:00", "1970-01-01 15:00:00"},
		{"EVENING", "1970-01-01 15:00:00", "1970-01-01 18:00:00"},
	}
	for _, date := range dates {
		for _, s := range slots {
			data.DB.Create(&data.PickupTimeSlot{
				DistrictCode:  "440305",
				PickupDate:    date,
				SlotCode:      s.code,
				StartTime:     s.start,
				EndTime:       s.end,
				Capacity:      100,
				ReservedCount: 0,
				Status:        1,
				Version:       0,
				CreatedAt:     now,
				UpdatedAt:     now,
			})
		}
	}

	fmt.Printf("  channels: %d, areas: %d\n", len(channels), len(areas))
	fmt.Printf("  lockers: %d, boxes: %d\n", len(lockers), len(boxTypes))
	fmt.Printf("  service_points: %d\n", len(points))
	fmt.Printf("  pickup_slots: %d\n", len(dates)*len(slots))

	fmt.Println("=== 种子数据写入完成 ===")
}
