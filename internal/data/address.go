package data

import (
	"context"
	"sync"
	"time"

	"shunfeng-miniprogram/internal/biz"
)

type addressRepo struct {
	data *Data

	mu     sync.RWMutex
	nextID int64
	addrs  map[int64]*biz.Address
}

func NewAddressRepo(data *Data) biz.AddressRepo {
	return &addressRepo{
		data:   data,
		nextID: 1,
		addrs:  make(map[int64]*biz.Address),
	}
}

func (r *addressRepo) ListByUser(_ context.Context, userID int64) ([]*biz.Address, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*biz.Address, 0)
	for _, addr := range r.addrs {
		if addr.UserID == userID && addr.DelFlag == 0 {
			result = append(result, addr)
		}
	}
	return result, nil
}

func (r *addressRepo) Create(_ context.Context, addr *biz.Address) (*biz.Address, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	addr.AddrID = r.nextID
	addr.CreateTime = now
	addr.DelFlag = 0
	r.addrs[addr.AddrID] = addr
	r.nextID++
	return addr, nil
}

func (r *addressRepo) Update(_ context.Context, addr *biz.Address) (*biz.Address, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cur, ok := r.addrs[addr.AddrID]
	if !ok || cur.DelFlag != 0 {
		return nil, biz.ErrAddressNotFound
	}
	addr.CreateTime = cur.CreateTime
	r.addrs[addr.AddrID] = addr
	return addr, nil
}

func (r *addressRepo) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cur, ok := r.addrs[id]
	if !ok || cur.DelFlag != 0 {
		return biz.ErrAddressNotFound
	}
	cur.DelFlag = 1
	return nil
}

func (r *addressRepo) BatchDelete(_ context.Context, ids []int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, id := range ids {
		cur, ok := r.addrs[id]
		if ok && cur.DelFlag == 0 {
			cur.DelFlag = 1
		}
	}
	return nil
}
