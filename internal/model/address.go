package model

import (
	"time"

	v1 "shunfeng-miniprogram/api/address/v1"
)

// Address 对应 sf_user_address 表，保存用户地址簿数据。
type Address struct {
	ID            int64          `db:"id" json:"id"`
	UserID        int64          `db:"user_id" json:"user_id"`
	AddrType      v1.AddressType `db:"addr_type" json:"addr_type"`
	ReceiverName  string         `db:"receiver_name" json:"receiver_name"`
	ReceiverPhone string         `db:"receiver_phone" json:"receiver_phone"`
	Province      string         `db:"province" json:"province"`
	City          string         `db:"city" json:"city"`
	District      string         `db:"district" json:"district"`
	DetailAddr    string         `db:"detail_addr" json:"detail_addr"`
	IsDefault     bool           `db:"is_default" json:"is_default"`
	CreateTime    time.Time      `db:"create_time" json:"create_time"`
	UpdateTime    time.Time      `db:"update_time" json:"update_time"`
	DelFlag       bool           `db:"del_flag" json:"del_flag"`
}
