package service

import (
	"context"
	"time"

	v3 "shunfeng-miniprogram/api/order/v3"
	"shunfeng-miniprogram/internal/biz"

	"go.einride.tech/aip/pagination"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const defaultExpressOrderPageSize = 20

// ExpressOrderService 寄件订单服务（order.v3）。
type ExpressOrderService struct {
	v3.UnimplementedOrderServiceServer

	uc *biz.ExpressOrderUsecase
}

// NewExpressOrderService 创建寄件订单服务。
func NewExpressOrderService(uc *biz.ExpressOrderUsecase) *ExpressOrderService {
	return &ExpressOrderService{uc: uc}
}

// CreateExpressOrder 创建寄件订单。
func (s *ExpressOrderService) CreateExpressOrder(ctx context.Context, req *v3.CreateExpressOrderRequest) (*v3.ExpressOrder, error) {
	order, err := s.uc.CreateExpressOrder(ctx, convertExpressOrder(req.GetOrder()))
	if err != nil {
		return nil, err
	}
	return convertExpressOrderReply(order), nil
}

// GetExpressOrder 获取订单详情。
func (s *ExpressOrderService) GetExpressOrder(ctx context.Context, req *v3.GetExpressOrderRequest) (*v3.ExpressOrder, error) {
	order, err := s.uc.GetExpressOrder(ctx, req.GetOrderNo())
	if err != nil {
		return nil, err
	}
	return convertExpressOrderReply(order), nil
}

// ListExpressOrders 列出订单。
func (s *ExpressOrderService) ListExpressOrders(ctx context.Context, req *v3.ListExpressOrdersRequest) (*v3.ExpressOrderSet, error) {
	pageToken, err := pagination.ParsePageToken(req)
	if err != nil {
		return nil, err
	}
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = defaultExpressOrderPageSize
	}
	orders, total, err := s.uc.ListExpressOrders(ctx, req.GetUserNo(), req.GetOrderStatus(),
		biz.ListOffset(int(pageToken.Offset)),
		biz.ListLimit(pageSize),
	)
	if err != nil {
		return nil, err
	}
	set := &v3.ExpressOrderSet{
		Orders:     make([]*v3.ExpressOrder, 0, len(orders)),
		TotalCount: total,
	}
	if len(orders) >= pageSize {
		set.NextPageToken = pageToken.Next(req).String()
	}
	for _, o := range orders {
		set.Orders = append(set.Orders, convertExpressOrderReply(o))
	}
	return set, nil
}

// UpdateExpressOrder 更新订单（字段掩码局部更新）。
func (s *ExpressOrderService) UpdateExpressOrder(ctx context.Context, req *v3.UpdateExpressOrderRequest) (*v3.ExpressOrder, error) {
	current, err := s.uc.GetExpressOrder(ctx, req.GetOrder().GetOrderNo())
	if err != nil {
		return nil, err
	}
	applyExpressOrderMask(current, req.GetOrder(), req.GetUpdateMask().GetPaths())
	order, err := s.uc.UpdateExpressOrder(ctx, current, req.GetUpdateMask().GetPaths())
	if err != nil {
		return nil, err
	}
	return convertExpressOrderReply(order), nil
}

// DeleteExpressOrder 删除订单。
func (s *ExpressOrderService) DeleteExpressOrder(ctx context.Context, req *v3.DeleteExpressOrderRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeleteExpressOrder(ctx, req.GetOrderNo()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// convertExpressOrder 将 v3.ExpressOrder 转换为 biz.ExpressOrder（DTO → DO）。
func convertExpressOrder(in *v3.ExpressOrder) *biz.ExpressOrder {
	if in == nil {
		return nil
	}
	do := &biz.ExpressOrder{
		OrderNo:          in.GetOrderNo(),
		TrackingNo:       in.GetTrackingNo(),
		UserNo:           in.GetUserNo(),
		ExpressCompany:   in.GetExpressCompany(),
		SenderName:       in.GetSenderName(),
		SenderPhone:      in.GetSenderPhone(),
		SenderProvince:   in.GetSenderProvince(),
		SenderCity:       in.GetSenderCity(),
		SenderDistrict:   in.GetSenderDistrict(),
		SenderAddress:    in.GetSenderAddress(),
		ReceiverName:     in.GetReceiverName(),
		ReceiverPhone:    in.GetReceiverPhone(),
		ReceiverProvince: in.GetReceiverProvince(),
		ReceiverCity:     in.GetReceiverCity(),
		ReceiverDistrict: in.GetReceiverDistrict(),
		ReceiverAddress:  in.GetReceiverAddress(),
		ItemName:         in.GetItemName(),
		ItemCategory:     in.GetItemCategory(),
		ItemQuantity:     int(in.GetItemQuantity()),
		Weight:           in.GetWeight(),
		Length:           in.GetLength(),
		Width:            in.GetWidth(),
		Height:           in.GetHeight(),
		DeclaredValue:    in.GetDeclaredValue(),
		IsFragile:        in.GetIsFragile(),
		IsBattery:        in.GetIsBattery(),
		FreightFee:       in.GetFreightFee(),
		InsuranceFee:     in.GetInsuranceFee(),
		PackageFee:       in.GetPackageFee(),
		TotalFee:         in.GetTotalFee(),
		PayMethod:        in.GetPayMethod(),
		PayStatus:        in.GetPayStatus(),
		OrderStatus:      in.GetOrderStatus(),
		SignedBy:         in.GetSignedBy(),
		Remark:           in.GetRemark(),
		// 包装材料列表：proto repeated → biz 结构体列表
		Packages:        convertPackages(in.GetPackages()),
		// 费用计算过程明细（JSON 字符串）
		PackageDetail:   in.GetPackageDetail(),
	}
	if t := in.GetPickupTime(); t != nil {
		tm := t.AsTime()
		do.PickupTime = &tm
	}
	if t := in.GetDeliveredTime(); t != nil {
		tm := t.AsTime()
		do.DeliveredTime = &tm
	}
	return do
}

// convertExpressOrderReply 将 biz.ExpressOrder 转换为 v3.ExpressOrder（DO → DTO）。
func convertExpressOrderReply(in *biz.ExpressOrder) *v3.ExpressOrder {
	if in == nil {
		return nil
	}
	return &v3.ExpressOrder{
		Id:              in.ID,
		OrderNo:         in.OrderNo,
		TrackingNo:      in.TrackingNo,
		UserNo:          in.UserNo,
		ExpressCompany:  in.ExpressCompany,
		SenderName:      in.SenderName,
		SenderPhone:     in.SenderPhone,
		SenderProvince:  in.SenderProvince,
		SenderCity:      in.SenderCity,
		SenderDistrict:  in.SenderDistrict,
		SenderAddress:   in.SenderAddress,
		ReceiverName:    in.ReceiverName,
		ReceiverPhone:   in.ReceiverPhone,
		ReceiverProvince: in.ReceiverProvince,
		ReceiverCity:    in.ReceiverCity,
		ReceiverDistrict: in.ReceiverDistrict,
		ReceiverAddress: in.ReceiverAddress,
		ItemName:        in.ItemName,
		ItemCategory:    in.ItemCategory,
		ItemQuantity:    int32(in.ItemQuantity),
		Weight:          in.Weight,
		Length:          in.Length,
		Width:           in.Width,
		Height:          in.Height,
		DeclaredValue:   in.DeclaredValue,
		IsFragile:       in.IsFragile,
		IsBattery:       in.IsBattery,
		FreightFee:      in.FreightFee,
		InsuranceFee:    in.InsuranceFee,
		PackageFee:      in.PackageFee,
		TotalFee:        in.TotalFee,
		PayMethod:       in.PayMethod,
		PayStatus:       in.PayStatus,
		OrderStatus:     in.OrderStatus,
		SignedBy:        in.SignedBy,
		Remark:          in.Remark,
		// 包装材料列表：biz 结构体列表 → proto repeated
		Packages:        convertPackagesReply(in.Packages),
		// 费用计算过程明细（JSON 字符串）
		PackageDetail:   in.PackageDetail,
		PickupTime:      tsOrNil(in.PickupTime),
		DeliveredTime:   tsOrNil(in.DeliveredTime),
		CreatedAt:       timestamppb.New(in.CreatedAt),
		UpdatedAt:       timestamppb.New(in.UpdatedAt),
	}
}

// convertPackages 将 proto 包装材料列表转换为 biz 结构体列表（DTO → DO）。
func convertPackages(in []*v3.PackageItem) []biz.PackageItem {
	if len(in) == 0 {
		return nil
	}
	out := make([]biz.PackageItem, 0, len(in))
	for _, p := range in {
		out = append(out, biz.PackageItem{
			MaterialID: p.GetMaterialId(),
			Quantity:   int(p.GetQuantity()),
		})
	}
	return out
}

// convertPackagesReply 将 biz 包装材料列表转换为 proto 包装材料列表（DO → DTO）。
func convertPackagesReply(in []biz.PackageItem) []*v3.PackageItem {
	if len(in) == 0 {
		return nil
	}
	out := make([]*v3.PackageItem, 0, len(in))
	for _, p := range in {
		out = append(out, &v3.PackageItem{
			MaterialId: p.MaterialID,
			Quantity:   int32(p.Quantity),
		})
	}
	return out
}

// tsOrNil 将 *time.Time 转换为 *timestamppb.Timestamp，nil 时返回 nil。
func tsOrNil(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

// applyExpressOrderMask 根据字段掩码更新订单字段。
func applyExpressOrderMask(dst *biz.ExpressOrder, src *v3.ExpressOrder, paths []string) {
	if dst == nil || src == nil {
		return
	}
	for _, path := range paths {
		switch path {
		case "*":
			*dst = *convertExpressOrder(src)
			return
		case "tracking_no":
			dst.TrackingNo = src.GetTrackingNo()
		case "express_company":
			dst.ExpressCompany = src.GetExpressCompany()
		case "sender_name":
			dst.SenderName = src.GetSenderName()
		case "sender_phone":
			dst.SenderPhone = src.GetSenderPhone()
		case "sender_province":
			dst.SenderProvince = src.GetSenderProvince()
		case "sender_city":
			dst.SenderCity = src.GetSenderCity()
		case "sender_district":
			dst.SenderDistrict = src.GetSenderDistrict()
		case "sender_address":
			dst.SenderAddress = src.GetSenderAddress()
		case "receiver_name":
			dst.ReceiverName = src.GetReceiverName()
		case "receiver_phone":
			dst.ReceiverPhone = src.GetReceiverPhone()
		case "receiver_province":
			dst.ReceiverProvince = src.GetReceiverProvince()
		case "receiver_city":
			dst.ReceiverCity = src.GetReceiverCity()
		case "receiver_district":
			dst.ReceiverDistrict = src.GetReceiverDistrict()
		case "receiver_address":
			dst.ReceiverAddress = src.GetReceiverAddress()
		case "item_name":
			dst.ItemName = src.GetItemName()
		case "item_category":
			dst.ItemCategory = src.GetItemCategory()
		case "item_quantity":
			dst.ItemQuantity = int(src.GetItemQuantity())
		case "weight":
			dst.Weight = src.GetWeight()
		case "length":
			dst.Length = src.GetLength()
		case "width":
			dst.Width = src.GetWidth()
		case "height":
			dst.Height = src.GetHeight()
		case "declared_value":
			dst.DeclaredValue = src.GetDeclaredValue()
		case "is_fragile":
			dst.IsFragile = src.GetIsFragile()
		case "is_battery":
			dst.IsBattery = src.GetIsBattery()
		case "freight_fee":
			dst.FreightFee = src.GetFreightFee()
		case "insurance_fee":
			dst.InsuranceFee = src.GetInsuranceFee()
		case "package_fee":
			dst.PackageFee = src.GetPackageFee()
		case "total_fee":
			dst.TotalFee = src.GetTotalFee()
		case "pay_method":
			dst.PayMethod = src.GetPayMethod()
		case "pay_status":
			dst.PayStatus = src.GetPayStatus()
		case "order_status":
			dst.OrderStatus = src.GetOrderStatus()
		case "pickup_time":
			if t := src.GetPickupTime(); t != nil {
				tm := t.AsTime()
				dst.PickupTime = &tm
			}
		case "delivered_time":
			if t := src.GetDeliveredTime(); t != nil {
				tm := t.AsTime()
				dst.DeliveredTime = &tm
			}
		case "signed_by":
			dst.SignedBy = src.GetSignedBy()
		case "remark":
			dst.Remark = src.GetRemark()
		case "packages":
			// 包装材料列表：按新传入值整体替换
			dst.Packages = convertPackages(src.GetPackages())
		case "package_detail":
			dst.PackageDetail = src.GetPackageDetail()
		}
	}
}
