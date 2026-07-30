package data

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"shunfeng-miniprogram/internal/biz"
	"shunfeng-miniprogram/internal/pkg/amap"
)

// addressRepo 地址仓储的 GORM 实现。
type addressRepo struct{}

// NewAddressRepo 创建地址仓储。
func NewAddressRepo() biz.AddressRepo {
	return &addressRepo{}
}

func (r *addressRepo) FindByID(ctx context.Context, userID, id int64) (*biz.Address, error) {
	var po Address
	if err := DB.WithContext(ctx).Where("id = ? AND user_id = ? AND del_flag = 0", id, userID).First(&po).Error; err != nil {
		return nil, biz.ErrAddressNotFound
	}
	return toAddressBiz(&po), nil
}

func (r *addressRepo) ListAddresses(ctx context.Context, userID int64, opts ...biz.AddressListOption) ([]*biz.Address, error) {
	options := biz.AddressListOptions{Limit: 20}
	for _, opt := range opts {
		opt(&options)
	}
	if options.Offset < 0 || options.Limit <= 0 {
		return nil, biz.ErrAddressInvalidArgument
	}

	tx := DB.WithContext(ctx).Where("user_id = ? AND del_flag = 0", userID)
	if options.AddrType != 0 {
		tx = tx.Where("addr_type = ?", options.AddrType)
	}
	var pos []Address
	if err := tx.Order("id ASC").Offset(options.Offset).Limit(options.Limit).Find(&pos).Error; err != nil {
		return nil, err
	}

	addresses := make([]*biz.Address, len(pos))
	for i, po := range pos {
		addresses[i] = toAddressBiz(&po)
	}
	return addresses, nil
}

func (r *addressRepo) CreateAddress(ctx context.Context, address *biz.Address) (*biz.Address, error) {
	po := newAddressPO(address)
	po.CreateTime = time.Now()
	po.UpdateTime = time.Now()

	if err := DB.WithContext(ctx).Create(po).Error; err != nil {
		return nil, err
	}
	return toAddressBiz(po), nil
}

func (r *addressRepo) UpdateAddress(ctx context.Context, address *biz.Address) (*biz.Address, error) {
	var po Address
	if err := DB.WithContext(ctx).Where("id = ? AND user_id = ? AND del_flag = 0", address.ID, address.UserID).First(&po).Error; err != nil {
		return nil, biz.ErrAddressNotFound
	}

	updated := newAddressPO(address)
	updated.ID = po.ID
	updated.CreateTime = po.CreateTime
	updated.UpdateTime = time.Now()
	updated.DelFlag = po.DelFlag

	if err := DB.WithContext(ctx).Model(&po).Updates(updated).Error; err != nil {
		return nil, err
	}
	return toAddressBiz(updated), nil
}

func (r *addressRepo) DeleteAddress(ctx context.Context, userID, id int64) error {
	result := DB.WithContext(ctx).Model(&Address{}).
		Where("id = ? AND user_id = ? AND del_flag = 0", id, userID).
		Update("del_flag", 1)
	if result.RowsAffected == 0 {
		return biz.ErrAddressNotFound
	}
	return result.Error
}

func (r *addressRepo) BatchDeleteAddresses(ctx context.Context, userID int64, ids []int64) (int32, error) {
	result := DB.WithContext(ctx).Model(&Address{}).
		Where("id IN ? AND user_id = ? AND del_flag = 0", ids, userID).
		Update("del_flag", 1)
	return int32(result.RowsAffected), result.Error
}

// newAddressPO 将 biz.Address 转换为 Address（DO → PO）。
func newAddressPO(in *biz.Address) *Address {
	if in == nil {
		return nil
	}
	return &Address{
		UserID:        in.UserID,
		AddrType:      in.AddrType,
		ReceiverName:  in.ReceiverName,
		ReceiverPhone: in.ReceiverPhone,
		Province:      in.Province,
		City:          in.City,
		District:      in.District,
		DetailAddr:    in.DetailAddr,
		IsDefault:     in.IsDefault,
		Latitude:      in.Latitude,
		Longitude:     in.Longitude,
	}
}

// toAddressBiz 将 Address 转换为 biz.Address（PO → DO）。
func toAddressBiz(in *Address) *biz.Address {
	if in == nil {
		return nil
	}
	return &biz.Address{
		ID:            in.ID,
		UserID:        in.UserID,
		AddrType:      in.AddrType,
		ReceiverName:  in.ReceiverName,
		ReceiverPhone: in.ReceiverPhone,
		Province:      in.Province,
		City:          in.City,
		District:      in.District,
		DetailAddr:    in.DetailAddr,
		IsDefault:     in.IsDefault,
		Latitude:      in.Latitude,
		Longitude:     in.Longitude,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
		DelFlag:       in.DelFlag,
	}
}

// amapGeocoder 基于高德地图 SDK 实现 biz.Geocoder，将文本地址转换为经纬度。
type amapGeocoder struct {
	client *amap.Client
}

// NewAmapGeocoder 从环境变量 AMAP_WEB_KEY 读取高德 Key 构造 Geocoder。
// 若未配置 Key，返回一个 noop 实现：反查始终失败，但不影响地址的创建/更新。
func NewAmapGeocoder() biz.Geocoder {
	key := os.Getenv("AMAP_WEB_KEY")
	if key == "" {
		return noopGeocoder{}
	}
	return &amapGeocoder{client: amap.NewClient(key)}
}

// GeocodeAddress 拼接省/市/区/详细地址后调用高德地理编码，返回纬度 lat、经度 lng。
func (g *amapGeocoder) GeocodeAddress(ctx context.Context, province, city, district, detail string) (float64, float64, error) {
	full := strings.Join([]string{province, city, district, detail}, "")
	if strings.TrimSpace(full) == "" {
		return 0, 0, fmt.Errorf("empty address")
	}
	opts := []amap.GeocodeOption{}
	if city != "" {
		opts = append(opts, amap.GeocodeWithCity(city))
	}
	resp, err := g.client.Geocode(ctx, full, opts...)
	if err != nil {
		return 0, 0, err
	}
	if len(resp.Geocodes) == 0 {
		return 0, 0, fmt.Errorf("no geocode result for %q", full)
	}
	// 高德 Location 格式为 "经度,纬度"。
	parts := strings.Split(string(resp.Geocodes[0].Location), ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid location: %q", resp.Geocodes[0].Location)
	}
	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse lng: %w", err)
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse lat: %w", err)
	}
	return lat, lng, nil
}

// noopGeocoder 未配置高德 Key 时的兜底实现。
type noopGeocoder struct{}

func (noopGeocoder) GeocodeAddress(ctx context.Context, province, city, district, detail string) (float64, float64, error) {
	return 0, 0, fmt.Errorf("geocoder not configured (AMAP_WEB_KEY missing)")
}
