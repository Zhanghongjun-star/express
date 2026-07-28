package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v3/errors"
)

type AddressType int32

const (
	AddrTypeMainland AddressType = 1
	AddrTypeHKMT     AddressType = 2
	AddrTypeIntl     AddressType = 3
)

type Address struct {
	AddrID        int64
	UserID        int64
	AddrType      AddressType
	ReceiverName  string
	ReceiverPhone string
	Province      string
	City          string
	District      string
	DetailAddr    string
	CreateTime    time.Time
	DelFlag       int32
}

var (
	ErrAddressNotFound   = errors.NotFound("ADDRESS_NOT_FOUND", "address not found")
	ErrAddressInvalidArg = errors.BadRequest("ADDRESS_INVALID_ARGUMENT", "invalid address argument")
)

type AddressRepo interface {
	ListByUser(ctx context.Context, userID int64) ([]*Address, error)
	Create(ctx context.Context, addr *Address) (*Address, error)
	Update(ctx context.Context, addr *Address) (*Address, error)
	Delete(ctx context.Context, id int64) error
	BatchDelete(ctx context.Context, ids []int64) error
}

type AddressUsecase struct {
	repo AddressRepo
}

func NewAddressUsecase(repo AddressRepo) *AddressUsecase {
	return &AddressUsecase{repo: repo}
}

func (uc *AddressUsecase) List(ctx context.Context, userID int64) ([]*Address, error) {
	if userID <= 0 {
		return nil, ErrAddressInvalidArg
	}
	return uc.repo.ListByUser(ctx, userID)
}

func (uc *AddressUsecase) Create(ctx context.Context, addr *Address) (*Address, error) {
	if addr.ReceiverName == "" || addr.ReceiverPhone == "" {
		return nil, ErrAddressInvalidArg
	}
	return uc.repo.Create(ctx, addr)
}

func (uc *AddressUsecase) Update(ctx context.Context, addr *Address) (*Address, error) {
	if addr.AddrID <= 0 {
		return nil, ErrAddressInvalidArg
	}
	return uc.repo.Update(ctx, addr)
}

func (uc *AddressUsecase) Delete(ctx context.Context, id int64, userID int64) error {
	if id <= 0 || userID <= 0 {
		return ErrAddressInvalidArg
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *AddressUsecase) BatchDelete(ctx context.Context, ids []int64, userID int64) error {
	if len(ids) == 0 || len(ids) > 50 || userID <= 0 {
		return ErrAddressInvalidArg
	}
	return uc.repo.BatchDelete(ctx, ids)
}
