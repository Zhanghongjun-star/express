package model

import "time"

// IdentityUser 对应 identity_user 表，保存用户登录账号、角色和账号状态字段。
type IdentityUser struct {
	ID            int64      `db:"id" json:"id"`
	UserNo        string     `db:"user_no" json:"user_no"`
	PhoneCipher   string     `db:"phone_cipher" json:"phone_cipher"`
	PhoneHash     string     `db:"phone_hash" json:"phone_hash"`
	Email         string     `db:"email" json:"email"`
	PasswordHash  string     `db:"password_hash" json:"password_hash"`
	RoleCode      string     `db:"role_code" json:"role_code"`
	AccountStatus string     `db:"account_status" json:"account_status"`
	LockedUntil   *time.Time `db:"locked_until" json:"locked_until"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

// IdentityRealNameAuth 对应 identity_real_name_auth 表，保存实名认证审核资料。
type IdentityRealNameAuth struct {
	ID             int64     `db:"id" json:"id"`
	UserID         int64     `db:"user_id" json:"user_id"`
	RealNameCipher string    `db:"real_name_cipher" json:"real_name_cipher"`
	IDCardCipher   string    `db:"id_card_cipher" json:"id_card_cipher"`
	IDCardHash     string    `db:"id_card_hash" json:"id_card_hash"`
	ImageURLs      string    `db:"image_urls" json:"image_urls"`
	AuthStatus     string    `db:"auth_status" json:"auth_status"`
	RejectReason   string    `db:"reject_reason" json:"reject_reason"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// IdentitySecurityLog 对应 identity_security_log 表，保存登录、改密、退出等安全事件。
type IdentitySecurityLog struct {
	ID            int64     `db:"id" json:"id"`
	UserID        int64     `db:"user_id" json:"user_id"`
	EventType     string    `db:"event_type" json:"event_type"`
	Result        string    `db:"result" json:"result"`
	IP            string    `db:"ip" json:"ip"`
	DeviceID      string    `db:"device_id" json:"device_id"`
	FailureReason string    `db:"failure_reason" json:"failure_reason"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
