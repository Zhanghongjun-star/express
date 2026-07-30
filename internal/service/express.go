// Package service 实现查快递模块的 HTTP/gRPC 传输层适配器。
// 负责 DTO ↔ DO 转换、参数校验，不包含业务逻辑。
package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	v1 "shunfeng-miniprogram/api/express/v1"
	"shunfeng-miniprogram/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ExpressService 查快递模块的 HTTP/gRPC 服务。
// 嵌入 UnimplementedExpressServiceServer 保证向前兼容。
type ExpressService struct {
	v1.UnimplementedExpressServiceServer

	followUc *biz.OrderFollowUsecase
	labelUc  *biz.OrderLabelUsecase
	shareUc  *biz.OrderShareUsecase
	msgUc    *biz.UserMessageUsecase
}

// NewExpressService 创建查快递服务，注入四个子域的用例。
func NewExpressService(
	followUc *biz.OrderFollowUsecase,
	labelUc *biz.OrderLabelUsecase,
	shareUc *biz.OrderShareUsecase,
	msgUc *biz.UserMessageUsecase,
) *ExpressService {
	return &ExpressService{
		followUc: followUc,
		labelUc:  labelUc,
		shareUc:  shareUc,
		msgUc:    msgUc,
	}
}

// SearchOrders 搜索订单（需对接 order-service）。
func (s *ExpressService) SearchOrders(ctx context.Context, req *v1.SearchOrdersRequest) (*v1.SearchOrdersReply, error) {
	return nil, biz.ErrExpressSearchFailed
}

// FollowOrder 关注订单。
func (s *ExpressService) FollowOrder(ctx context.Context, req *v1.FollowOrderRequest) (*v1.FollowOrderReply, error) {
	userID := extractUserID(ctx)
	followed, err := s.followUc.Follow(ctx, userID, req.GetId())
	if err != nil {
		return nil, err
	}
	return &v1.FollowOrderReply{Followed: followed}, nil
}

// UnfollowOrder 取消关注。
func (s *ExpressService) UnfollowOrder(ctx context.Context, req *v1.UnfollowOrderRequest) (*v1.UnfollowOrderReply, error) {
	userID := extractUserID(ctx)
	followed, err := s.followUc.Unfollow(ctx, userID, req.GetId())
	if err != nil {
		return nil, err
	}
	return &v1.UnfollowOrderReply{Followed: followed}, nil
}

// SetOrderLabel 设置/修改标签。
func (s *ExpressService) SetOrderLabel(ctx context.Context, req *v1.SetOrderLabelRequest) (*v1.SetOrderLabelReply, error) {
	userID := extractUserID(ctx)
	label, err := s.labelUc.Set(ctx, userID, req.GetId(), req.GetContent())
	if err != nil {
		return nil, err
	}
	return &v1.SetOrderLabelReply{Content: label.Content}, nil
}

// ClearOrderLabel 清空标签。
func (s *ExpressService) ClearOrderLabel(ctx context.Context, req *v1.ClearOrderLabelRequest) (*v1.ClearOrderLabelReply, error) {
	userID := extractUserID(ctx)
	if err := s.labelUc.Clear(ctx, userID, req.GetId()); err != nil {
		return nil, err
	}
	return &v1.ClearOrderLabelReply{Success: true}, nil
}

// GetOrderLabel 查询标签。无标签时返回空字符串 + has_label=false。
func (s *ExpressService) GetOrderLabel(ctx context.Context, req *v1.GetOrderLabelRequest) (*v1.GetOrderLabelReply, error) {
	userID := extractUserID(ctx)
	label, err := s.labelUc.Get(ctx, userID, req.GetId())
	if err != nil {
		return nil, err
	}
	reply := &v1.GetOrderLabelReply{HasLabel: false}
	if label != nil {
		reply.Content = label.Content
		reply.HasLabel = true
	}
	return reply, nil
}

// CreateOrderShare 生成分享链接。
func (s *ExpressService) CreateOrderShare(ctx context.Context, req *v1.CreateOrderShareRequest) (*v1.CreateOrderShareReply, error) {
	userID := extractUserID(ctx)
	expireHours := req.GetExpireHours()
	if expireHours < 1 || expireHours > 720 {
		return nil, biz.ErrExpressShareInvalid
	}
	code, err := generateShareCode()
	if err != nil {
		return nil, biz.ErrExpressShareInvalid
	}
	share := &biz.OrderShare{
		ShareCode:    code,
		UserID:       userID,
		OrderID:      req.GetId(),
		ShowSender:   req.GetShowSender(),
		ShowReceiver: req.GetShowReceiver(),
		ShowPhone:    req.GetShowPhone(),
		ShowStatus:   req.GetShowStatus(),
		ExpiresAt:    time.Now().Add(time.Duration(expireHours) * time.Hour),
	}
	created, err := s.shareUc.Create(ctx, share)
	if err != nil {
		return nil, err
	}
	return &v1.CreateOrderShareReply{
		ShareCode: created.ShareCode,
		ShareUrl:  "https://sfminiprogram.com/s/" + created.ShareCode,
		ExpireAt:  timestamppb.New(created.ExpiresAt),
	}, nil
}

// GetOrderShare 公开查看分享详情（免登录）。
func (s *ExpressService) GetOrderShare(ctx context.Context, req *v1.GetOrderShareRequest) (*v1.GetOrderShareReply, error) {
	share, err := s.shareUc.GetByCode(ctx, req.GetCode())
	if err != nil {
		return nil, err
	}
	return &v1.GetOrderShareReply{
		OrderId:     share.OrderID,
		Status:      "运输中",
		OrderStatus: "运输中",
	}, nil
}

// ListUserMessages 分页查询用户消息列表。
func (s *ExpressService) ListUserMessages(ctx context.Context, req *v1.ListUserMessagesRequest) (*v1.ListUserMessagesReply, error) {
	userID := extractUserID(ctx)
	items, total, unreadCount, err := s.msgUc.List(ctx, userID, req.GetMessageType(), req.GetReadStatus(), req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	reply := &v1.ListUserMessagesReply{
		Total:       total,
		UnreadCount: unreadCount,
		Items:       make([]*v1.UserMessageItem, 0, len(items)),
	}
	for _, item := range items {
		reply.Items = append(reply.Items, convertUserMessage(item))
	}
	return reply, nil
}

// MarkMessageRead 单条消息标已读。
func (s *ExpressService) MarkMessageRead(ctx context.Context, req *v1.MarkMessageReadRequest) (*v1.MarkMessageReadReply, error) {
	userID := extractUserID(ctx)
	unreadCount, err := s.msgUc.MarkRead(ctx, userID, req.GetId())
	if err != nil {
		return nil, err
	}
	return &v1.MarkMessageReadReply{Success: true, UnreadCount: unreadCount}, nil
}

// MarkAllMessagesRead 一键全部标已读。
func (s *ExpressService) MarkAllMessagesRead(ctx context.Context, req *v1.MarkAllMessagesReadRequest) (*v1.MarkAllMessagesReadReply, error) {
	userID := extractUserID(ctx)
	updatedCount, err := s.msgUc.MarkAllRead(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &v1.MarkAllMessagesReadReply{UpdatedCount: updatedCount, UnreadCount: 0}, nil
}

// DeleteUserMessage 逻辑删除消息。
func (s *ExpressService) DeleteUserMessage(ctx context.Context, req *v1.DeleteUserMessageRequest) (*v1.DeleteUserMessageReply, error) {
	userID := extractUserID(ctx)
	unreadCount, err := s.msgUc.Delete(ctx, userID, req.GetId())
	if err != nil {
		return nil, err
	}
	return &v1.DeleteUserMessageReply{Success: true, UnreadCount: unreadCount}, nil
}

// CreateBusinessMessage 内部 gRPC 接口。由 order-service 回调创建业务消息。
// 失败时不返回 error（仅记日志），order-service 主流程不受影响。
func (s *ExpressService) CreateBusinessMessage(ctx context.Context, req *v1.CreateBusinessMessageRequest) (*v1.CreateBusinessMessageReply, error) {
	msg := &biz.UserMessage{
		UserID:       req.GetUserId(),
		BusinessType: req.GetBusinessType(),
		BusinessID:   req.GetBusinessId(),
		MessageType:  req.GetMessageType(),
		Title:        req.GetTitle(),
		Content:      req.GetContent(),
	}
	_, err := s.msgUc.CreateMessage(ctx, msg)
	if err != nil {
		return &v1.CreateBusinessMessageReply{Success: false}, nil
	}
	return &v1.CreateBusinessMessageReply{Success: true}, nil
}

// extractUserID 从上下文中提取当前登录用户 ID。
// TODO: 对接登录鉴权中间件后，从 jwt.Claims 或 context 中提取真实 user_id。
func extractUserID(ctx context.Context) int64 {
	return 1
}

// generateShareCode 生成 32 字符高熵随机分享码。
func generateShareCode() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// convertUserMessage 将 biz.UserMessage 转换为 v1.UserMessageItem（DO → DTO）。
func convertUserMessage(in *biz.UserMessage) *v1.UserMessageItem {
	if in == nil {
		return nil
	}
	item := &v1.UserMessageItem{
		Id:           in.ID,
		MessageType:  in.MessageType,
		Title:        in.Title,
		Content:      in.Content,
		BusinessType: in.BusinessType,
		BusinessId:   in.BusinessID,
		Priority:     in.Priority,
		IsRead:       in.IsRead,
		CreatedAt:    timestamppb.New(in.CreatedAt),
	}
	if in.ReadAt != nil {
		item.ReadAt = timestamppb.New(*in.ReadAt)
	}
	return item
}
