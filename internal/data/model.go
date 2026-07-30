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

// OrderFollow 对应 sf_order_follow 表，记录用户关注的订单。
// 用户+订单唯一索引(uk_user_order)防重复关注，deleted_at 支持逻辑删除。
type OrderFollow struct {
	ID        int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	UserID    int64     `gorm:"column:user_id;index:idx_user_id;uniqueIndex:uk_user_order;not null" db:"user_id" json:"user_id"`
	OrderID   int64     `gorm:"column:order_id;uniqueIndex:uk_user_order;not null" db:"order_id" json:"order_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null" db:"created_at" json:"created_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;default:null" db:"deleted_at" json:"deleted_at"`
}

func (OrderFollow) TableName() string {
	return "sf_order_follow"
}

// OrderLabel 对应 sf_order_label 表，记录用户给订单设置的自定义标签。
// upsert 写入：相同 user_id+order_id 覆盖旧标签内容。
type OrderLabel struct {
	ID        int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	UserID    int64     `gorm:"column:user_id;index:idx_user_id;uniqueIndex:uk_user_order;not null" db:"user_id" json:"user_id"`
	OrderID   int64     `gorm:"column:order_id;uniqueIndex:uk_user_order;not null" db:"order_id" json:"order_id"`
	Content   string    `gorm:"column:content;type:varchar(100);not null" db:"content" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;not null" db:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;default:null" db:"deleted_at" json:"deleted_at"`
}

func (OrderLabel) TableName() string {
	return "sf_order_label"
}

// OrderShare 对应 sf_order_share 表，记录用户生成的订单分享。
// share_code 为 32 字符高熵随机码，公开查看时按 code 查询。
type OrderShare struct {
	ID           int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	ShareCode    string    `gorm:"column:share_code;type:varchar(32);uniqueIndex;not null" db:"share_code" json:"share_code"`
	UserID       int64     `gorm:"column:user_id;not null" db:"user_id" json:"user_id"`
	OrderID      int64     `gorm:"column:order_id;not null" db:"order_id" json:"order_id"`
	ShowSender   bool      `gorm:"column:show_sender;type:tinyint(1);default:0" db:"show_sender" json:"show_sender"`
	ShowReceiver bool      `gorm:"column:show_receiver;type:tinyint(1);default:0" db:"show_receiver" json:"show_receiver"`
	ShowPhone    bool      `gorm:"column:show_phone;type:tinyint(1);default:0" db:"show_phone" json:"show_phone"`
	ShowStatus   bool      `gorm:"column:show_status;type:tinyint(1);default:1" db:"show_status" json:"show_status"`
	Status       int32     `gorm:"column:status;type:tinyint;default:1" db:"status" json:"status"`
	ExpiresAt    time.Time `gorm:"column:expires_at;index;not null" db:"expires_at" json:"expires_at"`
	CreatedAt    time.Time `gorm:"column:created_at;not null" db:"created_at" json:"created_at"`
}

func (OrderShare) TableName() string {
	return "sf_order_share"
}

// UserMessage 对应 sf_user_message 表，记录用户收到的业务消息。
// uk_business_msg(business_type+business_id+message_type) 保证幂等。
type UserMessage struct {
	ID           int64      `gorm:"primaryKey;column:id" db:"id" json:"id"`
	UserID       int64      `gorm:"column:user_id;not null;index:idx_user_read_time" db:"user_id" json:"user_id"`
	MessageType  string     `gorm:"column:message_type;type:varchar(32);not null;uniqueIndex:uk_business_msg;priority:3" db:"message_type" json:"message_type"`
	Title        string     `gorm:"column:title;type:varchar(128);not null" db:"title" json:"title"`
	Content      string     `gorm:"column:content;type:text;not null" db:"content" json:"content"`
	BusinessType string     `gorm:"column:business_type;type:varchar(32);not null;uniqueIndex:uk_business_msg;priority:1" db:"business_type" json:"business_type"`
	BusinessID   string     `gorm:"column:business_id;type:varchar(64);not null;uniqueIndex:uk_business_msg;priority:2" db:"business_id" json:"business_id"`
	Priority     int32      `gorm:"column:priority;type:tinyint;default:0" db:"priority" json:"priority"`
	IsRead       bool       `gorm:"column:is_read;type:tinyint(1);default:0;index:idx_user_read_time" db:"is_read" json:"is_read"`
	ReadAt       *time.Time `gorm:"column:read_at;default:null" db:"read_at" json:"read_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index;default:null" db:"deleted_at" json:"deleted_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;index:idx_user_read_time" db:"created_at" json:"created_at"`
}

func (UserMessage) TableName() string {
	return "sf_user_message"
}

// ──────────────────────────────────────────────
// 寄件渠道
// ──────────────────────────────────────────────

// ShippingChannel 对应 order_shipping_channel 表。
type ShippingChannel struct {
	ID              int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	ChannelCode     string    `gorm:"column:channel_code;type:varchar(32)" db:"channel_code" json:"channel_code"`
	ChannelName     string    `gorm:"column:channel_name;type:varchar(64)" db:"channel_name" json:"channel_name"`
	ChannelDesc     string    `gorm:"column:channel_desc;type:varchar(255)" db:"channel_desc" json:"channel_desc"`
	Status          int32     `gorm:"column:status;type:tinyint" db:"status" json:"status"`
	BaseFee         float64   `gorm:"column:base_fee;type:decimal(10,2)" db:"base_fee" json:"base_fee"`
	NeedAddressGeo  int32     `gorm:"column:need_address_geo;type:tinyint" db:"need_address_geo" json:"need_address_geo"`
	NeedPickupSlot  int32     `gorm:"column:need_pickup_slot;type:tinyint" db:"need_pickup_slot" json:"need_pickup_slot"`
	NeedLockerBox   int32     `gorm:"column:need_locker_box;type:tinyint" db:"need_locker_box" json:"need_locker_box"`
	NeedServicePoint int32    `gorm:"column:need_service_point;type:tinyint" db:"need_service_point" json:"need_service_point"`
	SortNo          int32     `gorm:"column:sort_no;type:int" db:"sort_no" json:"sort_no"`
	CreatedAt       time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
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
	ID               int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	LockerCode       string    `gorm:"column:locker_code;type:varchar(64)" db:"locker_code" json:"locker_code"`
	LockerName       string    `gorm:"column:locker_name;type:varchar(128)" db:"locker_name" json:"locker_name"`
	ProvinceCode     string    `gorm:"column:province_code;type:varchar(32)" db:"province_code" json:"province_code"`
	ProvinceName     string    `gorm:"column:province_name;type:varchar(32)" db:"province_name" json:"province_name"`
	CityCode         string    `gorm:"column:city_code;type:varchar(32)" db:"city_code" json:"city_code"`
	CityName         string    `gorm:"column:city_name;type:varchar(32)" db:"city_name" json:"city_name"`
	DistrictCode     string    `gorm:"column:district_code;type:varchar(32)" db:"district_code" json:"district_code"`
	DistrictName     string    `gorm:"column:district_name;type:varchar(32)" db:"district_name" json:"district_name"`
	DetailAddress    string    `gorm:"column:detail_address;type:varchar(256)" db:"detail_address" json:"detail_address"`
	Longitude        float64   `gorm:"column:longitude;type:decimal(12,6)" db:"longitude" json:"longitude"`
	Latitude         float64   `gorm:"column:latitude;type:decimal(12,6)" db:"latitude" json:"latitude"`
	BusinessStartTime string   `gorm:"column:business_start_time;type:time" db:"business_start_time" json:"business_start_time"`
	BusinessEndTime  string    `gorm:"column:business_end_time;type:time" db:"business_end_time" json:"business_end_time"`
	Status           int32     `gorm:"column:status;type:tinyint" db:"status" json:"status"`
	CreatedAt        time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
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
	ID               int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	PointCode        string    `gorm:"column:point_code;type:varchar(64)" db:"point_code" json:"point_code"`
	PointName        string    `gorm:"column:point_name;type:varchar(128)" db:"point_name" json:"point_name"`
	PointType        int32     `gorm:"column:point_type;type:tinyint" db:"point_type" json:"point_type"`
	ProvinceCode     string    `gorm:"column:province_code;type:varchar(32)" db:"province_code" json:"province_code"`
	ProvinceName     string    `gorm:"column:province_name;type:varchar(32)" db:"province_name" json:"province_name"`
	CityCode         string    `gorm:"column:city_code;type:varchar(32)" db:"city_code" json:"city_code"`
	CityName         string    `gorm:"column:city_name;type:varchar(32)" db:"city_name" json:"city_name"`
	DistrictCode     string    `gorm:"column:district_code;type:varchar(32)" db:"district_code" json:"district_code"`
	DistrictName     string    `gorm:"column:district_name;type:varchar(32)" db:"district_name" json:"district_name"`
	DetailAddress    string    `gorm:"column:detail_address;type:varchar(256)" db:"detail_address" json:"detail_address"`
	Longitude        float64   `gorm:"column:longitude;type:decimal(12,6)" db:"longitude" json:"longitude"`
	Latitude         float64   `gorm:"column:latitude;type:decimal(12,6)" db:"latitude" json:"latitude"`
	BusinessStartTime string   `gorm:"column:business_start_time;type:time" db:"business_start_time" json:"business_start_time"`
	BusinessEndTime  string    `gorm:"column:business_end_time;type:time" db:"business_end_time" json:"business_end_time"`
	ContactPhoneMask string    `gorm:"column:contact_phone_mask;type:varchar(32)" db:"contact_phone_mask" json:"contact_phone_mask"`
	Status           int32     `gorm:"column:status;type:tinyint" db:"status" json:"status"`
	ExtraFee         float64   `gorm:"column:extra_fee;type:decimal(10,2)" db:"extra_fee" json:"extra_fee"`
	CreatedAt        time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
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

// ExpressOrder 对应 sf_express_order 表，保存完整寄件订单生命周期数据。
type ExpressOrder struct {
	ID                      int64     `gorm:"primaryKey;column:id" db:"id" json:"id"`
	UserID                  int64     `gorm:"column:user_id;index" db:"user_id" json:"user_id"`
	OrderNo                 string    `gorm:"column:order_no;type:varchar(64);uniqueIndex" db:"order_no" json:"order_no"`
	ExpressNo               string    `gorm:"column:express_no;type:varchar(64)" db:"express_no" json:"express_no"`

	// 寄件人
	SenderName              string    `gorm:"column:sender_name;type:varchar(32)" db:"sender_name" json:"sender_name"`
	SenderPhone             string    `gorm:"column:sender_phone;type:varchar(11)" db:"sender_phone" json:"sender_phone"`
	SenderProvince          string    `gorm:"column:sender_province;type:varchar(32)" db:"sender_province" json:"sender_province"`
	SenderCity              string    `gorm:"column:sender_city;type:varchar(32)" db:"sender_city" json:"sender_city"`
	SenderDistrict          string    `gorm:"column:sender_district;type:varchar(32)" db:"sender_district" json:"sender_district"`
	SenderDetail            string    `gorm:"column:sender_detail;type:varchar(256)" db:"sender_detail" json:"sender_detail"`
	SenderLat               float64   `gorm:"column:sender_lat;type:decimal(10,6)" db:"sender_lat" json:"sender_lat"`
	SenderLng               float64   `gorm:"column:sender_lng;type:decimal(10,6)" db:"sender_lng" json:"sender_lng"`

	// 收件人
	ReceiverName            string    `gorm:"column:receiver_name;type:varchar(32)" db:"receiver_name" json:"receiver_name"`
	ReceiverPhone           string    `gorm:"column:receiver_phone;type:varchar(11)" db:"receiver_phone" json:"receiver_phone"`
	ReceiverProvince        string    `gorm:"column:receiver_province;type:varchar(32)" db:"receiver_province" json:"receiver_province"`
	ReceiverCity            string    `gorm:"column:receiver_city;type:varchar(32)" db:"receiver_city" json:"receiver_city"`
	ReceiverDistrict        string    `gorm:"column:receiver_district;type:varchar(32)" db:"receiver_district" json:"receiver_district"`
	ReceiverDetail          string    `gorm:"column:receiver_detail;type:varchar(256)" db:"receiver_detail" json:"receiver_detail"`
	ReceiverLat             float64   `gorm:"column:receiver_lat;type:decimal(10,6)" db:"receiver_lat" json:"receiver_lat"`
	ReceiverLng             float64   `gorm:"column:receiver_lng;type:decimal(10,6)" db:"receiver_lng" json:"receiver_lng"`

	// 包裹
	Weight                  float64   `gorm:"column:weight;type:decimal(8,2)" db:"weight" json:"weight"`
	Length                  int32     `gorm:"column:length;type:int" db:"length" json:"length"`
	Width                   int32     `gorm:"column:width;type:int" db:"width" json:"width"`
	Height                  int32     `gorm:"column:height;type:int" db:"height" json:"height"`

	// 运费
	BaseFreight             float64   `gorm:"column:base_freight;type:decimal(10,2)" db:"base_freight" json:"base_freight"`
	InsureFee               float64   `gorm:"column:insure_fee;type:decimal(10,2)" db:"insure_fee" json:"insure_fee"`
	TotalFreight            float64   `gorm:"column:total_freight;type:decimal(10,2)" db:"total_freight" json:"total_freight"`

	// 寄件渠道
	ChannelCode             string    `gorm:"column:channel_code;type:varchar(32)" db:"channel_code" json:"channel_code"`
	LockerID                int64     `gorm:"column:locker_id;type:bigint" db:"locker_id" json:"locker_id"`
	BoxType                 string    `gorm:"column:box_type;type:varchar(16)" db:"box_type" json:"box_type"`
	ServicePointID          int64     `gorm:"column:service_point_id;type:bigint" db:"service_point_id" json:"service_point_id"`
	PickupDate              string    `gorm:"column:pickup_date;type:date" db:"pickup_date" json:"pickup_date"`
	PickupSlotCode          string    `gorm:"column:pickup_slot_code;type:varchar(32)" db:"pickup_slot_code" json:"pickup_slot_code"`
	PickupStartTime         string    `gorm:"column:pickup_start_time;type:time" db:"pickup_start_time" json:"pickup_start_time"`
	PickupEndTime           string    `gorm:"column:pickup_end_time;type:time" db:"pickup_end_time" json:"pickup_end_time"`

	// 付款
	PaymentMethod           string    `gorm:"column:payment_method;type:varchar(16)" db:"payment_method" json:"payment_method"`
	PaymentStatus           string    `gorm:"column:payment_status;type:varchar(16);default:pending" db:"payment_status" json:"payment_status"`
	PaidAt                  time.Time `gorm:"column:paid_at" db:"paid_at" json:"paid_at"`

	// 隐私
	PrivacyProtection       bool      `gorm:"column:privacy_protection;type:tinyint;default:0" db:"privacy_protection" json:"privacy_protection"`

	// 状态
	Status                  string    `gorm:"column:status;type:varchar(32);default:pending_pickup" db:"status" json:"status"`
	DelFlag                 bool      `gorm:"column:del_flag;type:tinyint;default:0" db:"del_flag" json:"del_flag"`

	// 标签/关注
	Tag                     string    `gorm:"column:tag;type:varchar(64)" db:"tag" json:"tag"`
	IsFollowed              bool      `gorm:"column:is_followed;type:tinyint;default:0" db:"is_followed" json:"is_followed"`

	// 分享链接
	ShareURL                string    `gorm:"column:share_url;type:varchar(512)" db:"share_url" json:"share_url"`

	CreatedAt               time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
	UpdatedAt               time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (ExpressOrder) TableName() string { return "sf_express_order" }

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
