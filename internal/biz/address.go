package biz

import (
	"context"
	"strings"
	"time"

	v1 "shunfeng-miniprogram/api/address/v1"

	"github.com/go-kratos/kratos/v3/errors"
)

const maxBatchDeleteAddresses = 50

var (
	// ErrAddressNotFound 地址不存在时返回。
	ErrAddressNotFound = errors.NotFound(v1.ErrorReason_ADDRESS_NOT_FOUND.String(), "address not found")
	// ErrAddressInvalidArgument 地址请求参数无效时返回。
	ErrAddressInvalidArgument = errors.BadRequest(v1.ErrorReason_ADDRESS_INVALID_ARGUMENT.String(), "invalid address argument")
	// ErrAddressParseFailed 地址解析失败时返回。
	ErrAddressParseFailed = errors.BadRequest(v1.ErrorReason_ADDRESS_PARSE_FAILED.String(), "address parse failed")
)

// Address 地址领域对象。
type Address struct {
	ID           int64
	UserID       int64
	AddrType     v1.AddressType
	ReceiverName string
	ReceiverPhone string
	Province     string
	City         string
	District     string
	DetailAddr   string
	IsDefault    bool
	CreateTime   time.Time
	UpdateTime   time.Time
	DelFlag      bool
}

// AddressRepo 地址仓储接口。
type AddressRepo interface {
	FindByID(context.Context, int64, int64) (*Address, error)
	ListAddresses(context.Context, int64, ...AddressListOption) ([]*Address, error)
	CreateAddress(context.Context, *Address) (*Address, error)
	UpdateAddress(context.Context, *Address) (*Address, error)
	DeleteAddress(context.Context, int64, int64) error
	BatchDeleteAddresses(context.Context, int64, []int64) (int32, error)
}

// AddressListOption 地址列表查询选项函数。
type AddressListOption func(*AddressListOptions)

// AddressListOptions 地址列表查询选项。
type AddressListOptions struct {
	AddrType v1.AddressType
	Offset   int
	Limit    int
}

// AddressListType 按地址类型筛选。
func AddressListType(addrType v1.AddressType) AddressListOption {
	return func(o *AddressListOptions) { o.AddrType = addrType }
}

// AddressListOffset 设置偏移量。
func AddressListOffset(offset int) AddressListOption {
	return func(o *AddressListOptions) { o.Offset = offset }
}

// AddressListLimit 设置限制条数。
func AddressListLimit(limit int) AddressListOption {
	return func(o *AddressListOptions) { o.Limit = limit }
}

// AddressUsecase 地址用例。
type AddressUsecase struct {
	repo AddressRepo
}

// NewAddressUsecase 创建地址用例。
func NewAddressUsecase(repo AddressRepo) *AddressUsecase {
	return &AddressUsecase{repo: repo}
}

// GetAddress 根据 ID 获取地址。
func (uc *AddressUsecase) GetAddress(ctx context.Context, userID, id int64) (*Address, error) {
	if userID <= 0 || id <= 0 {
		return nil, ErrAddressInvalidArgument
	}
	return uc.repo.FindByID(ctx, userID, id)
}

// ListAddresses 获取地址列表。
func (uc *AddressUsecase) ListAddresses(ctx context.Context, userID int64, opts ...AddressListOption) ([]*Address, error) {
	if userID <= 0 {
		return nil, ErrAddressInvalidArgument
	}
	return uc.repo.ListAddresses(ctx, userID, opts...)
}

// CreateAddress 创建地址。
func (uc *AddressUsecase) CreateAddress(ctx context.Context, address *Address) (*Address, error) {
	if err := validateAddress(address, false); err != nil {
		return nil, err
	}
	return uc.repo.CreateAddress(ctx, address)
}

// UpdateAddress 更新地址。
func (uc *AddressUsecase) UpdateAddress(ctx context.Context, address *Address) (*Address, error) {
	if err := validateAddress(address, true); err != nil {
		return nil, err
	}
	return uc.repo.UpdateAddress(ctx, address)
}

// DeleteAddress 删除地址。
func (uc *AddressUsecase) DeleteAddress(ctx context.Context, userID, id int64) error {
	if userID <= 0 || id <= 0 {
		return ErrAddressInvalidArgument
	}
	return uc.repo.DeleteAddress(ctx, userID, id)
}

// BatchDeleteAddresses 批量删除地址。
func (uc *AddressUsecase) BatchDeleteAddresses(ctx context.Context, userID int64, ids []int64) (int32, error) {
	if userID <= 0 || len(ids) == 0 || len(ids) > maxBatchDeleteAddresses {
		return 0, ErrAddressInvalidArgument
	}
	return uc.repo.BatchDeleteAddresses(ctx, userID, ids)
}

// ParseAddress 解析文本内容为地址对象。
func (uc *AddressUsecase) ParseAddress(content string) (*Address, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, ErrAddressInvalidArgument
	}

	fields := strings.FieldsFunc(content, func(r rune) bool {
		return r == ',' || r == '，' || r == ';' || r == '；' || r == '\n'
	})
	if len(fields) < 2 {
		return nil, ErrAddressParseFailed
	}

	address := &Address{
		ReceiverName: strings.TrimSpace(fields[0]),
		DetailAddr:   strings.TrimSpace(fields[len(fields)-1]),
		AddrType:     v1.AddressType_ADDRESS_TYPE_MAINLAND,
	}

	if len(fields) >= 3 {
		address.Province = strings.TrimSpace(fields[1])
	}
	if len(fields) >= 4 {
		address.City = strings.TrimSpace(fields[2])
	}
	if len(fields) >= 5 {
		address.District = strings.TrimSpace(fields[3])
	}
	for _, field := range fields {
		digits := onlyDigits(field)
		if len(digits) >= 7 {
			address.ReceiverPhone = digits
			break
		}
	}
	if address.ReceiverPhone == "" {
		return nil, ErrAddressParseFailed
	}
	return address, nil
}

// validateAddress 校验地址字段合法性。
func validateAddress(address *Address, requireID bool) error {
	if address == nil {
		return ErrAddressInvalidArgument
	}
	if requireID && address.ID <= 0 {
		return ErrAddressInvalidArgument
	}
	if address.UserID <= 0 {
		return ErrAddressInvalidArgument
	}
	if strings.TrimSpace(address.ReceiverName) == "" || strings.TrimSpace(address.ReceiverPhone) == "" || strings.TrimSpace(address.DetailAddr) == "" {
		return ErrAddressInvalidArgument
	}
	switch address.AddrType {
	case v1.AddressType_ADDRESS_TYPE_MAINLAND, v1.AddressType_ADDRESS_TYPE_HK_MO_TW, v1.AddressType_ADDRESS_TYPE_INTERNATIONAL:
	default:
		return ErrAddressInvalidArgument
	}
	return nil
}

// onlyDigits 提取字符串中的数字部分。
func onlyDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
