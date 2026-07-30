// Package data 实现查快递模块的仓储接口，使用 GORM 操作 MySQL。
package data

import (
	"context"
	"time"

	"shunfeng-miniprogram/internal/biz"

	"gorm.io/gorm"
)

// orderFollowRepo 订单关注仓储的 GORM 实现。
type orderFollowRepo struct{}

// NewOrderFollowRepo 创建订单关注仓储，返回 biz 接口。
func NewOrderFollowRepo() biz.OrderFollowRepo {
	return &orderFollowRepo{}
}

// Create 创建关注记录。使用 FirstOrCreate 实现幂等，重复关注不报错。
func (r *orderFollowRepo) Create(ctx context.Context, userID, orderID int64) error {
	po := &OrderFollow{
		UserID:    userID,
		OrderID:   orderID,
		CreatedAt: time.Now(),
	}
	if err := DB.WithContext(ctx).Where("user_id = ? AND order_id = ?", userID, orderID).
		FirstOrCreate(po).Error; err != nil {
		return err
	}
	return nil
}

// Delete 取消关注。逻辑删除（设 deleted_at），物理数据保留。
func (r *orderFollowRepo) Delete(ctx context.Context, userID, orderID int64) error {
	result := DB.WithContext(ctx).Model(&OrderFollow{}).
		Where("user_id = ? AND order_id = ? AND deleted_at IS NULL", userID, orderID).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return biz.ErrExpressOrderNotFound
	}
	return nil
}

// Exists 查询是否已关注（deleted_at IS NULL）。
func (r *orderFollowRepo) Exists(ctx context.Context, userID, orderID int64) (bool, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&OrderFollow{}).
		Where("user_id = ? AND order_id = ? AND deleted_at IS NULL", userID, orderID).
		Count(&count).Error
	return count > 0, err
}

// ListByUser 查询用户关注的所有订单 ID 列表。
func (r *orderFollowRepo) ListByUser(ctx context.Context, userID int64) ([]int64, error) {
	var pos []OrderFollow
	if err := DB.WithContext(ctx).Model(&OrderFollow{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Find(&pos).Error; err != nil {
		return nil, err
	}
	ids := make([]int64, len(pos))
	for i, po := range pos {
		ids[i] = po.OrderID
	}
	return ids, nil
}

// orderLabelRepo 订单标签仓储的 GORM 实现。
type orderLabelRepo struct{}

// NewOrderLabelRepo 创建订单标签仓储，返回 biz 接口。
func NewOrderLabelRepo() biz.OrderLabelRepo {
	return &orderLabelRepo{}
}

// Upsert 覆盖式写入标签。使用 GORM Assign + FirstOrCreate 实现 upsert。
func (r *orderLabelRepo) Upsert(ctx context.Context, userID, orderID int64, content string) (*biz.OrderLabel, error) {
	now := time.Now()
	po := &OrderLabel{
		UserID:    userID,
		OrderID:   orderID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := DB.WithContext(ctx).Where("user_id = ? AND order_id = ?", userID, orderID).
		Assign(map[string]interface{}{
			"content":    content,
			"updated_at": now,
			"deleted_at": nil,
		}).FirstOrCreate(po).Error; err != nil {
		return nil, err
	}
	return toOrderLabelBiz(po), nil
}

// Delete 清空标签（逻辑删除）。
func (r *orderLabelRepo) Delete(ctx context.Context, userID, orderID int64) error {
	result := DB.WithContext(ctx).Model(&OrderLabel{}).
		Where("user_id = ? AND order_id = ? AND deleted_at IS NULL", userID, orderID).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

// Get 查询标签。无标签时返回 (nil, nil)，由 biz 层判断返回空内容。
func (r *orderLabelRepo) Get(ctx context.Context, userID, orderID int64) (*biz.OrderLabel, error) {
	var po OrderLabel
	if err := DB.WithContext(ctx).Where("user_id = ? AND order_id = ? AND deleted_at IS NULL", userID, orderID).
		First(&po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toOrderLabelBiz(&po), nil
}

// BatchGetByOrderIDs 批量查询多个订单的标签，返回 orderID -> content 映射。
func (r *orderLabelRepo) BatchGetByOrderIDs(ctx context.Context, userID int64, orderIDs []int64) (map[int64]string, error) {
	if len(orderIDs) == 0 {
		return nil, nil
	}
	var pos []OrderLabel
	if err := DB.WithContext(ctx).Where("user_id = ? AND order_id IN ? AND deleted_at IS NULL", userID, orderIDs).
		Find(&pos).Error; err != nil {
		return nil, err
	}
	result := make(map[int64]string, len(pos))
	for _, po := range pos {
		result[po.OrderID] = po.Content
	}
	return result, nil
}

// orderShareRepo 订单分享仓储的 GORM 实现。
type orderShareRepo struct{}

// NewOrderShareRepo 创建订单分享仓储，返回 biz 接口。
func NewOrderShareRepo() biz.OrderShareRepo {
	return &orderShareRepo{}
}

// Create 创建分享记录。
func (r *orderShareRepo) Create(ctx context.Context, share *biz.OrderShare) (*biz.OrderShare, error) {
	po := &OrderShare{
		ShareCode:    share.ShareCode,
		UserID:       share.UserID,
		OrderID:      share.OrderID,
		ShowSender:   share.ShowSender,
		ShowReceiver: share.ShowReceiver,
		ShowPhone:    share.ShowPhone,
		ShowStatus:   share.ShowStatus,
		Status:       1,
		ExpiresAt:    share.ExpiresAt,
		CreatedAt:    time.Now(),
	}
	if err := DB.WithContext(ctx).Create(po).Error; err != nil {
		return nil, err
	}
	return toOrderShareBiz(po), nil
}

// GetByCode 按分享码查询。不存在时返回 ErrExpressShareExpired。
func (r *orderShareRepo) GetByCode(ctx context.Context, code string) (*biz.OrderShare, error) {
	var po OrderShare
	if err := DB.WithContext(ctx).Where("share_code = ?", code).First(&po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, biz.ErrExpressShareExpired
		}
		return nil, err
	}
	return toOrderShareBiz(&po), nil
}

// userMessageRepo 用户消息仓储的 GORM 实现。
type userMessageRepo struct{}

// NewUserMessageRepo 创建用户消息仓储，返回 biz 接口。
func NewUserMessageRepo() biz.UserMessageRepo {
	return &userMessageRepo{}
}

// Create 创建消息。重复回调时 unique index 防重复写入。
func (r *userMessageRepo) Create(ctx context.Context, msg *biz.UserMessage) (*biz.UserMessage, error) {
	po := &UserMessage{
		UserID:       msg.UserID,
		MessageType:  msg.MessageType,
		Title:        msg.Title,
		Content:      msg.Content,
		BusinessType: msg.BusinessType,
		BusinessID:   msg.BusinessID,
		Priority:     msg.Priority,
		IsRead:       false,
		CreatedAt:    time.Now(),
	}
	if err := DB.WithContext(ctx).Create(po).Error; err != nil {
		return nil, err
	}
	return toUserMessageBiz(po), nil
}

// List 分页查询消息列表。支持按 message_type 和 read_status 筛选。
// 返回消息列表、总行数、未读数。
func (r *userMessageRepo) List(ctx context.Context, userID int64, msgType, readStatus string, offset, limit int) ([]*biz.UserMessage, int32, int32, error) {
	tx := DB.WithContext(ctx).Model(&UserMessage{}).Where("user_id = ? AND deleted_at IS NULL", userID)

	if msgType != "" {
		tx = tx.Where("message_type = ?", msgType)
	}
	if readStatus == "read" {
		tx = tx.Where("is_read = 1")
	} else if readStatus == "unread" {
		tx = tx.Where("is_read = 0")
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	var pos []UserMessage
	if err := tx.Order("created_at DESC").Offset(offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, 0, 0, err
	}

	var unreadCount int64
	DB.WithContext(ctx).Model(&UserMessage{}).
		Where("user_id = ? AND is_read = 0 AND deleted_at IS NULL", userID).
		Count(&unreadCount)

	items := make([]*biz.UserMessage, len(pos))
	for i, po := range pos {
		items[i] = toUserMessageBiz(&po)
	}
	return items, int32(total), int32(unreadCount), nil
}

// MarkRead 单条标已读。校验 msgID + userID 归属，返回剩余未读数。
func (r *userMessageRepo) MarkRead(ctx context.Context, userID, msgID int64) (int32, error) {
	now := time.Now()
	result := DB.WithContext(ctx).Model(&UserMessage{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", msgID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		})
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, biz.ErrExpressMsgNotFound
	}
	return r.GetUnreadCount(ctx, userID)
}

// MarkAllRead 批量标已读。返回本次标记的行数。
func (r *userMessageRepo) MarkAllRead(ctx context.Context, userID int64) (int32, error) {
	now := time.Now()
	result := DB.WithContext(ctx).Model(&UserMessage{}).
		Where("user_id = ? AND is_read = 0 AND deleted_at IS NULL", userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		})
	if result.Error != nil {
		return 0, result.Error
	}
	return int32(result.RowsAffected), nil
}

// Delete 逻辑删除消息。校验 user_id 归属，返回剩余未读数。
func (r *userMessageRepo) Delete(ctx context.Context, userID, msgID int64) (int32, error) {
	result := DB.WithContext(ctx).Model(&UserMessage{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", msgID, userID).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, biz.ErrExpressMsgNotFound
	}
	return r.GetUnreadCount(ctx, userID)
}

// GetUnreadCount 查询用户未读消息总数。
func (r *userMessageRepo) GetUnreadCount(ctx context.Context, userID int64) (int32, error) {
	var count int64
	if err := DB.WithContext(ctx).Model(&UserMessage{}).
		Where("user_id = ? AND is_read = 0 AND deleted_at IS NULL", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

// toOrderLabelBiz 将 PO 转换为 biz.OrderLabel。
func toOrderLabelBiz(in *OrderLabel) *biz.OrderLabel {
	if in == nil {
		return nil
	}
	return &biz.OrderLabel{
		ID:        in.ID,
		UserID:    in.UserID,
		OrderID:   in.OrderID,
		Content:   in.Content,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

// toOrderShareBiz 将 PO 转换为 biz.OrderShare。
func toOrderShareBiz(in *OrderShare) *biz.OrderShare {
	if in == nil {
		return nil
	}
	return &biz.OrderShare{
		ID:           in.ID,
		ShareCode:    in.ShareCode,
		UserID:       in.UserID,
		OrderID:      in.OrderID,
		ShowSender:   in.ShowSender,
		ShowReceiver: in.ShowReceiver,
		ShowPhone:    in.ShowPhone,
		ShowStatus:   in.ShowStatus,
		Status:       in.Status,
		ExpiresAt:    in.ExpiresAt,
		CreatedAt:    in.CreatedAt,
	}
}

// toUserMessageBiz 将 PO 转换为 biz.UserMessage。
func toUserMessageBiz(in *UserMessage) *biz.UserMessage {
	if in == nil {
		return nil
	}
	return &biz.UserMessage{
		ID:           in.ID,
		UserID:       in.UserID,
		MessageType:  in.MessageType,
		Title:        in.Title,
		Content:      in.Content,
		BusinessType: in.BusinessType,
		BusinessID:   in.BusinessID,
		Priority:     in.Priority,
		IsRead:       in.IsRead,
		ReadAt:       in.ReadAt,
		CreatedAt:    in.CreatedAt,
	}
}


