package data

import "time"

// Todo 对应 sf_todo 表，保存待办事项数据。
type Todo struct {
	ID         int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	Title      string    `gorm:"column:title;type:varchar(128)" db:"title" json:"title"`
	Content    string    `gorm:"column:content;type:varchar(512)" db:"content" json:"content"`
	Completed  bool      `gorm:"column:completed;type:tinyint;default:0" db:"completed" json:"completed"`
	CreateTime time.Time `gorm:"column:create_time" db:"create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" db:"update_time" json:"update_time"`
}

// TableName 指定表名。
func (Todo) TableName() string {
	return "sf_todo"
}

// User 对应 sf_user 表，保存用户个人资料。
type User struct {
	ID              int64     `gorm:"primaryKey;column:user_id" db:"user_id" json:"user_id"`
	Phone           string    `gorm:"column:phone;type:varchar(11)" db:"phone" json:"phone"`
	Email           string    `gorm:"column:email;type:varchar(64)" db:"email" json:"email"`
	AvatarURL       string    `gorm:"column:avatar_url;type:varchar(256)" db:"avatar_url" json:"avatar_url"`
	NickName        string    `gorm:"column:nick_name;type:varchar(32)" db:"nick_name" json:"nick_name"`
	RealNameAuth    int32     `gorm:"column:real_name_auth;type:tinyint" db:"real_name_auth" json:"real_name_auth"`
	AccountStatus   int32     `gorm:"column:account_status;type:tinyint" db:"account_status" json:"account_status"`
	IsEnterprise    bool      `gorm:"column:is_enterprise;type:tinyint" db:"is_enterprise" json:"is_enterprise"`
	QueryCountToday int32     `gorm:"column:query_count_today;type:int;default:0" db:"query_count_today" json:"query_count_today"`
	CreateTime      time.Time `gorm:"column:create_time" db:"create_time" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time" db:"update_time" json:"update_time"`
}

// TableName 指定表名。
func (User) TableName() string {
	return "sf_user"
}

// Address 对应 sf_user_address 表，保存用户地址簿数据。
type Address struct {
	ID            int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	UserID        int64     `gorm:"column:user_id;index" db:"user_id" json:"user_id"`
	AddrType      int32     `gorm:"column:addr_type;type:tinyint" db:"addr_type" json:"addr_type"`
	ReceiverName  string    `gorm:"column:receiver_name;type:varchar(32)" db:"receiver_name" json:"receiver_name"`
	ReceiverPhone string    `gorm:"column:receiver_phone;type:varchar(11)" db:"receiver_phone" json:"receiver_phone"`
	Province      string    `gorm:"column:province;type:varchar(32)" db:"province" json:"province"`
	City          string    `gorm:"column:city;type:varchar(32)" db:"city" json:"city"`
	District      string    `gorm:"column:district;type:varchar(32)" db:"district" json:"district"`
	DetailAddr    string    `gorm:"column:detail_addr;type:varchar(256)" db:"detail_addr" json:"detail_addr"`
	IsDefault     bool      `gorm:"column:is_default;type:tinyint" db:"is_default" json:"is_default"`
	Latitude      float64   `gorm:"column:latitude;type:decimal(10,6)" db:"latitude" json:"latitude"`
	Longitude     float64   `gorm:"column:longitude;type:decimal(10,6)" db:"longitude" json:"longitude"`
	CreateTime    time.Time `gorm:"column:create_time" db:"create_time" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time" db:"update_time" json:"update_time"`
	DelFlag       bool      `gorm:"column:del_flag;type:tinyint;default:0" db:"del_flag" json:"del_flag"`
}

// TableName 指定表名。
func (Address) TableName() string {
	return "sf_user_address"
}

// Order 对应 sf_order_history 表，保存用户历史快件记录。
type Order struct {
	ID           int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	UserID       int64     `gorm:"column:user_id;index" db:"user_id" json:"user_id"`
	OrderNo      string    `gorm:"column:order_no;type:varchar(64)" db:"order_no" json:"order_no"`
	ExpressNo    string    `gorm:"column:express_no;type:varchar(64)" db:"express_no" json:"express_no"`
	SenderName   string    `gorm:"column:sender_name;type:varchar(32)" db:"sender_name" json:"sender_name"`
	ReceiverName string    `gorm:"column:receiver_name;type:varchar(32)" db:"receiver_name" json:"receiver_name"`
	Status       string    `gorm:"column:status;type:varchar(32)" db:"status" json:"status"`
	CreateTime   time.Time `gorm:"column:create_time" db:"create_time" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" db:"update_time" json:"update_time"`
}

// TableName 指定表名。
func (Order) TableName() string {
	return "sf_order_history"
}

// IdentityUser 对应 identity_user 表，保存登录账号、角色和状态。
type IdentityUser struct {
	ID            int64      `gorm:"primaryKey;column:id" db:"id" json:"id"`
	UserNo        string     `gorm:"column:user_no;type:varchar(32);uniqueIndex" db:"user_no" json:"user_no"`
	PhoneCipher   string     `gorm:"column:phone_cipher;type:varchar(255)" db:"phone_cipher" json:"phone_cipher"`
	PhoneHash     string     `gorm:"column:phone_hash;type:char(64);uniqueIndex" db:"phone_hash" json:"phone_hash"`
	Email         string     `gorm:"column:email;type:varchar(128);uniqueIndex" db:"email" json:"email"`
	PasswordHash  string     `gorm:"column:password_hash;type:varchar(255)" db:"password_hash" json:"password_hash"`
	RoleCode      string     `gorm:"column:role_code;type:varchar(32);default:user" db:"role_code" json:"role_code"`
	AccountStatus string     `gorm:"column:account_status;type:varchar(32);default:normal" db:"account_status" json:"account_status"`
	LockedUntil   *time.Time `gorm:"column:locked_until" db:"locked_until" json:"locked_until"`
	CreatedAt     time.Time  `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

// TableName 指定表名。
func (IdentityUser) TableName() string {
	return "identity_user"
}

// IdentityRealNameAuth 对应 identity_real_name_auth 表，保存实名认证审核资料。
type IdentityRealNameAuth struct {
	UserID         int64     `gorm:"primaryKey;column:user_id" db:"user_id" json:"user_id"`
	RealNameCipher string    `gorm:"column:real_name_cipher;type:varchar(255)" db:"real_name_cipher" json:"real_name_cipher"`
	IDCardCipher   string    `gorm:"column:id_card_cipher;type:varchar(255)" db:"id_card_cipher" json:"id_card_cipher"`
	IDCardHash     string    `gorm:"column:id_card_hash;type:char(64);uniqueIndex" db:"id_card_hash" json:"id_card_hash"`
	ImageURLs      string    `gorm:"column:image_urls;type:json" db:"image_urls" json:"image_urls"`
	AuthStatus     string    `gorm:"column:auth_status;type:varchar(32);default:pending" db:"auth_status" json:"auth_status"`
	RejectReason   string    `gorm:"column:reject_reason;type:varchar(256)" db:"reject_reason" json:"reject_reason"`
	CreatedAt      time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

// TableName 指定表名。
func (IdentityRealNameAuth) TableName() string {
	return "identity_real_name_auth"
}

// ──────────────────────────────────────────────
// 寄件渠道
// ──────────────────────────────────────────────

// ShippingChannel 对应 order_shipping_channel 表。
type ShippingChannel struct {
	ID               int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	ChannelCode      string    `gorm:"column:channel_code;type:varchar(32)" db:"channel_code" json:"channel_code"`
	ChannelName      string    `gorm:"column:channel_name;type:varchar(64)" db:"channel_name" json:"channel_name"`
	ChannelDesc      string    `gorm:"column:channel_desc;type:varchar(255)" db:"channel_desc" json:"channel_desc"`
	Status           int32     `gorm:"column:status;type:tinyint" db:"status" json:"status"`
	BaseFee          float64   `gorm:"column:base_fee;type:decimal(10,2)" db:"base_fee" json:"base_fee"`
	NeedAddressGeo   int32     `gorm:"column:need_address_geo;type:tinyint" db:"need_address_geo" json:"need_address_geo"`
	NeedPickupSlot   int32     `gorm:"column:need_pickup_slot;type:tinyint" db:"need_pickup_slot" json:"need_pickup_slot"`
	NeedLockerBox    int32     `gorm:"column:need_locker_box;type:tinyint" db:"need_locker_box" json:"need_locker_box"`
	NeedServicePoint int32     `gorm:"column:need_service_point;type:tinyint" db:"need_service_point" json:"need_service_point"`
	SortNo           int32     `gorm:"column:sort_no;type:int" db:"sort_no" json:"sort_no"`
	CreatedAt        time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (ShippingChannel) TableName() string { return "order_shipping_channel" }

// ShippingChannelArea 对应 order_shipping_channel_area 表。
type ShippingChannelArea struct {
	ID                int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	ChannelCode       string    `gorm:"column:channel_code;type:varchar(32)" db:"channel_code" json:"channel_code"`
	ProvinceCode      string    `gorm:"column:province_code;type:varchar(32)" db:"province_code" json:"province_code"`
	CityCode          string    `gorm:"column:city_code;type:varchar(32)" db:"city_code" json:"city_code"`
	DistrictCode      string    `gorm:"column:district_code;type:varchar(32)" db:"district_code" json:"district_code"`
	Status            int32     `gorm:"column:status;type:tinyint" db:"status" json:"status"`
	UnavailableReason string    `gorm:"column:unavailable_reason;type:varchar(255)" db:"unavailable_reason" json:"unavailable_reason"`
	ExtraFee          float64   `gorm:"column:extra_fee;type:decimal(10,2)" db:"extra_fee" json:"extra_fee"`
	CreatedAt         time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (ShippingChannelArea) TableName() string { return "order_shipping_channel_area" }

// Locker 对应 order_locker 表。
type Locker struct {
	ID                int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	LockerCode        string    `gorm:"column:locker_code;type:varchar(64)" db:"locker_code" json:"locker_code"`
	LockerName        string    `gorm:"column:locker_name;type:varchar(128)" db:"locker_name" json:"locker_name"`
	ProvinceCode      string    `gorm:"column:province_code;type:varchar(32)" db:"province_code" json:"province_code"`
	ProvinceName      string    `gorm:"column:province_name;type:varchar(32)" db:"province_name" json:"province_name"`
	CityCode          string    `gorm:"column:city_code;type:varchar(32)" db:"city_code" json:"city_code"`
	CityName          string    `gorm:"column:city_name;type:varchar(32)" db:"city_name" json:"city_name"`
	DistrictCode      string    `gorm:"column:district_code;type:varchar(32)" db:"district_code" json:"district_code"`
	DistrictName      string    `gorm:"column:district_name;type:varchar(32)" db:"district_name" json:"district_name"`
	DetailAddress     string    `gorm:"column:detail_address;type:varchar(256)" db:"detail_address" json:"detail_address"`
	Longitude         float64   `gorm:"column:longitude;type:decimal(12,6)" db:"longitude" json:"longitude"`
	Latitude          float64   `gorm:"column:latitude;type:decimal(12,6)" db:"latitude" json:"latitude"`
	BusinessStartTime string    `gorm:"column:business_start_time;type:time" db:"business_start_time" json:"business_start_time"`
	BusinessEndTime   string    `gorm:"column:business_end_time;type:time" db:"business_end_time" json:"business_end_time"`
	Status            int32     `gorm:"column:status;type:tinyint" db:"status" json:"status"`
	CreatedAt         time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (Locker) TableName() string { return "order_locker" }

// LockerBox 对应 order_locker_box 表。
type LockerBox struct {
	ID              int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	LockerID        int64     `gorm:"column:locker_id;type:bigint" db:"locker_id" json:"locker_id"`
	BoxNo           string    `gorm:"column:box_no;type:varchar(32)" db:"box_no" json:"box_no"`
	BoxType         string    `gorm:"column:box_type;type:varchar(16)" db:"box_type" json:"box_type"`
	BoxLength       int32     `gorm:"column:box_length;type:int" db:"box_length" json:"box_length"`
	BoxWidth        int32     `gorm:"column:box_width;type:int" db:"box_width" json:"box_width"`
	BoxHeight       int32     `gorm:"column:box_height;type:int" db:"box_height" json:"box_height"`
	MaxWeight       float64   `gorm:"column:max_weight;type:decimal(8,2)" db:"max_weight" json:"max_weight"`
	BoxFee          float64   `gorm:"column:box_fee;type:decimal(10,2)" db:"box_fee" json:"box_fee"`
	Status          int32     `gorm:"column:status;type:tinyint" db:"status" json:"status"`
	ReservedOrderID int64     `gorm:"column:reserved_order_id;type:bigint" db:"reserved_order_id" json:"reserved_order_id"`
	ReservedAt      time.Time `gorm:"column:reserved_at;type:datetime" db:"reserved_at" json:"reserved_at"`
	ExpireAt        time.Time `gorm:"column:expire_at;type:datetime" db:"expire_at" json:"expire_at"`
	Version         int64     `gorm:"column:version;type:bigint" db:"version" json:"version"`
	CreatedAt       time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (LockerBox) TableName() string { return "order_locker_box" }

// ServicePoint 对应 order_service_point 表。
type ServicePoint struct {
	ID                int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	PointCode         string    `gorm:"column:point_code;type:varchar(64)" db:"point_code" json:"point_code"`
	PointName         string    `gorm:"column:point_name;type:varchar(128)" db:"point_name" json:"point_name"`
	PointType         int32     `gorm:"column:point_type;type:tinyint" db:"point_type" json:"point_type"`
	ProvinceCode      string    `gorm:"column:province_code;type:varchar(32)" db:"province_code" json:"province_code"`
	ProvinceName      string    `gorm:"column:province_name;type:varchar(32)" db:"province_name" json:"province_name"`
	CityCode          string    `gorm:"column:city_code;type:varchar(32)" db:"city_code" json:"city_code"`
	CityName          string    `gorm:"column:city_name;type:varchar(32)" db:"city_name" json:"city_name"`
	DistrictCode      string    `gorm:"column:district_code;type:varchar(32)" db:"district_code" json:"district_code"`
	DistrictName      string    `gorm:"column:district_name;type:varchar(32)" db:"district_name" json:"district_name"`
	DetailAddress     string    `gorm:"column:detail_address;type:varchar(256)" db:"detail_address" json:"detail_address"`
	Longitude         float64   `gorm:"column:longitude;type:decimal(12,6)" db:"longitude" json:"longitude"`
	Latitude          float64   `gorm:"column:latitude;type:decimal(12,6)" db:"latitude" json:"latitude"`
	BusinessStartTime string    `gorm:"column:business_start_time;type:time" db:"business_start_time" json:"business_start_time"`
	BusinessEndTime   string    `gorm:"column:business_end_time;type:time" db:"business_end_time" json:"business_end_time"`
	ContactPhoneMask  string    `gorm:"column:contact_phone_mask;type:varchar(32)" db:"contact_phone_mask" json:"contact_phone_mask"`
	Status            int32     `gorm:"column:status;type:tinyint" db:"status" json:"status"`
	ExtraFee          float64   `gorm:"column:extra_fee;type:decimal(10,2)" db:"extra_fee" json:"extra_fee"`
	CreatedAt         time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (ServicePoint) TableName() string { return "order_service_point" }

// PickupTimeSlot 对应 order_pickup_time_slot 表。
type PickupTimeSlot struct {
	ID            int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	DistrictCode  string    `gorm:"column:district_code;type:varchar(32)" db:"district_code" json:"district_code"`
	PickupDate    string    `gorm:"column:pickup_date;type:date" db:"pickup_date" json:"pickup_date"`
	SlotCode      string    `gorm:"column:slot_code;type:varchar(32)" db:"slot_code" json:"slot_code"`
	StartTime     string    `gorm:"column:start_time;type:time" db:"start_time" json:"start_time"`
	EndTime       string    `gorm:"column:end_time;type:time" db:"end_time" json:"end_time"`
	Capacity      int32     `gorm:"column:capacity;type:int" db:"capacity" json:"capacity"`
	ReservedCount int32     `gorm:"column:reserved_count;type:int" db:"reserved_count" json:"reserved_count"`
	Status        int32     `gorm:"column:status;type:tinyint" db:"status" json:"status"`
	Version       int64     `gorm:"column:version;type:bigint" db:"version" json:"version"`
	CreatedAt     time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (PickupTimeSlot) TableName() string { return "order_pickup_time_slot" }

// PickupReservation 对应 order_pickup_reservation 表。
type PickupReservation struct {
	ID           int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	OrderID      int64     `gorm:"column:order_id;type:bigint" db:"order_id" json:"order_id"`
	UserID       int64     `gorm:"column:user_id;type:bigint" db:"user_id" json:"user_id"`
	SlotID       int64     `gorm:"column:slot_id;type:bigint" db:"slot_id" json:"slot_id"`
	DistrictCode string    `gorm:"column:district_code;type:varchar(32)" db:"district_code" json:"district_code"`
	PickupDate   string    `gorm:"column:pickup_date;type:date" db:"pickup_date" json:"pickup_date"`
	SlotCode     string    `gorm:"column:slot_code;type:varchar(32)" db:"slot_code" json:"slot_code"`
	Status       int32     `gorm:"column:status;type:tinyint" db:"status" json:"status"`
	ReservedAt   time.Time `gorm:"column:reserved_at;type:datetime" db:"reserved_at" json:"reserved_at"`
	ReleasedAt   time.Time `gorm:"column:released_at;type:datetime" db:"released_at" json:"released_at"`
	ExpireAt     time.Time `gorm:"column:expire_at;type:datetime" db:"expire_at" json:"expire_at"`
	CreatedAt    time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (PickupReservation) TableName() string { return "order_pickup_reservation" }

// IdentitySecurityLog 对应 identity_security_log 表，保存登录、改密、退出等安全事件。
type IdentitySecurityLog struct {
	ID            int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	UserID        int64     `gorm:"column:user_id;index:idx_identity_security_log_user_created_at" db:"user_id" json:"user_id"`
	EventType     string    `gorm:"column:event_type;type:varchar(64)" db:"event_type" json:"event_type"`
	Result        string    `gorm:"column:result;type:varchar(32)" db:"result" json:"result"`
	IP            string    `gorm:"column:ip;type:varchar(64)" db:"ip" json:"ip"`
	DeviceID      string    `gorm:"column:device_id;type:varchar(128)" db:"device_id" json:"device_id"`
	FailureReason string    `gorm:"column:failure_reason;type:varchar(255)" db:"failure_reason" json:"failure_reason"`
	CreatedAt     time.Time `gorm:"column:created_at;index:idx_identity_security_log_user_created_at" db:"created_at" json:"created_at"`
}

// TableName 指定表名。
func (IdentitySecurityLog) TableName() string {
	return "identity_security_log"
}

// 快递订单表
type ExpressOrder struct {
	ID             int64  `gorm:"primaryKey;column:id" db:"id" json:"id" comment:"主键ID"`
	OrderNo        string `gorm:"column:order_no;type:varchar(32);uniqueIndex" db:"order_no" json:"order_no" comment:"内部订单编号，唯一"`
	TrackingNo     string `gorm:"column:tracking_no;type:varchar(32);index" db:"tracking_no" json:"tracking_no" comment:"快递运单号"`
	UserNo         string `gorm:"column:user_no;type:varchar(32);index" db:"user_no" json:"user_no" comment:"关联下单用户编号，关联IdentityUser表"`
	ExpressCompany string `gorm:"column:express_company;type:varchar(32)" db:"express_company" json:"express_company" comment:"快递公司名称"`

	// 寄件人信息
	SenderName     string `gorm:"column:sender_name;type:varchar(64)" db:"sender_name" json:"sender_name" comment:"寄件人姓名"`
	SenderPhone    string `gorm:"column:sender_phone;type:varchar(20);index" db:"sender_phone" json:"sender_phone" comment:"寄件人联系电话"`
	SenderProvince string `gorm:"column:sender_province;type:varchar(32)" db:"sender_province" json:"sender_province" comment:"寄件省份"`
	SenderCity     string `gorm:"column:sender_city;type:varchar(32)" db:"sender_city" json:"sender_city" comment:"寄件城市"`
	SenderDistrict string `gorm:"column:sender_district;type:varchar(32)" db:"sender_district" json:"sender_district" comment:"寄件区县"`
	SenderAddress  string `gorm:"column:sender_address;type:varchar(255)" db:"sender_address" json:"sender_address" comment:"寄件详细地址"`

	// 收件人信息
	ReceiverName     string `gorm:"column:receiver_name;type:varchar(64)" db:"receiver_name" json:"receiver_name" comment:"收件人姓名"`
	ReceiverPhone    string `gorm:"column:receiver_phone;type:varchar(20);index" db:"receiver_phone" json:"receiver_phone" comment:"收件人联系电话"`
	ReceiverProvince string `gorm:"column:receiver_province;type:varchar(32)" db:"receiver_province" json:"receiver_province" comment:"收件省份"`
	ReceiverCity     string `gorm:"column:receiver_city;type:varchar(32)" db:"receiver_city" json:"receiver_city" comment:"收件城市"`
	ReceiverDistrict string `gorm:"column:receiver_district;type:varchar(32)" db:"receiver_district" json:"receiver_district" comment:"收件区县"`
	ReceiverAddress  string `gorm:"column:receiver_address;type:varchar(255)" db:"receiver_address" json:"receiver_address" comment:"收件详细地址"`

	// 包裹货品信息
	ItemName      string  `gorm:"column:item_name;type:varchar(128)" db:"item_name" json:"item_name" comment:"物品名称"`
	ItemCategory  string  `gorm:"column:item_category;type:varchar(32)" db:"item_category" json:"item_category" comment:"物品分类"`
	ItemQuantity  int     `gorm:"column:item_quantity;type:int;default:1" db:"item_quantity" json:"item_quantity" comment:"物品数量"`
	Weight        float64 `gorm:"column:weight;type:decimal(10,2);comment:单位kg" db:"weight" json:"weight" comment:"包裹总重量(kg)"`
	Length        float64 `gorm:"column:length;type:decimal(10,2);comment:单位cm" db:"length" json:"length" comment:"包裹长(cm)"`
	Width         float64 `gorm:"column:width;type:decimal(10,2);comment:单位cm" db:"width" json:"width" comment:"包裹宽(cm)"`
	Height        float64 `gorm:"column:height;type:decimal(10,2);comment:单位cm" db:"height" json:"height" comment:"包裹高(cm)"`
	DeclaredValue float64 `gorm:"column:declared_value;type:decimal(10,2)" db:"declared_value" json:"declared_value" comment:"物品申报价值(保价依据)"`
	IsFragile     bool    `gorm:"column:is_fragile;type:tinyint(1);default:0" db:"is_fragile" json:"is_fragile" comment:"是否易碎品 0否 1是"`
	IsBattery     bool    `gorm:"column:is_battery;type:tinyint(1);default:0" db:"is_battery" json:"is_battery" comment:"是否含电池 0否 1是"`

	// 运费及费用明细
	FreightFee   float64 `gorm:"column:freight_fee;type:decimal(10,2);default:0" db:"freight_fee" json:"freight_fee" comment:"基础运费"`
	InsuranceFee float64 `gorm:"column:insurance_fee;type:decimal(10,2);default:0" db:"insurance_fee" json:"insurance_fee" comment:"保价服务费"`
	PackageFee   float64 `gorm:"column:package_fee;type:decimal(10,2);default:0" db:"package_fee" json:"package_fee" comment:"打包材料费"`
	TotalFee     float64 `gorm:"column:total_fee;type:decimal(10,2);default:0" db:"total_fee" json:"total_fee" comment:"订单总费用"`
	PayMethod    string  `gorm:"column:pay_method;type:varchar(16);default:prepaid" db:"pay_method" json:"pay_method" comment:"支付方式：prepaid寄付、collect到付、monthly月结"`
	PayStatus    string  `gorm:"column:pay_status;type:varchar(16);default:unpaid" db:"pay_status" json:"pay_status" comment:"支付状态：unpaid未支付、paid已支付、refunded已退款"`

	// 物流流转状态
	OrderStatus   string     `gorm:"column:order_status;type:varchar(32);default:pending;index" db:"order_status" json:"order_status" comment:"订单状态：待揽收/运输中/派送中/已签收/拒收/退回/已取消"`
	PickupTime    *time.Time `gorm:"column:pickup_time" db:"pickup_time" json:"pickup_time" comment:"快递揽收时间，为空代表未揽收"`
	DeliveredTime *time.Time `gorm:"column:delivered_time" db:"delivered_time" json:"delivered_time" comment:"签收完成时间"`
	SignedBy      string     `gorm:"column:signed_by;type:varchar(64)" db:"signed_by" json:"signed_by" comment:"签收人"`
	Remark        string     `gorm:"column:remark;type:varchar(255)" db:"remark" json:"remark" comment:"订单备注、特殊要求"`

	Packages      string `gorm:"column:packages;type:text" db:"packages" json:"packages" comment:"包装材料明细(JSON:材料+数量)"`
	PackageDetail string `gorm:"column:package_detail;type:text" db:"package_detail" json:"package_detail" comment:"费用计算过程(JSON:打包费/保价费/合计/公式)"`

	CreatedAt time.Time `gorm:"column:created_at" db:"created_at" json:"created_at" comment:"订单创建时间"`
	UpdatedAt time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at" comment:"订单最后更新时间"`
}

// TableName 自定义数据库表名
func (ExpressOrder) TableName() string {
	return "express_order"
}

// PackagingMaterial 包装材料表（打包费计费数据源，下单时从中选择材料）。
type PackagingMaterial struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" db:"id" json:"id" comment:"主键"`
	Name        string    `gorm:"column:name;type:varchar(64);uniqueIndex" db:"name" json:"name" comment:"材料名称（唯一）"`
	UnitPrice   float64   `gorm:"column:unit_price;type:decimal(10,2)" db:"unit_price" json:"unit_price" comment:"单价（元）"`
	Unit        string    `gorm:"column:unit;type:varchar(16)" db:"unit" json:"unit" comment:"计价单位（个/卷/张/件）"`
	Description string    `gorm:"column:description;type:varchar(255)" db:"description" json:"description" comment:"说明"`
	CreatedAt   time.Time `gorm:"column:created_at" db:"created_at" json:"created_at" comment:"创建时间"`
	UpdatedAt   time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at" comment:"更新时间"`
}

// TableName 自定义数据库表名
func (PackagingMaterial) TableName() string {
	return "packaging_material"
}
