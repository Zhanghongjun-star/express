package data

import (
	"context"
	"time"

	"shunfeng-miniprogram/internal/biz"
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
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
		DelFlag:       in.DelFlag,
	}
}
