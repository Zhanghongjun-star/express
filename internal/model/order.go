package model

import "time"

// Order 对应 sf_order_history 表，保存用户历史快件记录。
type Order struct {
	ID           int64     `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	OrderNo      string    `db:"order_no" json:"order_no"`
	ExpressNo    string    `db:"express_no" json:"express_no"`
	SenderName   string    `db:"sender_name" json:"sender_name"`
	ReceiverName string    `db:"receiver_name" json:"receiver_name"`
	Status       string    `db:"status" json:"status"`
	CreateTime   time.Time `db:"create_time" json:"create_time"`
	UpdateTime   time.Time `db:"update_time" json:"update_time"`
}
