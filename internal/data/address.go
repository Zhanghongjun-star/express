package data

import (
	"cmp"
	"context"
	"slices"
	"sync"
	"time"

	v1 "shunfeng-miniprogram/api/address/v1"
	"shunfeng-miniprogram/internal/biz"
	"shunfeng-miniprogram/internal/model"
)

// addressRepo 地址仓储的内存实现。
type addressRepo struct {
	data *Data

	mu     sync.RWMutex
	nextID int64
	items  map[int64]*model.Address
}

// NewAddressRepo 创建地址仓储。
func NewAddressRepo(data *Data) biz.AddressRepo {
	return &addressRepo{
		data:   data,
		nextID: 1,
		items:  make(map[int64]*model.Address),
	}
}

func (r *addressRepo) FindByID(_ context.Context, userID, id int64) (*biz.Address, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok || item.DelFlag || item.UserID != userID {
		return nil, biz.ErrAddressNotFound
	}
	return toBiz(item), nil
}

func (r *addressRepo) ListAddresses(_ context.Context, userID int64, opts ...biz.AddressListOption) ([]*biz.Address, error) {
	options := biz.AddressListOptions{Limit: 20}
	for _, opt := range opts {
		opt(&options)
	}
	if options.Offset < 0 || options.Limit <= 0 {
		return nil, biz.ErrAddressInvalidArgument
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	records := make([]*biz.Address, 0, len(r.items))
	for _, item := range r.items {
		if item.DelFlag || item.UserID != userID {
			continue
		}
		if options.AddrType != v1.AddressType_ADDRESS_TYPE_UNSPECIFIED && item.AddrType != options.AddrType {
			continue
		}
		records = append(records, toBiz(item))
	}
	slices.SortFunc(records, func(a, b *biz.Address) int {
		return cmp.Compare(a.ID, b.ID)
	})

	if options.Offset >= len(records) {
		return []*biz.Address{}, nil
	}
	end := options.Offset + options.Limit
	if end > len(records) {
		end = len(records)
	}
	return records[options.Offset:end], nil
}

func (r *addressRepo) CreateAddress(_ context.Context, address *biz.Address) (*biz.Address, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	addr := newAddress(address)
	addr.ID = r.nextID
	addr.CreateTime = now
	addr.UpdateTime = now
	r.items[addr.ID] = addr
	r.nextID++
	return toBiz(addr), nil
}

func (r *addressRepo) UpdateAddress(_ context.Context, address *biz.Address) (*biz.Address, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.items[address.ID]
	if !ok || current.DelFlag || current.UserID != address.UserID {
		return nil, biz.ErrAddressNotFound
	}

	updated := newAddress(address)
	updated.CreateTime = current.CreateTime
	updated.UpdateTime = time.Now()
	updated.DelFlag = current.DelFlag
	r.items[address.ID] = updated
	return toBiz(updated), nil
}

func (r *addressRepo) DeleteAddress(_ context.Context, userID, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.items[id]
	if !ok || current.DelFlag || current.UserID != userID {
		return biz.ErrAddressNotFound
	}
	current.DelFlag = true
	current.UpdateTime = time.Now()
	return nil
}

func (r *addressRepo) BatchDeleteAddresses(_ context.Context, userID int64, ids []int64) (int32, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var deleted int32
	for _, id := range ids {
		current, ok := r.items[id]
		if !ok || current.DelFlag || current.UserID != userID {
			continue
		}
		current.DelFlag = true
		current.UpdateTime = time.Now()
		deleted++
	}
	return deleted, nil
}

// newAddress 将 biz.Address 转换为 model.Address（DO → PO）。
func newAddress(in *biz.Address) *model.Address {
	if in == nil {
		return nil
	}
	return &model.Address{
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

// toBiz 将 model.Address 转换为 biz.Address（PO → DO）。
func toBiz(in *model.Address) *biz.Address {
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
