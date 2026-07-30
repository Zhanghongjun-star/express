// Package biz 包含查快递模块的领域对象、仓储接口和用例。
// 关注·标签·分享·消息 四个子域各自独立，通过仓储接口与 data 层解耦。
package biz

import (
	"context"
	"time"

	v1 "shunfeng-miniprogram/api/express/v1"

	"github.com/go-kratos/kratos/v3/errors"
)

// 查快递模块业务错误码。
var (
	// ErrExpressOrderNotFound 订单不存在或无权操作。
	ErrExpressOrderNotFound = errors.NotFound(v1.ErrorReason_EXPRESS_ORDER_NOT_FOUND.String(), "order not found or no permission")
	// ErrExpressLabelInvalid 标签内容长度不合法（1-100字符）。
	ErrExpressLabelInvalid = errors.BadRequest(v1.ErrorReason_EXPRESS_LABEL_INVALID.String(), "label content must be 1-100 characters")
	// ErrExpressShareInvalid 分享参数校验不通过。
	ErrExpressShareInvalid = errors.BadRequest(v1.ErrorReason_EXPRESS_SHARE_INVALID.String(), "invalid share parameters")
	// ErrExpressShareExpired 分享链接已过期或不存在。
	ErrExpressShareExpired = errors.NotFound(v1.ErrorReason_EXPRESS_SHARE_EXPIRED.String(), "share link has expired")
	// ErrExpressShareFreqLimit 分享生成频率超限。
	ErrExpressShareFreqLimit = errors.Forbidden(v1.ErrorReason_EXPRESS_SHARE_FREQ_LIMIT.String(), "too frequent, please try later")
	// ErrExpressMsgNotFound 消息不存在或无权操作。
	ErrExpressMsgNotFound = errors.NotFound(v1.ErrorReason_EXPRESS_MSG_NOT_FOUND.String(), "message not found or no permission")
	// ErrExpressSearchFailed 搜索服务暂不可用（order-service 超时）。
	ErrExpressSearchFailed = errors.InternalServer(v1.ErrorReason_EXPRESS_SEARCH_FAILED.String(), "search service unavailable")
)

// OrderFollow 订单关注领域对象。
// 一个用户对一个订单只能关注一次，通过 uk_user_order 唯一索引保障。
type OrderFollow struct {
	ID        int64
	UserID    int64
	OrderID   int64
	CreatedAt time.Time
}

// OrderLabel 订单标签领域对象。
// 用户给订单设置的自定义标签，仅自己可见，通过覆盖式 upsert 修改。
type OrderLabel struct {
	ID        int64
	UserID    int64
	OrderID   int64
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OrderShare 订单分享领域对象。
// share_code 为 32 字符高熵随机码，公开查看时凭 code 获取脱敏数据。
type OrderShare struct {
	ID           int64
	ShareCode    string
	UserID       int64
	OrderID      int64
	ShowSender   bool
	ShowReceiver bool
	ShowPhone    bool
	ShowStatus   bool
	Status       int32
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

// UserMessage 用户消息领域对象。
// 由 order-service 通过内部 gRPC 回调创建，支持已读/未读状态管理。
type UserMessage struct {
	ID           int64
	UserID       int64
	MessageType  string
	Title        string
	Content      string
	BusinessType string
	BusinessID   string
	Priority     int32
	IsRead       bool
	ReadAt       *time.Time
	CreatedAt    time.Time
}

// SearchResult 搜索结果领域对象，合并了订单基础信息和关注/标签状态。
type SearchResult struct {
	ID           int64
	OrderNo      string
	ExpressNo    string
	SenderName   string
	ReceiverName string
	Status       string
	HasFollowed  bool
	Label        string
	CreateTime   time.Time
}

// OrderFollowRepo 订单关注仓储接口。
type OrderFollowRepo interface {
	// Create 创建关注记录（幂等，重复关注不报错）。
	Create(ctx context.Context, userID, orderID int64) error
	// Delete 取消关注（逻辑删除，设 deleted_at）。
	Delete(ctx context.Context, userID, orderID int64) error
	// Exists 查询是否已关注。
	Exists(ctx context.Context, userID, orderID int64) (bool, error)
	// ListByUser 查询用户关注的所有订单 ID 列表。
	ListByUser(ctx context.Context, userID int64) ([]int64, error)
}

// OrderLabelRepo 订单标签仓储接口。
type OrderLabelRepo interface {
	// Upsert 覆盖式写入标签（唯一索引保障每个订单只有一个标签）。
	Upsert(ctx context.Context, userID, orderID int64, content string) (*OrderLabel, error)
	// Delete 清空标签（逻辑删除）。
	Delete(ctx context.Context, userID, orderID int64) error
	// Get 查询单个订单的标签。
	Get(ctx context.Context, userID, orderID int64) (*OrderLabel, error)
	// BatchGetByOrderIDs 批量查询多个订单的标签，返回 orderID -> content 映射。
	BatchGetByOrderIDs(ctx context.Context, userID int64, orderIDs []int64) (map[int64]string, error)
}

// OrderShareRepo 订单分享仓储接口。
type OrderShareRepo interface {
	// Create 创建分享记录。
	Create(ctx context.Context, share *OrderShare) (*OrderShare, error)
	// GetByCode 按分享码查询分享记录。
	GetByCode(ctx context.Context, code string) (*OrderShare, error)
}

// UserMessageRepo 用户消息仓储接口。
type UserMessageRepo interface {
	// Create 创建消息（幂等，重复回调不重复写入）。
	Create(ctx context.Context, msg *UserMessage) (*UserMessage, error)
	// List 分页查询消息列表，返回消息、总数、未读数。
	List(ctx context.Context, userID int64, msgType, readStatus string, offset, limit int) ([]*UserMessage, int32, int32, error)
	// MarkRead 单条标已读，返回剩余未读数。
	MarkRead(ctx context.Context, userID, msgID int64) (int32, error)
	// MarkAllRead 全部标已读，返回本次标记条数。
	MarkAllRead(ctx context.Context, userID int64) (int32, error)
	// Delete 逻辑删除消息，返回剩余未读数。
	Delete(ctx context.Context, userID, msgID int64) (int32, error)
	// GetUnreadCount 查询用户未读消息总数。
	GetUnreadCount(ctx context.Context, userID int64) (int32, error)
}

// OrderFollowUsecase 订单关注用例。
type OrderFollowUsecase struct {
	followRepo OrderFollowRepo
	labelRepo  OrderLabelRepo
	shareRepo  OrderShareRepo
	msgRepo    UserMessageRepo
}

// NewOrderFollowUsecase 创建订单关注用例。
func NewOrderFollowUsecase(followRepo OrderFollowRepo) *OrderFollowUsecase {
	return &OrderFollowUsecase{followRepo: followRepo}
}

// Follow 关注订单。先校验参数，委托 repo 写入关注表。
func (uc *OrderFollowUsecase) Follow(ctx context.Context, userID, orderID int64) (bool, error) {
	if userID <= 0 || orderID <= 0 {
		return false, ErrExpressOrderNotFound
	}
	if err := uc.followRepo.Create(ctx, userID, orderID); err != nil {
		return false, err
	}
	return true, nil
}

// Unfollow 取消关注。逻辑删除关注记录。
func (uc *OrderFollowUsecase) Unfollow(ctx context.Context, userID, orderID int64) (bool, error) {
	if userID <= 0 || orderID <= 0 {
		return false, ErrExpressOrderNotFound
	}
	if err := uc.followRepo.Delete(ctx, userID, orderID); err != nil {
		return false, err
	}
	return false, nil
}

// OrderLabelUsecase 订单标签用例。
type OrderLabelUsecase struct {
	labelRepo OrderLabelRepo
}

// NewOrderLabelUsecase 创建订单标签用例。
func NewOrderLabelUsecase(labelRepo OrderLabelRepo) *OrderLabelUsecase {
	return &OrderLabelUsecase{labelRepo: labelRepo}
}

// Set 设置/修改标签。校验参数后执行 upsert 覆盖写入。
func (uc *OrderLabelUsecase) Set(ctx context.Context, userID, orderID int64, content string) (*OrderLabel, error) {
	if userID <= 0 || orderID <= 0 {
		return nil, ErrExpressOrderNotFound
	}
	if len(content) < 1 || len(content) > 100 {
		return nil, ErrExpressLabelInvalid
	}
	return uc.labelRepo.Upsert(ctx, userID, orderID, content)
}

// Clear 清空标签。逻辑删除标签记录。
func (uc *OrderLabelUsecase) Clear(ctx context.Context, userID, orderID int64) error {
	if userID <= 0 || orderID <= 0 {
		return ErrExpressOrderNotFound
	}
	return uc.labelRepo.Delete(ctx, userID, orderID)
}

// Get 查询标签。无标签时返回 nil（非 error）。
func (uc *OrderLabelUsecase) Get(ctx context.Context, userID, orderID int64) (*OrderLabel, error) {
	if userID <= 0 || orderID <= 0 {
		return nil, ErrExpressOrderNotFound
	}
	return uc.labelRepo.Get(ctx, userID, orderID)
}

// OrderShareUsecase 订单分享用例。
type OrderShareUsecase struct {
	shareRepo OrderShareRepo
}

// NewOrderShareUsecase 创建订单分享用例。
func NewOrderShareUsecase(shareRepo OrderShareRepo) *OrderShareUsecase {
	return &OrderShareUsecase{shareRepo: shareRepo}
}

// Create 创建分享记录。
func (uc *OrderShareUsecase) Create(ctx context.Context, share *OrderShare) (*OrderShare, error) {
	if share == nil || share.OrderID <= 0 || share.ExpiresAt.IsZero() {
		return nil, ErrExpressShareInvalid
	}
	return uc.shareRepo.Create(ctx, share)
}

// GetByCode 按分享码查看分享。校验状态=启用 且 未过期。
func (uc *OrderShareUsecase) GetByCode(ctx context.Context, code string) (*OrderShare, error) {
	if len(code) == 0 {
		return nil, ErrExpressShareExpired
	}
	share, err := uc.shareRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, ErrExpressShareExpired
	}
	if share.Status != 1 || time.Now().After(share.ExpiresAt) {
		return nil, ErrExpressShareExpired
	}
	return share, nil
}

// UserMessageUsecase 用户消息用例。
type UserMessageUsecase struct {
	msgRepo UserMessageRepo
}

// NewUserMessageUsecase 创建用户消息用例。
func NewUserMessageUsecase(msgRepo UserMessageRepo) *UserMessageUsecase {
	return &UserMessageUsecase{msgRepo: msgRepo}
}

// List 分页查询消息列表。默认第1页20条，最大50条。
func (uc *UserMessageUsecase) List(ctx context.Context, userID int64, msgType, readStatus string, page, pageSize int32) ([]*UserMessage, int32, int32, error) {
	if userID <= 0 {
		return nil, 0, 0, ErrExpressMsgNotFound
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 20
	}
	offset := int((page - 1) * pageSize)
	return uc.msgRepo.List(ctx, userID, msgType, readStatus, offset, int(pageSize))
}

// MarkRead 单条消息标已读。校验 user_id 归属。
func (uc *UserMessageUsecase) MarkRead(ctx context.Context, userID, msgID int64) (int32, error) {
	if userID <= 0 || msgID <= 0 {
		return 0, ErrExpressMsgNotFound
	}
	return uc.msgRepo.MarkRead(ctx, userID, msgID)
}

// MarkAllRead 一键全部已读。
func (uc *UserMessageUsecase) MarkAllRead(ctx context.Context, userID int64) (int32, error) {
	if userID <= 0 {
		return 0, ErrExpressMsgNotFound
	}
	return uc.msgRepo.MarkAllRead(ctx, userID)
}

// Delete 逻辑删除消息。返回剩余未读数，供前端刷新角标。
func (uc *UserMessageUsecase) Delete(ctx context.Context, userID, msgID int64) (int32, error) {
	if userID <= 0 || msgID <= 0 {
		return 0, ErrExpressMsgNotFound
	}
	return uc.msgRepo.Delete(ctx, userID, msgID)
}

// CreateMessage 创建业务消息（由 order-service 内部 gRPC 回调调用）。
func (uc *UserMessageUsecase) CreateMessage(ctx context.Context, msg *UserMessage) (*UserMessage, error) {
	if msg == nil {
		return nil, ErrExpressMsgNotFound
	}
	return uc.msgRepo.Create(ctx, msg)
}
