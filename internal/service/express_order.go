package service

import (
	"context"

	v3 "shunfeng-miniprogram/api/order/v3"
	"shunfeng-miniprogram/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ExpressOrderService struct {
	v3.UnimplementedOrderServiceServer
	uc *biz.ExpressOrderUsecase
}

func NewExpressOrderService(uc *biz.ExpressOrderUsecase) *ExpressOrderService {
	return &ExpressOrderService{uc: uc}
}

func (s *ExpressOrderService) CreateExpressOrder(ctx context.Context, req *v3.CreateExpressOrderRequest) (*v3.ExpressOrderReply, error) {
	order, err := s.uc.CreateOrder(ctx, &biz.CreateExpressOrderCommand{
		UserID:           req.GetUserId(),
		SenderName:       req.GetSenderName(),
		SenderPhone:      req.GetSenderPhone(),
		SenderProvince:   req.GetSenderProvince(),
		SenderCity:       req.GetSenderCity(),
		SenderDistrict:   req.GetSenderDistrict(),
		SenderDetail:     req.GetSenderDetail(),
		ReceiverName:     req.GetReceiverName(),
		ReceiverPhone:    req.GetReceiverPhone(),
		ReceiverProvince: req.GetReceiverProvince(),
		ReceiverCity:     req.GetReceiverCity(),
		ReceiverDistrict: req.GetReceiverDistrict(),
		ReceiverDetail:   req.GetReceiverDetail(),
		Weight:           req.GetWeight(),
		Length:           req.GetLength(),
		Width:            req.GetWidth(),
		Height:           req.GetHeight(),
		ChannelCode:      req.GetChannelCode().String(),
		LockerID:         req.GetLockerId(),
		BoxType:          req.GetBoxType(),
		ServicePointID:   req.GetServicePointId(),
		PickupDate:       req.GetPickupDate(),
		PickupSlotCode:   req.GetPickupSlotCode(),
		PaymentMethod:    convertPaymentMethod(req.GetPaymentMethod()),
		PrivacyProtection: req.GetPrivacyProtection(),
		ExpressType:      req.GetExpressType(),
		InsureValue:      req.GetInsureValue(),
	})
	if err != nil {
		return nil, err
	}
	return convertExpressOrderReply(order), nil
}

func (s *ExpressOrderService) GetExpressOrder(ctx context.Context, req *v3.GetExpressOrderRequest) (*v3.ExpressOrderReply, error) {
	order, err := s.uc.GetOrder(ctx, req.GetUserId(), req.GetId())
	if err != nil {
		return nil, err
	}
	return convertExpressOrderReply(order), nil
}

func (s *ExpressOrderService) SearchOrders(ctx context.Context, req *v3.SearchOrdersRequest) (*v3.ExpressOrderSet, error) {
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 || pageSize > 20 {
		pageSize = 20
	}
	orders, total, err := s.uc.SearchOrders(ctx, req.GetUserId(), req.GetKeyword(), 0, pageSize)
	if err != nil {
		return nil, err
	}
	return &v3.ExpressOrderSet{
		Orders:     convertExpressOrderReplyList(orders),
		TotalCount: total,
	}, nil
}

func (s *ExpressOrderService) ListOrdersByCategory(ctx context.Context, req *v3.ListOrdersByCategoryRequest) (*v3.ExpressOrderSet, error) {
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 || pageSize > 20 {
		pageSize = 20
	}
	orders, total, err := s.uc.ListByCategory(ctx, req.GetUserId(), req.GetCategory(), 0, pageSize)
	if err != nil {
		return nil, err
	}
	return &v3.ExpressOrderSet{
		Orders:     convertExpressOrderReplyList(orders),
		TotalCount: total,
	}, nil
}

func (s *ExpressOrderService) DeleteExpressOrder(ctx context.Context, req *v3.DeleteExpressOrderRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeleteOrder(ctx, req.GetUserId(), req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ExpressOrderService) SetOrderTag(ctx context.Context, req *v3.SetOrderTagRequest) (*v3.ExpressOrderReply, error) {
	order, err := s.uc.SetOrderTag(ctx, req.GetUserId(), req.GetId(), req.GetTag())
	if err != nil {
		return nil, err
	}
	return convertExpressOrderReply(order), nil
}

func (s *ExpressOrderService) FollowOrder(ctx context.Context, req *v3.FollowOrderRequest) (*v3.ExpressOrderReply, error) {
	order, err := s.uc.FollowOrder(ctx, req.GetUserId(), req.GetId(), req.GetFollow())
	if err != nil {
		return nil, err
	}
	return convertExpressOrderReply(order), nil
}

func (s *ExpressOrderService) ShareOrder(ctx context.Context, req *v3.ShareOrderRequest) (*v3.ShareOrderReply, error) {
	url, err := s.uc.ShareOrder(ctx, req.GetUserId(), req.GetId())
	if err != nil {
		return nil, err
	}
	return &v3.ShareOrderReply{ShareUrl: url}, nil
}

func (s *ExpressOrderService) PayExpressOrder(ctx context.Context, req *v3.PayExpressOrderRequest) (*v3.PayExpressOrderReply, error) {
	payURL, tradeNo, err := s.uc.PayOrder(ctx, req.GetUserId(), req.GetId())
	if err != nil {
		return nil, err
	}
	return &v3.PayExpressOrderReply{PayUrl: payURL, TradeNo: tradeNo}, nil
}

func (s *ExpressOrderService) HandlePaymentCallback(ctx context.Context, req *v3.PaymentCallbackRequest) (*v3.PaymentCallbackReply, error) {
	if err := s.uc.HandlePaymentCallback(ctx, req.GetTradeNo(), req.GetOrderNo(), req.GetPaymentStatus(), req.GetSignature()); err != nil {
		return nil, err
	}
	return &v3.PaymentCallbackReply{Success: true}, nil
}

func convertPaymentMethod(pm v3.PaymentMethod) string {
	switch pm {
	case v3.PaymentMethod_SENDER_PAY:
		return biz.PaymentMethodSenderPay
	case v3.PaymentMethod_RECEIVER_PAY:
		return biz.PaymentMethodReceiverPay
	case v3.PaymentMethod_MONTHLY_SETTLE:
		return biz.PaymentMethodMonthly
	default:
		return biz.PaymentMethodSenderPay
	}
}

func convertChannelCode(code string) v3.ChannelCode {
	switch code {
	case "DOOR_PICKUP":
		return v3.ChannelCode_DOOR_PICKUP
	case "LOCKER":
		return v3.ChannelCode_LOCKER
	case "SF_STATION":
		return v3.ChannelCode_SF_STATION
	case "PARTNER_STORE":
		return v3.ChannelCode_PARTNER_STORE
	default:
		return v3.ChannelCode_CHANNEL_UNSPECIFIED
	}
}

func convertOrderStatus(status string) v3.OrderStatus {
	switch status {
	case biz.OrderStatusPendingPickup:
		return v3.OrderStatus_PENDING_PICKUP
	case biz.OrderStatusAccepted:
		return v3.OrderStatus_ACCEPTED
	case biz.OrderStatusAwaitingPickup:
		return v3.OrderStatus_AWAITING_PICKUP
	case biz.OrderStatusPickedUp:
		return v3.OrderStatus_PICKED_UP
	case biz.OrderStatusInTransit:
		return v3.OrderStatus_IN_TRANSIT
	case biz.OrderStatusDelivering:
		return v3.OrderStatus_DELIVERING
	case biz.OrderStatusSigned:
		return v3.OrderStatus_SIGNED
	case biz.OrderStatusCancelled:
		return v3.OrderStatus_CANCELLED
	default:
		return v3.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func convertPaymentStatus(status string) v3.PaymentStatus {
	switch status {
	case biz.PaymentStatusPending:
		return v3.PaymentStatus_PENDING
	case biz.PaymentStatusPaid:
		return v3.PaymentStatus_PAID
	case biz.PaymentStatusExpired:
		return v3.PaymentStatus_EXPIRED
	default:
		return v3.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}

func convertExpressOrderReply(in *biz.ExpressOrder) *v3.ExpressOrderReply {
	if in == nil {
		return nil
	}
	reply := &v3.ExpressOrderReply{
		Id:                in.ID,
		UserId:            in.UserID,
		OrderNo:           in.OrderNo,
		ExpressNo:         in.ExpressNo,
		SenderName:        in.SenderName,
		SenderPhone:       in.SenderPhone,
		SenderProvince:    in.SenderProvince,
		SenderCity:        in.SenderCity,
		SenderDistrict:    in.SenderDistrict,
		SenderDetail:      in.SenderDetail,
		ReceiverName:      in.ReceiverName,
		ReceiverPhone:     in.ReceiverPhone,
		ReceiverProvince:  in.ReceiverProvince,
		ReceiverCity:      in.ReceiverCity,
		ReceiverDistrict:  in.ReceiverDistrict,
		ReceiverDetail:    in.ReceiverDetail,
		Weight:            in.Weight,
		Length:            in.Length,
		Width:             in.Width,
		Height:            in.Height,
		BaseFreight:       in.BaseFreight,
		InsureFee:         in.InsureFee,
		TotalFreight:      in.TotalFreight,
		ChannelCode:       convertChannelCode(in.ChannelCode),
		LockerId:          in.LockerID,
		BoxType:           in.BoxType,
		ServicePointId:    in.ServicePointID,
		PickupDate:        in.PickupDate,
		PickupSlotCode:    in.PickupSlotCode,
		PrivacyProtection: in.PrivacyProtection,
		Status:            convertOrderStatus(in.Status),
		Tag:               in.Tag,
		IsFollowed:        in.IsFollowed,
		ShareUrl:          in.ShareURL,
		PayUrl:            in.PayURL,
		CreatedAt:         timestamppb.New(in.CreatedAt),
		UpdatedAt:         timestamppb.New(in.UpdatedAt),
	}

	// Set payment method
	switch in.PaymentMethod {
	case biz.PaymentMethodSenderPay:
		reply.PaymentMethod = v3.PaymentMethod_SENDER_PAY
	case biz.PaymentMethodReceiverPay:
		reply.PaymentMethod = v3.PaymentMethod_RECEIVER_PAY
	case biz.PaymentMethodMonthly:
		reply.PaymentMethod = v3.PaymentMethod_MONTHLY_SETTLE
	}

	// Set payment status
	switch in.PaymentStatus {
	case biz.PaymentStatusPending:
		reply.PaymentStatus = v3.PaymentStatus_PENDING
	case biz.PaymentStatusPaid:
		reply.PaymentStatus = v3.PaymentStatus_PAID
	case biz.PaymentStatusExpired:
		reply.PaymentStatus = v3.PaymentStatus_EXPIRED
	}

	return reply
}

func convertExpressOrderReplyList(in []*biz.ExpressOrder) []*v3.ExpressOrderReply {
	out := make([]*v3.ExpressOrderReply, len(in))
	for i := range in {
		out[i] = convertExpressOrderReply(in[i])
	}
	return out
}
