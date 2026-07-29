package data

import (
	"context"
	"time"

	"shunfeng-miniprogram/internal/biz"

	"gorm.io/gorm"
)

// ──────────────────────────────────────────────
// Channel repo
// ──────────────────────────────────────────────

type channelRepo struct{}

func NewChannelRepo() biz.ChannelRepo {
	return &channelRepo{}
}

func (r *channelRepo) ListChannels(ctx context.Context) ([]*biz.ShippingChannel, error) {
	var pos []ShippingChannel
	if err := DB.WithContext(ctx).Where("status = 1").Order("sort_no ASC").Find(&pos).Error; err != nil {
		return nil, err
	}
	return toChannelBizList(pos), nil
}

func (r *channelRepo) FindChannelByCode(ctx context.Context, channelCode string) (*biz.ShippingChannel, error) {
	var po ShippingChannel
	if err := DB.WithContext(ctx).Where("channel_code = ?", channelCode).First(&po).Error; err != nil {
		return nil, biz.ErrChannelUnavailable
	}
	return toChannelBiz(&po), nil
}

func (r *channelRepo) FindAreaByChannelDistrict(ctx context.Context, channelCode, districtCode string) (*biz.ChannelArea, error) {
	var po ShippingChannelArea
	if err := DB.WithContext(ctx).
		Where("channel_code = ? AND district_code = ?", channelCode, districtCode).
		First(&po).Error; err != nil {
		return nil, biz.ErrChannelUnavailable
	}
	return toChannelAreaBiz(&po), nil
}

func (r *channelRepo) ListAreasByDistrict(ctx context.Context, districtCode string) ([]*biz.ChannelArea, error) {
	var pos []ShippingChannelArea
	if err := DB.WithContext(ctx).Where("district_code = ?", districtCode).Find(&pos).Error; err != nil {
		return nil, err
	}
	return toChannelAreaBizList(pos), nil
}

// ──────────────────────────────────────────────
// Locker repo
// ──────────────────────────────────────────────

type lockerRepo struct{}

func NewLockerRepo() biz.LockerRepo {
	return &lockerRepo{}
}

func (r *lockerRepo) ListByLocation(ctx context.Context, districtCode string, lat, lng float64, radius int32, offset, limit int) ([]*biz.Locker, int, error) {
	var total int64
	query := DB.WithContext(ctx).Model(&Locker{}).Where("status = 1 AND district_code = ?", districtCode)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var pos []Locker
	if err := query.Offset(offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	return toLockerBizList(pos), int(total), nil
}

func (r *lockerRepo) FindByID(ctx context.Context, id int64) (*biz.Locker, error) {
	var po Locker
	if err := DB.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, biz.ErrLockerNotFound
	}
	return toLockerBiz(&po), nil
}

// ──────────────────────────────────────────────
// LockerBox repo
// ──────────────────────────────────────────────

type lockerBoxRepo struct{}

func NewLockerBoxRepo() biz.LockerBoxRepo {
	return &lockerBoxRepo{}
}

func (r *lockerBoxRepo) GetBoxTypesByLocker(ctx context.Context, lockerID int64) ([]*biz.LockerBox, error) {
	var pos []LockerBox
	if err := DB.WithContext(ctx).Where("locker_id = ?", lockerID).Find(&pos).Error; err != nil {
		return nil, err
	}
	return toLockerBoxBizList(pos), nil
}

func (r *lockerBoxRepo) LockBox(ctx context.Context, id int64, orderID int64, expireAt time.Time, version int64) error {
	result := DB.WithContext(ctx).Model(&LockerBox{}).
		Where("id = ? AND status = 1 AND version = ?", id, version).
		Updates(map[string]any{
			"status":            2,
			"reserved_order_id": orderID,
			"reserved_at":       time.Now(),
			"expire_at":         expireAt,
			"version":           gorm.Expr("version + 1"),
		})
	if result.RowsAffected == 0 {
		return biz.ErrBoxTypeFull
	}
	return result.Error
}

func (r *lockerBoxRepo) ReleaseBox(ctx context.Context, id int64, orderID int64, version int64) error {
	result := DB.WithContext(ctx).Model(&LockerBox{}).
		Where("id = ? AND reserved_order_id = ? AND version = ?", id, orderID, version).
		Updates(map[string]any{
			"status":            1,
			"reserved_order_id": nil,
			"reserved_at":       nil,
			"expire_at":         nil,
			"version":           gorm.Expr("version + 1"),
		})
	if result.RowsAffected == 0 {
		return nil
	}
	return result.Error
}

func (r *lockerBoxRepo) ReleaseExpiredBoxes(ctx context.Context) (int64, error) {
	result := DB.WithContext(ctx).Model(&LockerBox{}).
		Where("status = 2 AND expire_at IS NOT NULL AND expire_at <= ?", time.Now()).
		Updates(map[string]any{
			"status":            1,
			"reserved_order_id": nil,
			"reserved_at":       nil,
			"expire_at":         nil,
			"version":           gorm.Expr("version + 1"),
		})
	return result.RowsAffected, result.Error
}

// ──────────────────────────────────────────────
// ServicePoint repo
// ──────────────────────────────────────────────

type servicePointRepo struct{}

func NewServicePointRepo() biz.ServicePointRepo {
	return &servicePointRepo{}
}

func (r *servicePointRepo) ListByLocation(ctx context.Context, pointType int32, districtCode string, lat, lng float64, offset, limit int) ([]*biz.ServicePoint, int, error) {
	var total int64
	query := DB.WithContext(ctx).Model(&ServicePoint{}).Where("status = 1 AND district_code = ?", districtCode)
	if pointType > 0 {
		query = query.Where("point_type = ?", pointType)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var pos []ServicePoint
	if err := query.Offset(offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	return toServicePointBizList(pos), int(total), nil
}

func (r *servicePointRepo) FindByID(ctx context.Context, id int64) (*biz.ServicePoint, error) {
	var po ServicePoint
	if err := DB.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, biz.ErrServicePointNotFound
	}
	return toServicePointBiz(&po), nil
}

// ──────────────────────────────────────────────
// Pickup repo
// ──────────────────────────────────────────────

type pickupRepo struct{}

func NewPickupRepo() biz.PickupRepo {
	return &pickupRepo{}
}

func (r *pickupRepo) ListDates(ctx context.Context, districtCode string) ([]string, error) {
	var dates []string
	if err := DB.WithContext(ctx).Model(&PickupTimeSlot{}).
		Select("DISTINCT pickup_date").
		Where("district_code = ? AND status = 1", districtCode).
		Order("pickup_date ASC").
		Pluck("pickup_date", &dates).Error; err != nil {
		return nil, err
	}
	return dates, nil
}

func (r *pickupRepo) ListByDate(ctx context.Context, districtCode, pickupDate string) ([]*biz.PickupTimeSlot, error) {
	var pos []PickupTimeSlot
	if err := DB.WithContext(ctx).
		Where("district_code = ? AND pickup_date = ? AND status = 1", districtCode, pickupDate).
		Order("start_time ASC").
		Find(&pos).Error; err != nil {
		return nil, err
	}
	return toPickupTimeSlotBizList(pos), nil
}

func (r *pickupRepo) FindSlotByCode(ctx context.Context, districtCode, pickupDate, slotCode string) (*biz.PickupTimeSlot, error) {
	var po PickupTimeSlot
	if err := DB.WithContext(ctx).
		Where("district_code = ? AND pickup_date = ? AND slot_code = ?", districtCode, pickupDate, slotCode).
		First(&po).Error; err != nil {
		return nil, biz.ErrSlotFull
	}
	return toPickupTimeSlotBiz(&po), nil
}

func (r *pickupRepo) ReserveSlot(ctx context.Context, id int64, version int64) error {
	result := DB.WithContext(ctx).Model(&PickupTimeSlot{}).
		Where("id = ? AND status = 1 AND reserved_count < capacity AND version = ?", id, version).
		Updates(map[string]any{
			"reserved_count": gorm.Expr("reserved_count + 1"),
			"version":        gorm.Expr("version + 1"),
		})
	if result.RowsAffected == 0 {
		return biz.ErrSlotFull
	}
	return result.Error
}

func (r *pickupRepo) ReleaseSlot(ctx context.Context, id int64, version int64) error {
	result := DB.WithContext(ctx).Model(&PickupTimeSlot{}).
		Where("id = ? AND reserved_count > 0 AND version = ?", id, version).
		Updates(map[string]any{
			"reserved_count": gorm.Expr("reserved_count - 1"),
			"version":        gorm.Expr("version + 1"),
		})
	if result.RowsAffected == 0 {
		return nil
	}
	return result.Error
}

func (r *pickupRepo) ReleaseExpiredSlots(ctx context.Context) (int64, error) {
	// 由定时任务清理超时预约占用流水后联动释放
	return 0, nil
}

// ──────────────────────────────────────────────
// PO → DO 转换函数
// ──────────────────────────────────────────────

func toChannelBiz(in *ShippingChannel) *biz.ShippingChannel {
	if in == nil {
		return nil
	}
	return &biz.ShippingChannel{
		ID:              in.ID,
		ChannelCode:     in.ChannelCode,
		ChannelName:     in.ChannelName,
		ChannelDesc:     in.ChannelDesc,
		Status:          in.Status,
		BaseFee:         in.BaseFee,
		NeedAddressGeo:  in.NeedAddressGeo == 1,
		NeedPickupSlot:  in.NeedPickupSlot == 1,
		NeedLockerBox:   in.NeedLockerBox == 1,
		NeedServicePoint: in.NeedServicePoint == 1,
		SortNo:          int(in.SortNo),
		CreatedAt:       in.CreatedAt,
		UpdatedAt:       in.UpdatedAt,
	}
}

func toChannelBizList(in []ShippingChannel) []*biz.ShippingChannel {
	out := make([]*biz.ShippingChannel, len(in))
	for i := range in {
		out[i] = toChannelBiz(&in[i])
	}
	return out
}

func toChannelAreaBiz(in *ShippingChannelArea) *biz.ChannelArea {
	if in == nil {
		return nil
	}
	return &biz.ChannelArea{
		ID:                in.ID,
		ChannelCode:       in.ChannelCode,
		ProvinceCode:      in.ProvinceCode,
		CityCode:          in.CityCode,
		DistrictCode:      in.DistrictCode,
		Status:            in.Status,
		UnavailableReason: in.UnavailableReason,
		ExtraFee:          in.ExtraFee,
	}
}

func toChannelAreaBizList(in []ShippingChannelArea) []*biz.ChannelArea {
	out := make([]*biz.ChannelArea, len(in))
	for i := range in {
		out[i] = toChannelAreaBiz(&in[i])
	}
	return out
}

func toLockerBiz(in *Locker) *biz.Locker {
	if in == nil {
		return nil
	}
	return &biz.Locker{
		ID:               in.ID,
		LockerCode:       in.LockerCode,
		LockerName:       in.LockerName,
		ProvinceCode:     in.ProvinceCode,
		ProvinceName:     in.ProvinceName,
		CityCode:         in.CityCode,
		CityName:         in.CityName,
		DistrictCode:     in.DistrictCode,
		DistrictName:     in.DistrictName,
		DetailAddress:    in.DetailAddress,
		Longitude:        in.Longitude,
		Latitude:         in.Latitude,
		BusinessStartTime: in.BusinessStartTime,
		BusinessEndTime:  in.BusinessEndTime,
		Status:           in.Status,
		CreatedAt:        in.CreatedAt,
		UpdatedAt:        in.UpdatedAt,
	}
}

func toLockerBizList(in []Locker) []*biz.Locker {
	out := make([]*biz.Locker, len(in))
	for i := range in {
		out[i] = toLockerBiz(&in[i])
	}
	return out
}

func toLockerBoxBiz(in *LockerBox) *biz.LockerBox {
	if in == nil {
		return nil
	}
	return &biz.LockerBox{
		ID:              in.ID,
		LockerID:        in.LockerID,
		BoxNo:           in.BoxNo,
		BoxType:         in.BoxType,
		BoxLength:       in.BoxLength,
		BoxWidth:        in.BoxWidth,
		BoxHeight:       in.BoxHeight,
		MaxWeight:       in.MaxWeight,
		BoxFee:          in.BoxFee,
		Status:          in.Status,
		ReservedOrderID: in.ReservedOrderID,
		ReservedAt:      in.ReservedAt,
		ExpireAt:        in.ExpireAt,
		Version:         in.Version,
		CreatedAt:       in.CreatedAt,
		UpdatedAt:       in.UpdatedAt,
	}
}

func toLockerBoxBizList(in []LockerBox) []*biz.LockerBox {
	out := make([]*biz.LockerBox, len(in))
	for i := range in {
		out[i] = toLockerBoxBiz(&in[i])
	}
	return out
}

func toServicePointBiz(in *ServicePoint) *biz.ServicePoint {
	if in == nil {
		return nil
	}
	return &biz.ServicePoint{
		ID:               in.ID,
		PointCode:        in.PointCode,
		PointName:        in.PointName,
		PointType:        in.PointType,
		ProvinceCode:     in.ProvinceCode,
		ProvinceName:     in.ProvinceName,
		CityCode:         in.CityCode,
		CityName:         in.CityName,
		DistrictCode:     in.DistrictCode,
		DistrictName:     in.DistrictName,
		DetailAddress:    in.DetailAddress,
		Longitude:        in.Longitude,
		Latitude:         in.Latitude,
		BusinessStartTime: in.BusinessStartTime,
		BusinessEndTime:  in.BusinessEndTime,
		ContactPhoneMask: in.ContactPhoneMask,
		Status:           in.Status,
		ExtraFee:         in.ExtraFee,
		CreatedAt:        in.CreatedAt,
		UpdatedAt:        in.UpdatedAt,
	}
}

func toServicePointBizList(in []ServicePoint) []*biz.ServicePoint {
	out := make([]*biz.ServicePoint, len(in))
	for i := range in {
		out[i] = toServicePointBiz(&in[i])
	}
	return out
}

func toPickupTimeSlotBiz(in *PickupTimeSlot) *biz.PickupTimeSlot {
	if in == nil {
		return nil
	}
	return &biz.PickupTimeSlot{
		ID:            in.ID,
		DistrictCode:  in.DistrictCode,
		PickupDate:    in.PickupDate,
		SlotCode:      in.SlotCode,
		StartTime:     in.StartTime,
		EndTime:       in.EndTime,
		Capacity:      in.Capacity,
		ReservedCount: in.ReservedCount,
		Status:        in.Status,
		Version:       in.Version,
		CreatedAt:     in.CreatedAt,
		UpdatedAt:     in.UpdatedAt,
	}
}

func toPickupTimeSlotBizList(in []PickupTimeSlot) []*biz.PickupTimeSlot {
	out := make([]*biz.PickupTimeSlot, len(in))
	for i := range in {
		out[i] = toPickupTimeSlotBiz(&in[i])
	}
	return out
}
