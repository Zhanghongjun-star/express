package data

import (
	"context"
	"encoding/json"
	"time"

	"shunfeng-miniprogram/internal/biz"
)

// expressOrderRepo 寄件订单仓储的 GORM 实现。
type expressOrderRepo struct{}

// NewExpressOrderRepo 创建寄件订单仓储。
func NewExpressOrderRepo() biz.ExpressOrderRepo {
	return &expressOrderRepo{}
}

func (r *expressOrderRepo) Create(ctx context.Context, o *biz.ExpressOrder) (*biz.ExpressOrder, error) {
	po := newExpressOrderPO(o)
	now := time.Now()
	po.CreatedAt = now
	po.UpdatedAt = now
	if err := DB.WithContext(ctx).Create(po).Error; err != nil {
		return nil, err
	}
	return toBizExpressOrder(po), nil
}

func (r *expressOrderRepo) GetByNo(ctx context.Context, orderNo string) (*biz.ExpressOrder, error) {
	var po ExpressOrder
	if err := DB.WithContext(ctx).Where("order_no = ?", orderNo).First(&po).Error; err != nil {
		return nil, biz.ErrExpressOrderNotFound
	}
	return toBizExpressOrder(&po), nil
}

func (r *expressOrderRepo) List(ctx context.Context, userNo, status string, opts ...biz.ListOption) ([]*biz.ExpressOrder, int32, error) {
	options := biz.ListOptions{Limit: 20}
	for _, opt := range opts {
		opt(&options)
	}
	if options.Offset < 0 || options.Limit <= 0 {
		return nil, 0, biz.ErrExpressOrderInvalidArgument
	}

	tx := DB.WithContext(ctx).Model(&ExpressOrder{}).Where("user_no = ?", userNo)
	if status != "" {
		tx = tx.Where("order_status = ?", status)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var pos []ExpressOrder
	q := DB.WithContext(ctx).Where("user_no = ?", userNo)
	if status != "" {
		q = q.Where("order_status = ?", status)
	}
	if err := q.Order("id DESC").Offset(options.Offset).Limit(options.Limit).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	orders := make([]*biz.ExpressOrder, len(pos))
	for i := range pos {
		orders[i] = toBizExpressOrder(&pos[i])
	}
	return orders, int32(total), nil
}

func (r *expressOrderRepo) Update(ctx context.Context, o *biz.ExpressOrder, mask []string) (*biz.ExpressOrder, error) {
	var po ExpressOrder
	if err := DB.WithContext(ctx).Where("order_no = ?", o.OrderNo).First(&po).Error; err != nil {
		return nil, biz.ErrExpressOrderNotFound
	}

	updates := expressOrderUpdateMap(o, mask)
	if len(updates) == 0 {
		return toBizExpressOrder(&po), nil
	}
	updates["updated_at"] = time.Now()

	if err := DB.WithContext(ctx).Model(&po).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := DB.WithContext(ctx).Where("order_no = ?", o.OrderNo).First(&po).Error; err != nil {
		return nil, err
	}
	return toBizExpressOrder(&po), nil
}

func (r *expressOrderRepo) Delete(ctx context.Context, orderNo string) error {
	result := DB.WithContext(ctx).Where("order_no = ?", orderNo).Delete(&ExpressOrder{})
	if result.RowsAffected == 0 {
		return biz.ErrExpressOrderNotFound
	}
	return result.Error
}

// newExpressOrderPO 将 biz.ExpressOrder 转换为 ExpressOrder（DO → PO）。
func newExpressOrderPO(in *biz.ExpressOrder) *ExpressOrder {
	if in == nil {
		return nil
	}
	return &ExpressOrder{
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
		ItemQuantity:    in.ItemQuantity,
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
		PickupTime:      in.PickupTime,
		DeliveredTime:   in.DeliveredTime,
		SignedBy:        in.SignedBy,
		Remark:          in.Remark,
		// 包装材料明细：biz 列表序列化为 JSON 字符串存储
		Packages:        marshalPackages(in.Packages),
		// 费用计算过程明细（JSON 字符串）
		PackageDetail:   in.PackageDetail,
	}
}

// marshalPackages 将 biz 包装材料列表序列化为 JSON 字符串。
func marshalPackages(items []biz.PackageItem) string {
	if len(items) == 0 {
		return ""
	}
	b, err := json.Marshal(items)
	if err != nil {
		return ""
	}
	return string(b)
}

// unmarshalPackages 将 JSON 字符串解析为 biz 包装材料列表。
func unmarshalPackages(s string) []biz.PackageItem {
	if s == "" {
		return nil
	}
	var items []biz.PackageItem
	if err := json.Unmarshal([]byte(s), &items); err != nil {
		return nil
	}
	return items
}

// toBizExpressOrder 将 ExpressOrder 转换为 biz.ExpressOrder（PO → DO）。
func toBizExpressOrder(in *ExpressOrder) *biz.ExpressOrder {
	if in == nil {
		return nil
	}
	return &biz.ExpressOrder{
		ID:              in.ID,
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
		ItemQuantity:    int(in.ItemQuantity),
		Weight:          in.Weight,
		Length:          in.Length,
		Width:           in.Width,
		Height:          in.Height,
		DeclaredValue:   in.DeclaredValue,
		IsFragile:       in.IsFragile,
		IsBattery:       in.IsBattery,
		FreightFee:      in.FreightFee,
		InsuranceFee:   in.InsuranceFee,
		PackageFee:      in.PackageFee,
		TotalFee:        in.TotalFee,
		PayMethod:       in.PayMethod,
		PayStatus:       in.PayStatus,
		OrderStatus:     in.OrderStatus,
		PickupTime:      in.PickupTime,
		DeliveredTime:   in.DeliveredTime,
		SignedBy:        in.SignedBy,
		Remark:          in.Remark,
		// 包装材料明细：JSON 字符串反序列化为 biz 列表
		Packages:        unmarshalPackages(in.Packages),
		// 费用计算过程明细（JSON 字符串）
		PackageDetail:   in.PackageDetail,
		CreatedAt:       in.CreatedAt,
		UpdatedAt:       in.UpdatedAt,
	}
}

// expressOrderUpdateMap 根据字段掩码构建 GORM 更新映射。
func expressOrderUpdateMap(in *biz.ExpressOrder, paths []string) map[string]interface{} {
	m := map[string]interface{}{}
	for _, p := range paths {
		switch p {
		case "*":
			return expressOrderFullMap(in)
		case "tracking_no":
			m["tracking_no"] = in.TrackingNo
		case "express_company":
			m["express_company"] = in.ExpressCompany
		case "sender_name":
			m["sender_name"] = in.SenderName
		case "sender_phone":
			m["sender_phone"] = in.SenderPhone
		case "sender_province":
			m["sender_province"] = in.SenderProvince
		case "sender_city":
			m["sender_city"] = in.SenderCity
		case "sender_district":
			m["sender_district"] = in.SenderDistrict
		case "sender_address":
			m["sender_address"] = in.SenderAddress
		case "receiver_name":
			m["receiver_name"] = in.ReceiverName
		case "receiver_phone":
			m["receiver_phone"] = in.ReceiverPhone
		case "receiver_province":
			m["receiver_province"] = in.ReceiverProvince
		case "receiver_city":
			m["receiver_city"] = in.ReceiverCity
		case "receiver_district":
			m["receiver_district"] = in.ReceiverDistrict
		case "receiver_address":
			m["receiver_address"] = in.ReceiverAddress
		case "item_name":
			m["item_name"] = in.ItemName
		case "item_category":
			m["item_category"] = in.ItemCategory
		case "item_quantity":
			m["item_quantity"] = in.ItemQuantity
		case "weight":
			m["weight"] = in.Weight
		case "length":
			m["length"] = in.Length
		case "width":
			m["width"] = in.Width
		case "height":
			m["height"] = in.Height
		case "declared_value":
			m["declared_value"] = in.DeclaredValue
		case "is_fragile":
			m["is_fragile"] = in.IsFragile
		case "is_battery":
			m["is_battery"] = in.IsBattery
		case "freight_fee":
			m["freight_fee"] = in.FreightFee
		case "insurance_fee":
			m["insurance_fee"] = in.InsuranceFee
		case "package_fee":
			m["package_fee"] = in.PackageFee
		case "total_fee":
			m["total_fee"] = in.TotalFee
		case "pay_method":
			m["pay_method"] = in.PayMethod
		case "pay_status":
			m["pay_status"] = in.PayStatus
		case "order_status":
			m["order_status"] = in.OrderStatus
		case "pickup_time":
			m["pickup_time"] = in.PickupTime
		case "delivered_time":
			m["delivered_time"] = in.DeliveredTime
		case "signed_by":
			m["signed_by"] = in.SignedBy
		case "remark":
			m["remark"] = in.Remark
		case "packages":
			// 包装材料明细：按新值整体替换
			m["packages"] = marshalPackages(in.Packages)
		case "package_detail":
			m["package_detail"] = in.PackageDetail
		}
	}
	return m
}

// expressOrderFullMap 全量更新映射（用于 mask 为 "*" 的情况）。
func expressOrderFullMap(in *biz.ExpressOrder) map[string]interface{} {
	return map[string]interface{}{
		"tracking_no":       in.TrackingNo,
		"express_company":   in.ExpressCompany,
		"sender_name":       in.SenderName,
		"sender_phone":      in.SenderPhone,
		"sender_province":   in.SenderProvince,
		"sender_city":       in.SenderCity,
		"sender_district":   in.SenderDistrict,
		"sender_address":    in.SenderAddress,
		"receiver_name":     in.ReceiverName,
		"receiver_phone":    in.ReceiverPhone,
		"receiver_province": in.ReceiverProvince,
		"receiver_city":     in.ReceiverCity,
		"receiver_district": in.ReceiverDistrict,
		"receiver_address":  in.ReceiverAddress,
		"item_name":         in.ItemName,
		"item_category":     in.ItemCategory,
		"item_quantity":     in.ItemQuantity,
		"weight":            in.Weight,
		"length":            in.Length,
		"width":             in.Width,
		"height":            in.Height,
		"declared_value":    in.DeclaredValue,
		"is_fragile":        in.IsFragile,
		"is_battery":        in.IsBattery,
		"freight_fee":       in.FreightFee,
		"insurance_fee":     in.InsuranceFee,
		"package_fee":       in.PackageFee,
		"total_fee":         in.TotalFee,
		"pay_method":        in.PayMethod,
		"pay_status":        in.PayStatus,
		"order_status":      in.OrderStatus,
		"pickup_time":       in.PickupTime,
		"delivered_time":    in.DeliveredTime,
		"signed_by":         in.SignedBy,
		"remark":            in.Remark,
		"packages":          marshalPackages(in.Packages),
		"package_detail":    in.PackageDetail,
	}
}
