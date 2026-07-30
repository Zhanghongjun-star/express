package biz

import (
	"context"
	"time"

	v1 "shunfeng-miniprogram/api/shipping/v1"

	"github.com/go-kratos/kratos/v3/errors"
)

var (
	ErrChannelUnavailable     = errors.BadRequest(v1.ErrorReason_SHIPPING_CHANNEL_UNAVAILABLE.String(), "channel unavailable in current region")
	ErrLockerNotFound         = errors.NotFound(v1.ErrorReason_SHIPPING_NOT_FOUND.String(), "locker not found")
	ErrServicePointNotFound   = errors.NotFound(v1.ErrorReason_SHIPPING_NOT_FOUND.String(), "service point not found")
	ErrSlotFull               = errors.Conflict(v1.ErrorReason_SHIPPING_RESOURCE_OCCUPIED.String(), "time slot is fully occupied")
	ErrBoxTypeFull            = errors.Conflict(v1.ErrorReason_SHIPPING_RESOURCE_OCCUPIED.String(), "box type has no free compartments")
	ErrPickupDateOutOfRange   = errors.BadRequest(v1.ErrorReason_SHIPPING_PICKUP_DATE_OUT_OF_RANGE.String(), "pickup date must be within the next 3 days")
	ErrShippingInvalidArg     = errors.BadRequest(v1.ErrorReason_SHIPPING_INVALID_ARGUMENT.String(), "invalid shipping argument")
	ErrParamsMismatch         = errors.BadRequest(v1.ErrorReason_SHIPPING_PARAMS_MISMATCH.String(), "channel params mismatch")
	ErrShippingRateLimit      = errors.TooManyRequests(v1.ErrorReason_SHIPPING_RATE_LIMIT_EXCEEDED.String(), "rate limit exceeded, please retry later")
)

type ShippingChannel struct {
	ID              int64
	ChannelCode     string
	ChannelName     string
	ChannelDesc     string
	Status          int32
	BaseFee         float64
	NeedAddressGeo  bool
	NeedPickupSlot  bool
	NeedLockerBox   bool
	NeedServicePoint bool
	SortNo          int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ChannelArea struct {
	ID                int64
	ChannelCode       string
	ProvinceCode      string
	CityCode          string
	DistrictCode      string
	Status            int32
	UnavailableReason string
	ExtraFee          float64
}

type Locker struct {
	ID                int64
	LockerCode        string
	LockerName        string
	ProvinceCode      string
	ProvinceName      string
	CityCode          string
	CityName          string
	DistrictCode      string
	DistrictName      string
	DetailAddress     string
	Longitude         float64
	Latitude          float64
	BusinessStartTime string
	BusinessEndTime   string
	Status            int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type LockerBox struct {
	ID              int64
	LockerID        int64
	BoxNo           string
	BoxType         string
	BoxLength       int32
	BoxWidth        int32
	BoxHeight       int32
	MaxWeight       float64
	BoxFee          float64
	Status          int32
	ReservedOrderID int64
	ReservedAt      time.Time
	ExpireAt        time.Time
	Version         int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ServicePoint struct {
	ID               int64
	PointCode        string
	PointName        string
	PointType        int32
	ProvinceCode     string
	ProvinceName     string
	CityCode         string
	CityName         string
	DistrictCode     string
	DistrictName     string
	DetailAddress    string
	Longitude        float64
	Latitude         float64
	BusinessStartTime string
	BusinessEndTime  string
	ContactPhoneMask string
	Status           int32
	ExtraFee         float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type PickupTimeSlot struct {
	ID            int64
	DistrictCode  string
	PickupDate    string
	SlotCode      string
	StartTime     string
	EndTime       string
	Capacity      int32
	ReservedCount int32
	Status        int32
	Version       int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PickupReservation struct {
	ID           int64
	OrderID      int64
	UserID       int64
	SlotID       int64
	DistrictCode string
	PickupDate   string
	SlotCode     string
	Status       int32
	ReservedAt   time.Time
	ReleasedAt   time.Time
	ExpireAt     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ChannelRepo interface {
	ListChannels(ctx context.Context) ([]*ShippingChannel, error)
	FindChannelByCode(ctx context.Context, channelCode string) (*ShippingChannel, error)
	FindAreaByChannelDistrict(ctx context.Context, channelCode, districtCode string) (*ChannelArea, error)
	ListAreasByDistrict(ctx context.Context, districtCode string) ([]*ChannelArea, error)
}

type LockerRepo interface {
	ListByLocation(ctx context.Context, districtCode string, lat, lng float64, radius int32, offset, limit int) ([]*Locker, int, error)
	FindByID(ctx context.Context, id int64) (*Locker, error)
}

type LockerBoxRepo interface {
	GetBoxTypesByLocker(ctx context.Context, lockerID int64) ([]*LockerBox, error)
	LockBox(ctx context.Context, id int64, orderID int64, expireAt time.Time, version int64) error
	ReleaseBox(ctx context.Context, id int64, orderID int64, version int64) error
	ReleaseExpiredBoxes(ctx context.Context) (int64, error)
}

type ServicePointRepo interface {
	ListByLocation(ctx context.Context, pointType int32, districtCode string, lat, lng float64, offset, limit int) ([]*ServicePoint, int, error)
	FindByID(ctx context.Context, id int64) (*ServicePoint, error)
}

type PickupRepo interface {
	ListDates(ctx context.Context, districtCode string) ([]string, error)
	ListByDate(ctx context.Context, districtCode, pickupDate string) ([]*PickupTimeSlot, error)
	FindSlotByCode(ctx context.Context, districtCode, pickupDate, slotCode string) (*PickupTimeSlot, error)
	ReserveSlot(ctx context.Context, id int64, version int64) error
	ReleaseSlot(ctx context.Context, id int64, version int64) error
	ReleaseExpiredSlots(ctx context.Context) (int64, error)
}

type ChannelUsecase struct {
	channelRepo ChannelRepo
}

func NewChannelUsecase(repo ChannelRepo) *ChannelUsecase {
	return &ChannelUsecase{channelRepo: repo}
}

func (uc *ChannelUsecase) ListChannels(ctx context.Context, districtCode string) ([]*ShippingChannel, []*ChannelArea, error) {
	channels, err := uc.channelRepo.ListChannels(ctx)
	if err != nil {
		return nil, nil, err
	}
	areas, err := uc.channelRepo.ListAreasByDistrict(ctx, districtCode)
	if err != nil {
		return nil, nil, err
	}
	return channels, areas, nil
}

type LockerUsecase struct {
	lockerRepo    LockerRepo
	lockerBoxRepo LockerBoxRepo
}

func NewLockerUsecase(lockerRepo LockerRepo, lockerBoxRepo LockerBoxRepo) *LockerUsecase {
	return &LockerUsecase{
		lockerRepo:    lockerRepo,
		lockerBoxRepo: lockerBoxRepo,
	}
}

func (uc *LockerUsecase) ListLockers(ctx context.Context, districtCode string, lat, lng float64, radius int32, pageSize, pageOffset int) ([]*Locker, int, error) {
	if lat == 0 || lng == 0 {
		return nil, 0, ErrShippingInvalidArg
	}
	return uc.lockerRepo.ListByLocation(ctx, districtCode, lat, lng, radius, pageOffset, pageSize)
}

func (uc *LockerUsecase) ListBoxTypes(ctx context.Context, lockerID int64) ([]*LockerBox, error) {
	if lockerID <= 0 {
		return nil, ErrShippingInvalidArg
	}
	locker, err := uc.lockerRepo.FindByID(ctx, lockerID)
	if err != nil {
		return nil, err
	}
	if locker.Status != 1 {
		return nil, ErrLockerNotFound
	}
	return uc.lockerBoxRepo.GetBoxTypesByLocker(ctx, lockerID)
}

type ServicePointUsecase struct {
	repo ServicePointRepo
}

func NewServicePointUsecase(repo ServicePointRepo) *ServicePointUsecase {
	return &ServicePointUsecase{repo: repo}
}

func (uc *ServicePointUsecase) ListServicePoints(ctx context.Context, pointType int32, districtCode string, lat, lng float64, pageSize, pageOffset int) ([]*ServicePoint, int, error) {
	if lat == 0 || lng == 0 {
		return nil, 0, ErrShippingInvalidArg
	}
	return uc.repo.ListByLocation(ctx, pointType, districtCode, lat, lng, pageOffset, pageSize)
}

type PickupUsecase struct {
	repo PickupRepo
}

func NewPickupUsecase(repo PickupRepo) *PickupUsecase {
	return &PickupUsecase{repo: repo}
}

func (uc *PickupUsecase) ListPickupDates(ctx context.Context, districtCode string) ([]string, error) {
	if districtCode == "" {
		return nil, ErrShippingInvalidArg
	}
	return uc.repo.ListDates(ctx, districtCode)
}

func (uc *PickupUsecase) ListPickupTimeSlots(ctx context.Context, districtCode, pickupDate string) ([]*PickupTimeSlot, error) {
	if districtCode == "" || pickupDate == "" {
		return nil, ErrShippingInvalidArg
	}
	return uc.repo.ListByDate(ctx, districtCode, pickupDate)
}
