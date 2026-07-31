package service

import (
	"context"

	v1 "shunfeng-miniprogram/api/order/v1"
	"shunfeng-miniprogram/internal/biz"

	"go.einride.tech/aip/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const defaultOrderPageSize = 20

// OrderService 订单服务。
type OrderService struct {
	v1.UnimplementedOrderServiceServer

	uc *biz.OrderUsecase
}

// NewOrderService 创建订单服务。
func NewOrderService(uc *biz.OrderUsecase) *OrderService {
	return &OrderService{uc: uc}
}

// GetHistoryCount 获取当日查询次数。
func (s *OrderService) GetHistoryCount(ctx context.Context, req *v1.GetHistoryCountRequest) (*v1.HistoryCount, error) {
	count, err := s.uc.GetHistoryCount(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &v1.HistoryCount{
		QueryCountToday: count.QueryCountToday,
		MaxQueryPerDay:  count.MaxQueryPerDay,
	}, nil
}

// ListHistoryRecords 获取 90 天内历史快件记录。
func (s *OrderService) ListHistoryRecords(ctx context.Context, req *v1.ListHistoryRecordsRequest) (*v1.HistoryRecordSet, error) {
	pageToken, err := pagination.ParsePageToken(req)
	if err != nil {
		return nil, err
	}
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = defaultOrderPageSize
	}
	records, total, err := s.uc.ListHistoryRecords(ctx, req.GetUserId(), int(pageToken.Offset), pageSize)
	if err != nil {
		return nil, err
	}
	set := &v1.HistoryRecordSet{
		Records:    make([]*v1.HistoryRecord, 0, len(records)),
		TotalCount: total,
	}
	if len(records) >= pageSize {
		set.NextPageToken = pageToken.Next(req).String()
	}
	for _, rec := range records {
		set.Records = append(set.Records, convertHistoryRecordReply(rec))
	}
	return set, nil
}

// ExportHistory 导出物流单据。
func (s *OrderService) ExportHistory(ctx context.Context, req *v1.ExportHistoryRequest) (*v1.ExportHistoryReply, error) {
	url, err := s.uc.ExportHistory(ctx, req.GetUserId(), req.GetRecordIds())
	if err != nil {
		return nil, err
	}
	return &v1.ExportHistoryReply{DownloadUrl: url}, nil
}

// convertHistoryRecordReply 将 biz.HistoryRecord 转换为 v1.HistoryRecord。
func convertHistoryRecordReply(in *biz.HistoryRecord) *v1.HistoryRecord {
	if in == nil {
		return nil
	}
	return &v1.HistoryRecord{
		Id:           in.ID,
		OrderNo:      in.OrderNo,
		ExpressNo:    in.ExpressNo,
		SenderName:   in.SenderName,
		ReceiverName: in.ReceiverName,
		Status:       in.Status,
		CreateTime:   timestamppb.New(in.CreateTime),
	}
}
