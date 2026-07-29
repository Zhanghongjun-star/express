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
