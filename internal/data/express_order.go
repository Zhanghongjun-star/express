package data

import (
	"context"
	"time"

	"shunfeng-miniprogram/internal/biz"
)

type expressOrderRepo struct{}

func NewExpressOrderRepo() biz.ExpressOrderRepo {
	return &expressOrderRepo{}
}

func (r *expressOrderRepo) Create(ctx context.Context, order *biz.ExpressOrder) (*biz.ExpressOrder, error) {
	po := toExpressOrderPO(order)
	if err := DB.WithContext(ctx).Create(po).Error; err != nil {
		return nil, biz.ErrOrderCreateFailed
	}
	return toExpressOrderBiz(po), nil
}

func (r *expressOrderRepo) FindByID(ctx context.Context, userID, orderID int64) (*biz.ExpressOrder, error) {
	var po ExpressOrder
	if err := DB.WithContext(ctx).
		Where("id = ? AND user_id = ? AND del_flag = 0", orderID, userID).
		First(&po).Error; err != nil {
		return nil, biz.ErrOrderNotFound
	}
	return toExpressOrderBiz(&po), nil
}

func (r *expressOrderRepo) Search(ctx context.Context, userID int64, keyword string, offset, limit int) ([]*biz.ExpressOrder, int32, error) {
	query := DB.WithContext(ctx).Model(&ExpressOrder{}).
		Where("user_id = ? AND del_flag = 0", userID)

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"order_no LIKE ? OR express_no LIKE ? OR sender_name LIKE ? OR receiver_name LIKE ? OR sender_phone LIKE ? OR receiver_phone LIKE ?",
			like, like, like, like, like, like,
		)
	}

	var total int64
	query.Count(&total)

	var pos []ExpressOrder
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	return toExpressOrderBizList(pos), int32(total), nil
}

func (r *expressOrderRepo) ListByCategory(ctx context.Context, userID int64, category string, offset, limit int) ([]*biz.ExpressOrder, int32, error) {
	query := DB.WithContext(ctx).Model(&ExpressOrder{}).Where("del_flag = 0")

	switch category {
	case "sent":
		query = query.Where("user_id = ?", userID)
	case "received":
		query = query.Where("receiver_phone IN (SELECT phone FROM sf_user WHERE user_id = ?)", userID)
	case "followed":
		query = query.Where("user_id = ? AND is_followed = 1", userID)
	case "unpaid":
		query = query.Where("user_id = ? AND payment_status = 'pending'", userID)
	}

	var total int64
	query.Count(&total)

	var pos []ExpressOrder
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	return toExpressOrderBizList(pos), int32(total), nil
}

func (r *expressOrderRepo) Delete(ctx context.Context, userID, orderID int64) error {
	result := DB.WithContext(ctx).Model(&ExpressOrder{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		Update("del_flag", 1)
	if result.RowsAffected == 0 {
		return biz.ErrOrderNotFound
	}
	return result.Error
}

func (r *expressOrderRepo) UpdateTag(ctx context.Context, userID, orderID int64, tag string) error {
	result := DB.WithContext(ctx).Model(&ExpressOrder{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		Updates(map[string]any{"tag": tag, "updated_at": time.Now()})
	if result.RowsAffected == 0 {
		return biz.ErrOrderNotFound
	}
	return result.Error
}

func (r *expressOrderRepo) UpdateFollow(ctx context.Context, userID, orderID int64, follow bool) error {
	result := DB.WithContext(ctx).Model(&ExpressOrder{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		Updates(map[string]any{"is_followed": follow, "updated_at": time.Now()})
	if result.RowsAffected == 0 {
		return biz.ErrOrderNotFound
	}
	return result.Error
}

func (r *expressOrderRepo) UpdateShareURL(ctx context.Context, userID, orderID int64, url string) error {
	return DB.WithContext(ctx).Model(&ExpressOrder{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		Updates(map[string]any{"share_url": url, "updated_at": time.Now()}).Error
}

func (r *expressOrderRepo) UpdatePayURL(ctx context.Context, userID, orderID int64, payURL, tradeNo string) error {
	return DB.WithContext(ctx).Model(&ExpressOrder{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		Updates(map[string]any{"pay_url": payURL, "trade_no": tradeNo, "updated_at": time.Now()}).Error
}

func (r *expressOrderRepo) UpdatePaymentStatusByTradeNo(ctx context.Context, tradeNo, paymentStatus string) error {
	updates := map[string]any{"payment_status": paymentStatus, "updated_at": time.Now()}
	if paymentStatus == biz.PaymentStatusPaid {
		updates["paid_at"] = time.Now()
	}
	return DB.WithContext(ctx).Model(&ExpressOrder{}).
		Where("trade_no = ?", tradeNo).
		Updates(updates).Error
}

func (r *expressOrderRepo) UpdateStatus(ctx context.Context, userID, orderID int64, status string) error {
	return DB.WithContext(ctx).Model(&ExpressOrder{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		Updates(map[string]any{"status": status, "updated_at": time.Now()}).Error
}

func (r *expressOrderRepo) UpdatePaymentStatus(ctx context.Context, userID, orderID int64, paymentStatus string) error {
	updates := map[string]any{"payment_status": paymentStatus, "updated_at": time.Now()}
	if paymentStatus == biz.PaymentStatusPaid {
		updates["paid_at"] = time.Now()
	}
	return DB.WithContext(ctx).Model(&ExpressOrder{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		Updates(updates).Error
}

func toExpressOrderPO(in *biz.ExpressOrder) *ExpressOrder {
	if in == nil {
		return nil
	}
	return &ExpressOrder{
		UserID:            in.UserID,
		OrderNo:           in.OrderNo,
		ExpressNo:         in.ExpressNo,
		SenderName:        in.SenderName,
		SenderPhone:       in.SenderPhone,
		SenderProvince:    in.SenderProvince,
		SenderCity:        in.SenderCity,
		SenderDistrict:    in.SenderDistrict,
		SenderDetail:      in.SenderDetail,
		SenderLat:         in.SenderLat,
		SenderLng:         in.SenderLng,
		ReceiverName:      in.ReceiverName,
		ReceiverPhone:     in.ReceiverPhone,
		ReceiverProvince:  in.ReceiverProvince,
		ReceiverCity:      in.ReceiverCity,
		ReceiverDistrict:  in.ReceiverDistrict,
		ReceiverDetail:    in.ReceiverDetail,
		ReceiverLat:       in.ReceiverLat,
		ReceiverLng:       in.ReceiverLng,
		Weight:            in.Weight,
		Length:            in.Length,
		Width:             in.Width,
		Height:            in.Height,
		BaseFreight:       in.BaseFreight,
		InsureFee:         in.InsureFee,
		TotalFreight:      in.TotalFreight,
		ChannelCode:       in.ChannelCode,
		LockerID:          in.LockerID,
		BoxType:           in.BoxType,
		ServicePointID:    in.ServicePointID,
		PickupDate:        in.PickupDate,
		PickupSlotCode:    in.PickupSlotCode,
		PickupStartTime:   stringPtrIfNotEmpty(in.PickupStartTime),
		PickupEndTime:     stringPtrIfNotEmpty(in.PickupEndTime),
		PaymentMethod:     in.PaymentMethod,
		PaymentStatus:     in.PaymentStatus,
		PaidAt:            timePtrIfNotZero(in.PaidAt),
		PrivacyProtection: in.PrivacyProtection,
		Status:            in.Status,
		DelFlag:           in.DelFlag,
		Tag:               in.Tag,
		IsFollowed:        in.IsFollowed,
		ShareURL:          in.ShareURL,
		PayURL:            in.PayURL,
		TradeNo:           in.TradeNo,
		CreatedAt:         in.CreatedAt,
		UpdatedAt:         in.UpdatedAt,
	}
}

func toExpressOrderBiz(in *ExpressOrder) *biz.ExpressOrder {
	if in == nil {
		return nil
	}
	return &biz.ExpressOrder{
		ID:                in.ID,
		UserID:            in.UserID,
		OrderNo:           in.OrderNo,
		ExpressNo:         in.ExpressNo,
		SenderName:        in.SenderName,
		SenderPhone:       in.SenderPhone,
		SenderProvince:    in.SenderProvince,
		SenderCity:        in.SenderCity,
		SenderDistrict:    in.SenderDistrict,
		SenderDetail:      in.SenderDetail,
		SenderLat:         in.SenderLat,
		SenderLng:         in.SenderLng,
		ReceiverName:      in.ReceiverName,
		ReceiverPhone:     in.ReceiverPhone,
		ReceiverProvince:  in.ReceiverProvince,
		ReceiverCity:      in.ReceiverCity,
		ReceiverDistrict:  in.ReceiverDistrict,
		ReceiverDetail:    in.ReceiverDetail,
		ReceiverLat:       in.ReceiverLat,
		ReceiverLng:       in.ReceiverLng,
		Weight:            in.Weight,
		Length:            in.Length,
		Width:             in.Width,
		Height:            in.Height,
		BaseFreight:       in.BaseFreight,
		InsureFee:         in.InsureFee,
		TotalFreight:      in.TotalFreight,
		ChannelCode:       in.ChannelCode,
		LockerID:          in.LockerID,
		BoxType:           in.BoxType,
		ServicePointID:    in.ServicePointID,
		PickupDate:        in.PickupDate,
		PickupSlotCode:    in.PickupSlotCode,
		PickupStartTime:   stringValue(in.PickupStartTime),
		PickupEndTime:     stringValue(in.PickupEndTime),
		PaymentMethod:     in.PaymentMethod,
		PaymentStatus:     in.PaymentStatus,
		PaidAt:            timeValue(in.PaidAt),
		PrivacyProtection: in.PrivacyProtection,
		Status:            in.Status,
		DelFlag:           in.DelFlag,
		Tag:               in.Tag,
		IsFollowed:        in.IsFollowed,
		ShareURL:          in.ShareURL,
		PayURL:            in.PayURL,
		TradeNo:           in.TradeNo,
		CreatedAt:         in.CreatedAt,
		UpdatedAt:         in.UpdatedAt,
	}
}

func toExpressOrderBizList(in []ExpressOrder) []*biz.ExpressOrder {
	out := make([]*biz.ExpressOrder, len(in))
	for i := range in {
		out[i] = toExpressOrderBiz(&in[i])
	}
	return out
}

func stringPtrIfNotEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func timePtrIfNotZero(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}

func timeValue(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return *value
}
