package biz

import (
	"context"
	"time"

	v1 "shunfeng-miniprogram/api/order/v1"

	"github.com/go-kratos/kratos/v3/errors"
)

var (
	// ErrOrderNotFound 订单/记录不存在时返回。
	ErrOrderNotFound = errors.NotFound(v1.ErrorReason_ORDER_NOT_FOUND.String(), "order not found")
	// ErrOrderInvalidArgument 请求参数无效时返回。
	ErrOrderInvalidArgument = errors.BadRequest(v1.ErrorReason_ORDER_INVALID_ARGUMENT.String(), "invalid order argument")
	// ErrOrderQueryLimitExceeded 当日查询次数超限时返回。
	ErrOrderQueryLimitExceeded = errors.Forbidden(v1.ErrorReason_ORDER_QUERY_LIMIT_EXCEEDED.String(), "query limit exceeded")
)

const (
	// MaxQueryPerDay 每日最大查询次数。
	MaxQueryPerDay = 3
	// MaxHistoryDays 历史记录保留天数。
	MaxHistoryDays = 90
	// DefaultPageSize 默认分页大小。
	DefaultPageSize = 20
)

// HistoryCount 当日查询次数统计。
type HistoryCount struct {
	QueryCountToday int32
	MaxQueryPerDay  int32
}

// HistoryRecord 历史快件记录领域对象。
type HistoryRecord struct {
	ID           int64
	UserID       int64
	OrderNo      string
	ExpressNo    string
	SenderName   string
	ReceiverName string
	Status       string
	CreateTime   time.Time
}

// OrderRepo 订单仓储接口。
type OrderRepo interface {
	GetHistoryCount(context.Context, int64) (*HistoryCount, error)
	ListHistoryRecords(context.Context, int64, int, int) ([]*HistoryRecord, int32, error)
	ExportHistory(context.Context, int64, []int64) (string, error)
}

// OrderUsecase 订单用例。
type OrderUsecase struct {
	repo OrderRepo
}

// NewOrderUsecase 创建订单用例。
func NewOrderUsecase(repo OrderRepo) *OrderUsecase {
	return &OrderUsecase{repo: repo}
}

// GetHistoryCount 获取当日查询次数。
func (uc *OrderUsecase) GetHistoryCount(ctx context.Context, userID int64) (*HistoryCount, error) {
	if userID <= 0 {
		return nil, ErrOrderInvalidArgument
	}
	return uc.repo.GetHistoryCount(ctx, userID)
}

// ListHistoryRecords 获取 90 天内历史快件记录。
func (uc *OrderUsecase) ListHistoryRecords(ctx context.Context, userID int64, offset, limit int) ([]*HistoryRecord, int32, error) {
	if userID <= 0 || limit <= 0 {
		return nil, 0, ErrOrderInvalidArgument
	}
	if limit > DefaultPageSize {
		limit = DefaultPageSize
	}
	count, err := uc.repo.GetHistoryCount(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	if count.QueryCountToday >= MaxQueryPerDay {
		return nil, 0, ErrOrderQueryLimitExceeded
	}
	records, total, err := uc.repo.ListHistoryRecords(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

// ExportHistory 导出物流单据。
func (uc *OrderUsecase) ExportHistory(ctx context.Context, userID int64, recordIDs []int64) (string, error) {
	if userID <= 0 || len(recordIDs) == 0 {
		return "", ErrOrderInvalidArgument
	}
	return uc.repo.ExportHistory(ctx, userID, recordIDs)
}
