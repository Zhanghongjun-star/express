package data

import (
	"context"
	"fmt"
	"time"

	"shunfeng-miniprogram/internal/biz"

	"gorm.io/gorm"
)

// orderRepo 订单仓储的 GORM 实现。
type orderRepo struct{}

// NewOrderRepo 创建订单仓储。
func NewOrderRepo() biz.OrderRepo {
	return &orderRepo{}
}

func (r *orderRepo) GetHistoryCount(ctx context.Context, userID int64) (*biz.HistoryCount, error) {
	var po User
	if err := DB.WithContext(ctx).First(&po, userID).Error; err != nil {
		return nil, biz.ErrUserNotFound
	}
	return &biz.HistoryCount{
		QueryCountToday: po.QueryCountToday,
		MaxQueryPerDay:  biz.MaxQueryPerDay,
	}, nil
}

func (r *orderRepo) ListHistoryRecords(ctx context.Context, userID int64, offset, limit int) ([]*biz.HistoryRecord, int32, error) {
	var userPo User
	if err := DB.WithContext(ctx).First(&userPo, userID).Error; err != nil {
		return nil, 0, biz.ErrUserNotFound
	}
	now := time.Now()
	if userPo.UpdateTime.Year() != now.Year() || userPo.UpdateTime.YearDay() != now.YearDay() {
		DB.WithContext(ctx).Model(&User{}).Where("user_id = ?", userID).
			Select("query_count_today", "update_time").
			Updates(map[string]any{"query_count_today": 0, "update_time": now})
	}

	// 增加查询次数
	DB.WithContext(ctx).Model(&User{}).Where("user_id = ?", userID).
		UpdateColumn("query_count_today", gorm.Expr("query_count_today + 1"))

	cutoff := time.Now().AddDate(0, 0, -biz.MaxHistoryDays)
	var total int64
	DB.WithContext(ctx).Model(&Order{}).
		Where("user_id = ? AND create_time >= ?", userID, cutoff).
		Count(&total)

	var pos []Order
	if err := DB.WithContext(ctx).
		Where("user_id = ? AND create_time >= ?", userID, cutoff).
		Order("id ASC").Offset(offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	records := make([]*biz.HistoryRecord, len(pos))
	for i, po := range pos {
		records[i] = toOrderBiz(&po)
	}
	return records, int32(total), nil
}

func (r *orderRepo) ExportHistory(ctx context.Context, userID int64, recordIDs []int64) (string, error) {
	var count int64
	DB.WithContext(ctx).Model(&Order{}).
		Where("id IN ? AND user_id = ?", recordIDs, userID).
		Count(&count)
	if int(count) != len(recordIDs) {
		return "", biz.ErrOrderNotFound
	}
	url := fmt.Sprintf("/api/v1/export/download?user_id=%d&t=%d", userID, time.Now().UnixMilli())
	return url, nil
}

// toOrderBiz 将 Order 转换为 biz.HistoryRecord（PO → DO）。
func toOrderBiz(in *Order) *biz.HistoryRecord {
	if in == nil {
		return nil
	}
	return &biz.HistoryRecord{
		ID:           in.ID,
		UserID:       in.UserID,
		OrderNo:      in.OrderNo,
		ExpressNo:    in.ExpressNo,
		SenderName:   in.SenderName,
		ReceiverName: in.ReceiverName,
		Status:       in.Status,
		CreateTime:   in.CreateTime,
	}
}
