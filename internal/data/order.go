package data

import (
	"context"
	"fmt"
	"sync"
	"time"

	"shunfeng-miniprogram/internal/biz"
	"shunfeng-miniprogram/internal/model"
)

// orderRepo 订单仓储的内存实现。
type orderRepo struct {
	data *Data

	mu      sync.RWMutex
	nextID  int64
	records map[int64]*model.Order
	counts  map[int64]int32
}

// NewOrderRepo 创建订单仓储。
func NewOrderRepo(data *Data) biz.OrderRepo {
	now := time.Now()
	r := &orderRepo{
		data:    data,
		nextID:  1,
		records: make(map[int64]*model.Order),
		counts:  make(map[int64]int32),
	}
	now = now.Add(-24 * time.Hour)
	for i := int64(1); i <= 3; i++ {
		r.records[i] = &model.Order{
			ID:           i,
			UserID:       1,
			OrderNo:      fmt.Sprintf("SF%d", 20260728001+i),
			ExpressNo:    fmt.Sprintf("SF%d", 1000000000+i),
			SenderName:   "张三",
			ReceiverName: "李四",
			Status:       "已签收",
			CreateTime:   now.Add(time.Duration(i) * time.Hour),
			UpdateTime:   time.Now(),
		}
		r.nextID++
	}
	r.counts[1] = 0
	return r
}

func (r *orderRepo) GetHistoryCount(_ context.Context, userID int64) (*biz.HistoryCount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := r.counts[userID]
	return &biz.HistoryCount{
		QueryCountToday: count,
		MaxQueryPerDay:  biz.MaxQueryPerDay,
	}, nil
}

func (r *orderRepo) ListHistoryRecords(_ context.Context, userID int64, offset, limit int) ([]*biz.HistoryRecord, int32, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counts[userID]++

	cutoff := time.Now().AddDate(0, 0, -biz.MaxHistoryDays)
	filtered := make([]*model.Order, 0)
	for _, rec := range r.records {
		if rec.UserID != userID {
			continue
		}
		if rec.CreateTime.Before(cutoff) {
			continue
		}
		filtered = append(filtered, rec)
	}

	total := int32(len(filtered))
	if offset >= len(filtered) {
		return []*biz.HistoryRecord{}, total, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	result := make([]*biz.HistoryRecord, 0, end-offset)
	for _, rec := range filtered[offset:end] {
		result = append(result, toRecordBiz(rec))
	}
	return result, total, nil
}

func (r *orderRepo) ExportHistory(_ context.Context, userID int64, recordIDs []int64) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, id := range recordIDs {
		rec, ok := r.records[id]
		if !ok || rec.UserID != userID {
			return "", biz.ErrOrderNotFound
		}
	}
	url := fmt.Sprintf("https://example.com/export/%d_%d.xlsx", userID, time.Now().UnixMilli())
	return url, nil
}

// newRecord 将 biz.HistoryRecord 转换为 model.Order（DO → PO）。
func newRecord(in *biz.HistoryRecord) *model.Order {
	if in == nil {
		return nil
	}
	return &model.Order{
		ID:           in.ID,
		UserID:       in.UserID,
		OrderNo:      in.OrderNo,
		ExpressNo:    in.ExpressNo,
		SenderName:   in.SenderName,
		ReceiverName: in.ReceiverName,
		Status:       in.Status,
		CreateTime:   in.CreateTime,
		UpdateTime:   time.Now(),
	}
}

// toRecordBiz 将 model.Order 转换为 biz.HistoryRecord（PO → DO）。
func toRecordBiz(in *model.Order) *biz.HistoryRecord {
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
