package biz

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	v3 "shunfeng-miniprogram/api/order/v3"

	"github.com/go-kratos/kratos/v3/errors"
)

var (
	ErrOrderCreateFailed           = errors.InternalServer(v3.ErrorReason_ORDER_CREATE_FAILED.String(), "order creation failed")
	ErrOrderPaymentRequired        = errors.Forbidden(v3.ErrorReason_ORDER_PAYMENT_REQUIRED.String(), "payment required")
	ErrOrderChannelUnavail         = errors.BadRequest(v3.ErrorReason_ORDER_CHANNEL_UNAVAILABLE.String(), "shipping channel unavailable")
	ErrOrderInvalidStatusTransition = errors.BadRequest(v3.ErrorReason_ORDER_CREATE_FAILED.String(), "invalid order status transition")
)

const (
	OrderStatusPendingPickup = "pending_pickup"
	OrderStatusAccepted      = "accepted"
	OrderStatusAwaitingPickup = "awaiting_pickup"
	OrderStatusPickedUp      = "picked_up"
	OrderStatusInTransit     = "in_transit"
	OrderStatusDelivering    = "delivering"
	OrderStatusSigned        = "signed"
	OrderStatusCancelled     = "cancelled"

	PaymentMethodSenderPay   = "sender_pay"
	PaymentMethodReceiverPay = "receiver_pay"
	PaymentMethodMonthly     = "monthly_settle"

	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
	PaymentStatusExpired = "expired"
)

type ExpressOrder struct {
	ID              int64
	UserID          int64
	OrderNo         string
	ExpressNo       string

	SenderName      string
	SenderPhone     string
	SenderProvince  string
	SenderCity      string
	SenderDistrict  string
	SenderDetail    string
	SenderLat       float64
	SenderLng       float64

	ReceiverName    string
	ReceiverPhone   string
	ReceiverProvince string
	ReceiverCity    string
	ReceiverDistrict string
	ReceiverDetail  string
	ReceiverLat     float64
	ReceiverLng     float64

	Weight          float64
	Length          int32
	Width           int32
	Height          int32

	BaseFreight     float64
	InsureFee       float64
	TotalFreight    float64

	ChannelCode     string
	LockerID        int64
	BoxType         string
	ServicePointID  int64
	PickupDate      string
	PickupSlotCode  string
	PickupStartTime string
	PickupEndTime   string

	PaymentMethod   string
	PaymentStatus   string
	PaidAt          time.Time

	PrivacyProtection bool

	Status          string
	DelFlag         bool
	Tag             string
	IsFollowed      bool
	ShareURL        string
	PayURL          string
	TradeNo         string

	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CreateExpressOrderCommand struct {
	UserID          int64

	SenderName      string
	SenderPhone     string
	SenderProvince  string
	SenderCity      string
	SenderDistrict  string
	SenderDetail    string

	ReceiverName    string
	ReceiverPhone   string
	ReceiverProvince string
	ReceiverCity    string
	ReceiverDistrict string
	ReceiverDetail  string

	Weight          float64
	Length          int32
	Width           int32
	Height          int32

	ChannelCode     string
	LockerID        int64
	BoxType         string
	ServicePointID  int64
	PickupDate      string
	PickupSlotCode  string

	PaymentMethod   string
	PrivacyProtection bool
	ExpressType     int32
	InsureValue     float64

	EstimatedFreight float64
	InsureFee        float64
}

type ExpressOrderRepo interface {
	Create(context.Context, *ExpressOrder) (*ExpressOrder, error)
	FindByID(context.Context, int64, int64) (*ExpressOrder, error)
	Search(context.Context, int64, string, int, int) ([]*ExpressOrder, int32, error)
	ListByCategory(context.Context, int64, string, int, int) ([]*ExpressOrder, int32, error)
	Delete(context.Context, int64, int64) error
	UpdateTag(context.Context, int64, int64, string) error
	UpdateFollow(context.Context, int64, int64, bool) error
	UpdateShareURL(context.Context, int64, int64, string) error
	UpdatePayURL(context.Context, int64, int64, string, string) error
	UpdatePaymentStatusByTradeNo(context.Context, string, string) error
	UpdateStatus(context.Context, int64, int64, string) error
	UpdatePaymentStatus(context.Context, int64, int64, string) error
}

type ExpressOrderUsecase struct {
	repo        ExpressOrderRepo
	freightUc   *FreightUsecase
	channelUc   *ChannelUsecase
	lockerUc    *LockerUsecase
	pickupUc    *PickupUsecase
	servicePtUc *ServicePointUsecase
	userRepo    UserRepo
}

func NewExpressOrderUsecase(repo ExpressOrderRepo, freightUc *FreightUsecase, channelUc *ChannelUsecase, lockerUc *LockerUsecase, pickupUc *PickupUsecase, servicePtUc *ServicePointUsecase, userRepo UserRepo) *ExpressOrderUsecase {
	return &ExpressOrderUsecase{
		repo:        repo,
		freightUc:   freightUc,
		channelUc:   channelUc,
		lockerUc:    lockerUc,
		pickupUc:    pickupUc,
		servicePtUc: servicePtUc,
		userRepo:    userRepo,
	}
}

// validTransitions 定义订单状态的合法转换路径。
// key 为当前状态，value 为允许的下一个状态集合。
var validTransitions = map[string]map[string]bool{
	OrderStatusPendingPickup:  {OrderStatusAccepted: true, OrderStatusCancelled: true},
	OrderStatusAccepted:       {OrderStatusAwaitingPickup: true, OrderStatusCancelled: true},
	OrderStatusAwaitingPickup: {OrderStatusPickedUp: true, OrderStatusCancelled: true},
	OrderStatusPickedUp:       {OrderStatusInTransit: true, OrderStatusCancelled: true},
	OrderStatusInTransit:      {OrderStatusDelivering: true},
	OrderStatusDelivering:     {OrderStatusSigned: true},
}

// ValidateStatusTransition 校验订单状态变更是否合法。
// cancelled 和 signed 为终态，不接受任何后续转换。
func ValidateStatusTransition(current, next string) bool {
	if current == OrderStatusCancelled || current == OrderStatusSigned {
		return next == current
	}
	allowed, ok := validTransitions[current]
	if !ok {
		return false
	}
	return allowed[next]
}

func (uc *ExpressOrderUsecase) CreateOrder(ctx context.Context, cmd *CreateExpressOrderCommand) (*ExpressOrder, error) {
	if err := uc.validateCreateOrder(cmd); err != nil {
		return nil, err
	}

	if cmd.PaymentMethod == PaymentMethodMonthly {
		user, err := uc.userRepo.FindByID(ctx, cmd.UserID)
		if err != nil || !user.IsEnterprise {
			return nil, ErrOrderChannelUnavail
		}
	}

	if cmd.ChannelCode == "DOOR_PICKUP" && cmd.PickupDate != "" && cmd.PickupSlotCode != "" {
		slots, err := uc.pickupUc.ListPickupTimeSlots(ctx, cmd.SenderDistrict, cmd.PickupDate)
		if err != nil {
			return nil, ErrOrderChannelUnavail
		}
		found := false
		for _, sl := range slots {
			if sl.SlotCode == cmd.PickupSlotCode && sl.ReservedCount < sl.Capacity {
				found = true
				cmd.PickupDate = sl.PickupDate
				break
			}
		}
		if !found {
			return nil, ErrSlotFull
		}
	}

	if cmd.EstimatedFreight <= 0 {
		req := &FreightRequest{
			SenderProvince:   cmd.SenderProvince,
			SenderCity:       cmd.SenderCity,
			SenderArea:       cmd.SenderDistrict,
			ReceiverProvince: cmd.ReceiverProvince,
			ReceiverCity:     cmd.ReceiverCity,
			ReceiverArea:     cmd.ReceiverDistrict,
			Weight:           cmd.Weight,
			Length:           cmd.Length,
			Width:            cmd.Width,
			Height:           cmd.Height,
			ExpressType:      ExpressType(cmd.ExpressType),
			InsureValue:      cmd.InsureValue,
		}
		result, err := uc.freightUc.Estimate(ctx, req)
		if err != nil {
			return nil, err
		}
		cmd.EstimatedFreight = result.TotalPrice
		cmd.InsureFee = result.InsureFee
	}

	order := &ExpressOrder{
		UserID:          cmd.UserID,
		OrderNo:         generateOrderNo(),
		SenderName:      cmd.SenderName,
		SenderPhone:     cmd.SenderPhone,
		SenderProvince:  cmd.SenderProvince,
		SenderCity:      cmd.SenderCity,
		SenderDistrict:  cmd.SenderDistrict,
		SenderDetail:    cmd.SenderDetail,
		ReceiverName:    cmd.ReceiverName,
		ReceiverPhone:   cmd.ReceiverPhone,
		ReceiverProvince: cmd.ReceiverProvince,
		ReceiverCity:    cmd.ReceiverCity,
		ReceiverDistrict: cmd.ReceiverDistrict,
		ReceiverDetail:  cmd.ReceiverDetail,
		Weight:          cmd.Weight,
		Length:          cmd.Length,
		Width:           cmd.Width,
		Height:          cmd.Height,
		BaseFreight:     cmd.EstimatedFreight - cmd.InsureFee,
		InsureFee:       cmd.InsureFee,
		TotalFreight:    cmd.EstimatedFreight,
		ChannelCode:     cmd.ChannelCode,
		LockerID:        cmd.LockerID,
		BoxType:         cmd.BoxType,
		ServicePointID:  cmd.ServicePointID,
		PickupDate:      cmd.PickupDate,
		PickupSlotCode:  cmd.PickupSlotCode,
		PaymentMethod:   cmd.PaymentMethod,
		PaymentStatus:   PaymentStatusPending,
		PrivacyProtection: cmd.PrivacyProtection,
		Status:          OrderStatusPendingPickup,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	created, err := uc.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *ExpressOrderUsecase) GetOrder(ctx context.Context, userID, orderID int64) (*ExpressOrder, error) {
	if userID <= 0 || orderID <= 0 {
		return nil, ErrOrderInvalidArgument
	}
	return uc.repo.FindByID(ctx, userID, orderID)
}

func (uc *ExpressOrderUsecase) SearchOrders(ctx context.Context, userID int64, keyword string, offset, limit int) ([]*ExpressOrder, int32, error) {
	if userID <= 0 {
		return nil, 0, ErrOrderInvalidArgument
	}
	if limit <= 0 || limit > 20 {
		limit = 20
	}
	return uc.repo.Search(ctx, userID, strings.TrimSpace(keyword), offset, limit)
}

func (uc *ExpressOrderUsecase) ListByCategory(ctx context.Context, userID int64, category string, offset, limit int) ([]*ExpressOrder, int32, error) {
	if userID <= 0 {
		return nil, 0, ErrOrderInvalidArgument
	}
	validCategories := map[string]bool{"sent": true, "received": true, "followed": true, "unpaid": true}
	if !validCategories[category] {
		return nil, 0, ErrOrderInvalidArgument
	}
	if limit <= 0 || limit > 20 {
		limit = 20
	}
	return uc.repo.ListByCategory(ctx, userID, category, offset, limit)
}

func (uc *ExpressOrderUsecase) DeleteOrder(ctx context.Context, userID, orderID int64) error {
	if userID <= 0 || orderID <= 0 {
		return ErrOrderInvalidArgument
	}
	return uc.repo.Delete(ctx, userID, orderID)
}

func (uc *ExpressOrderUsecase) SetOrderTag(ctx context.Context, userID, orderID int64, tag string) (*ExpressOrder, error) {
	if userID <= 0 || orderID <= 0 {
		return nil, ErrOrderInvalidArgument
	}
	if err := uc.repo.UpdateTag(ctx, userID, orderID, strings.TrimSpace(tag)); err != nil {
		return nil, err
	}
	return uc.repo.FindByID(ctx, userID, orderID)
}

func (uc *ExpressOrderUsecase) FollowOrder(ctx context.Context, userID, orderID int64, follow bool) (*ExpressOrder, error) {
	if userID <= 0 || orderID <= 0 {
		return nil, ErrOrderInvalidArgument
	}
	if err := uc.repo.UpdateFollow(ctx, userID, orderID, follow); err != nil {
		return nil, err
	}
	return uc.repo.FindByID(ctx, userID, orderID)
}

func (uc *ExpressOrderUsecase) ShareOrder(ctx context.Context, userID, orderID int64) (string, error) {
	if userID <= 0 || orderID <= 0 {
		return "", ErrOrderInvalidArgument
	}
	order, err := uc.repo.FindByID(ctx, userID, orderID)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://example.com/share/%s", order.OrderNo)
	if err := uc.repo.UpdateShareURL(ctx, userID, orderID, url); err != nil {
		return "", err
	}
	return url, nil
}

func (uc *ExpressOrderUsecase) PayOrder(ctx context.Context, userID, orderID int64) (payURL, tradeNo string, err error) {
	if userID <= 0 || orderID <= 0 {
		return "", "", ErrOrderInvalidArgument
	}
	order, err := uc.repo.FindByID(ctx, userID, orderID)
	if err != nil {
		return "", "", err
	}
	if order.PaymentStatus == PaymentStatusPaid {
		return order.PayURL, order.TradeNo, nil
	}
	if order.PaymentStatus != PaymentStatusPending {
		return "", "", ErrOrderPaymentRequired
	}
	if order.PayURL != "" {
		return order.PayURL, order.TradeNo, nil
	}
	b := make([]byte, 4)
	rand.Read(b)
	tradeNo = fmt.Sprintf("T%s%s", time.Now().Format("20060102150405"), hex.EncodeToString(b))
	payURL = fmt.Sprintf("https://pay.example.com/pay?trade_no=%s&amount=%.2f", tradeNo, order.TotalFreight)
	if err := uc.repo.UpdatePayURL(ctx, userID, orderID, payURL, tradeNo); err != nil {
		return "", "", err
	}
	return payURL, tradeNo, nil
}

func (uc *ExpressOrderUsecase) HandlePaymentCallback(ctx context.Context, tradeNo, orderNo, paymentStatus, signature string) error {
	if strings.TrimSpace(tradeNo) == "" && strings.TrimSpace(orderNo) == "" {
		return ErrOrderInvalidArgument
	}
	if paymentStatus != PaymentStatusPaid {
		return nil
	}
	if strings.TrimSpace(tradeNo) == "" {
		return nil
	}
	return uc.repo.UpdatePaymentStatusByTradeNo(ctx, tradeNo, PaymentStatusPaid)
}

func (uc *ExpressOrderUsecase) validateCreateOrder(cmd *CreateExpressOrderCommand) error {
	if cmd.UserID <= 0 {
		return ErrOrderInvalidArgument
	}
	if strings.TrimSpace(cmd.SenderName) == "" || strings.TrimSpace(cmd.SenderPhone) == "" {
		return ErrOrderInvalidArgument
	}
	if strings.TrimSpace(cmd.ReceiverName) == "" || strings.TrimSpace(cmd.ReceiverPhone) == "" {
		return ErrOrderInvalidArgument
	}
	if cmd.Weight <= 0 {
		return ErrOrderInvalidArgument
	}
	validChannels := map[string]bool{"DOOR_PICKUP": true, "LOCKER": true, "SF_STATION": true, "PARTNER_STORE": true}
	if !validChannels[cmd.ChannelCode] {
		return ErrOrderInvalidArgument
	}
	validPayments := map[string]bool{PaymentMethodSenderPay: true, PaymentMethodReceiverPay: true, PaymentMethodMonthly: true}
	if !validPayments[cmd.PaymentMethod] {
		return ErrOrderInvalidArgument
	}
	if cmd.ChannelCode == "LOCKER" && cmd.LockerID <= 0 {
		return ErrOrderInvalidArgument
	}
	if (cmd.ChannelCode == "SF_STATION" || cmd.ChannelCode == "PARTNER_STORE") && cmd.ServicePointID <= 0 {
		return ErrOrderInvalidArgument
	}
	if cmd.ChannelCode == "DOOR_PICKUP" && cmd.PickupDate == "" {
		return ErrOrderInvalidArgument
	}
	return nil
}

func generateOrderNo() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("SF%s%s", time.Now().Format("20060102150405"), hex.EncodeToString(b)[:6])
}
