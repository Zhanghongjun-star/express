package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	v3 "shunfeng-miniprogram/api/order/v3"

	"github.com/go-kratos/kratos/v3/errors"
)

var (
	// ErrExpressOrderNotFound 订单不存在时返回。
	ErrExpressOrderNotFound = errors.NotFound(v3.ErrorReason_EXPRESS_ORDER_NOT_FOUND.String(), "express order not found")
	// ErrExpressOrderInvalidArgument 订单参数无效时返回。
	ErrExpressOrderInvalidArgument = errors.BadRequest(v3.ErrorReason_EXPRESS_ORDER_INVALID_ARGUMENT.String(), "invalid express order argument")
)

// ExpressOrder 寄件订单领域对象。
type ExpressOrder struct {
	ID             int64
	OrderNo        string
	TrackingNo     string
	UserNo         string
	ExpressCompany string

	SenderName      string
	SenderPhone     string
	SenderProvince  string
	SenderCity      string
	SenderDistrict  string
	SenderAddress   string
	ReceiverName    string
	ReceiverPhone   string
	ReceiverProvince string
	ReceiverCity    string
	ReceiverDistrict string
	ReceiverAddress string

	ItemName      string
	ItemCategory  string
	ItemQuantity  int
	Weight        float64
	Length        float64
	Width         float64
	Height        float64
	DeclaredValue float64
	IsFragile     bool
	IsBattery     bool

	FreightFee   float64
	InsuranceFee float64
	PackageFee   float64
	TotalFee     float64
	PayMethod    string
	PayStatus    string

	OrderStatus  string
	PickupTime   *time.Time
	DeliveredTime *time.Time
	SignedBy     string
	Remark       string

	// Packages 选用的包装材料及数量列表。
	Packages []PackageItem
	// PackageDetail 费用计算过程明细（JSON 字符串，含打包费/保价费/合计/公式）。
	PackageDetail string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// PackageItem 订单选用的包装材料及数量（仅存材料 id，名称/单价由 packaging_material 表解析）。
type PackageItem struct {
	MaterialID int64 // 包装材料ID，对应 packaging_material 表的 ID
	Quantity   int   // 数量
}

// materialRate 包装材料单价（含名称），用于下单时按 id 查表并生成过程明细。
type materialRate struct {
	Name      string  `json:"name"`       // 材料名称（来自 packaging_material 表）
	UnitPrice float64 `json:"unit_price"` // 单价（元）
}

// packageLine 单条包装材料的计费明细。
type packageLine struct {
	Material  string  `json:"material"`   // 材料名称
	Quantity  int     `json:"quantity"`   // 数量
	UnitPrice float64 `json:"unit_price"` // 单价（元）
	Subtotal  float64 `json:"subtotal"`   // 小计（单价×数量）
}

// insuranceLine 保价费计算过程。
type insuranceLine struct {
	DeclaredValue float64 `json:"declared_value"` // 申报价值
	Fragile       bool    `json:"fragile"`        // 是否易碎
	Formula       string  `json:"formula"`        // 计费公式说明
	Fee           float64 `json:"fee"`            // 保价费（元）
}

// feeBreakdown 费用计算过程明细，序列化为 package_detail 存储/返回。
type feeBreakdown struct {
	Packages      []packageLine `json:"packages"`        // 各包装材料明细
	PackageFee    float64       `json:"package_fee"`     // 打包总费用
	Insurance     insuranceLine `json:"insurance"`       // 保价费明细
	TwoItemsTotal float64       `json:"two_items_total"` // 打包费 + 保价费 合计
	FreightFee    float64       `json:"freight_fee"`     // 基础运费
	GrandTotal    float64       `json:"grand_total"`     // 订单总费用（含运费）
}

// ExpressOrderRepo 寄件订单仓储接口（biz 声明，data 实现）。
type ExpressOrderRepo interface {
	Create(ctx context.Context, o *ExpressOrder) (*ExpressOrder, error)
	GetByNo(ctx context.Context, orderNo string) (*ExpressOrder, error)
	List(ctx context.Context, userNo, status string, opts ...ListOption) ([]*ExpressOrder, int32, error)
	Update(ctx context.Context, o *ExpressOrder, mask []string) (*ExpressOrder, error)
	Delete(ctx context.Context, orderNo string) error
}

// ExpressOrderUsecase 寄件订单用例。
type ExpressOrderUsecase struct {
	repo         ExpressOrderRepo
	materialRepo PackagingMaterialRepo // 包装材料仓储，用于下单时查询材料单价
}

// NewExpressOrderUsecase 创建寄件订单用例。
func NewExpressOrderUsecase(repo ExpressOrderRepo, materialRepo PackagingMaterialRepo) *ExpressOrderUsecase {
	return &ExpressOrderUsecase{repo: repo, materialRepo: materialRepo}
}

// CreateExpressOrder 创建寄件订单，自动生成订单号并填充默认状态。
func (uc *ExpressOrderUsecase) CreateExpressOrder(ctx context.Context, o *ExpressOrder) (*ExpressOrder, error) {
	if o == nil {
		return nil, ErrExpressOrderInvalidArgument
	}
	if o.OrderNo == "" {
		o.OrderNo = generateOrderNo()
	}
	if o.OrderStatus == "" {
		o.OrderStatus = "pending"
	}
	if o.PayStatus == "" {
		o.PayStatus = "unpaid"
	}
	// 服务端权威计算：先按 id 从 packaging_material 表查询所选材料单价，再算打包费、保价费，最后汇总并生成过程明细。
	rates := make(map[int64]*materialRate, len(o.Packages))
	for _, it := range o.Packages {
		if _, ok := rates[it.MaterialID]; ok {
			continue
		}
		m, err := uc.materialRepo.GetByID(ctx, it.MaterialID)
		if err != nil {
			return nil, errors.BadRequest(
				v3.ErrorReason_EXPRESS_ORDER_INVALID_ARGUMENT.String(),
				fmt.Sprintf("未知包装材料(id=%d)：%v", it.MaterialID, err),
			)
		}
		rates[it.MaterialID] = &materialRate{Name: m.Name, UnitPrice: m.UnitPrice}
	}
	pkgFee, pkgLines, err := calcPackageFee(o.Packages, rates)
	if err != nil {
		return nil, err
	}
	insFee, insFormula, err := calcInsuranceFee(o.DeclaredValue, o.IsFragile)
	if err != nil {
		return nil, err
	}
	o.PackageFee = pkgFee
	o.InsuranceFee = insFee
	// 订单总费用 = 基础运费 + 打包费 + 保价费。
	o.TotalFee = o.FreightFee + pkgFee + insFee
	// 组装费用计算过程明细（JSON），方便前端/调用方查看分项与公式。
	breakdown := feeBreakdown{
		Packages:      pkgLines,
		PackageFee:    pkgFee,
		Insurance:     insuranceLine{DeclaredValue: o.DeclaredValue, Fragile: o.IsFragile, Formula: insFormula, Fee: insFee},
		TwoItemsTotal: pkgFee + insFee, // 打包费 + 保价费 两项合计
		FreightFee:    o.FreightFee,
		GrandTotal:    o.TotalFee,
	}
	if detail, err := json.Marshal(breakdown); err == nil {
		o.PackageDetail = string(detail)
	}
	return uc.repo.Create(ctx, o)
}

// GetExpressOrder 根据订单号获取订单。
func (uc *ExpressOrderUsecase) GetExpressOrder(ctx context.Context, orderNo string) (*ExpressOrder, error) {
	if orderNo == "" {
		return nil, ErrExpressOrderInvalidArgument
	}
	return uc.repo.GetByNo(ctx, orderNo)
}

// ListExpressOrders 列出指定用户的寄件订单。
func (uc *ExpressOrderUsecase) ListExpressOrders(ctx context.Context, userNo, status string, opts ...ListOption) ([]*ExpressOrder, int32, error) {
	if userNo == "" {
		return nil, 0, ErrExpressOrderInvalidArgument
	}
	return uc.repo.List(ctx, userNo, status, opts...)
}

// UpdateExpressOrder 更新寄件订单的指定字段。
func (uc *ExpressOrderUsecase) UpdateExpressOrder(ctx context.Context, o *ExpressOrder, mask []string) (*ExpressOrder, error) {
	if o == nil || o.OrderNo == "" {
		return nil, ErrExpressOrderInvalidArgument
	}
	return uc.repo.Update(ctx, o, mask)
}

// DeleteExpressOrder 删除寄件订单。
func (uc *ExpressOrderUsecase) DeleteExpressOrder(ctx context.Context, orderNo string) error {
	if orderNo == "" {
		return ErrExpressOrderInvalidArgument
	}
	return uc.repo.Delete(ctx, orderNo)
}

// generateOrderNo 生成内部订单号，格式：SF + 时间戳(14位) + 4位随机。
func generateOrderNo() string {
	ts := time.Now().Format("20060102150405")
	return fmt.Sprintf("SF%s%04d", ts, rand.Intn(10000))
}

// maxInsuredValue 内地个人小件快递单票最高保价金额（元）。超出则无法投保。
const maxInsuredValue = 20000.0

// calcInsuranceFee 根据申报保价金额与是否易碎计算保价费（元）。
// 计费标准：
//   - 普通物品：≤500 元固定 1 元；500<金额≤1000 元固定 2 元；>1000 元按 5‰（0.5%）计费。
//   - 易碎物品：≤500 元固定 3 元；500<金额≤1000 元固定 7 元；>1000 元按 1.5% 计费。
//
// 计算后四舍五入取整，无最低门槛；申报金额超过单票上限时返回错误。
// 返回值：保价费、计费公式说明、错误。
func calcInsuranceFee(declaredValue float64, fragile bool) (float64, string, error) {
	if declaredValue <= 0 {
		return 0, "", nil
	}
	if declaredValue > maxInsuredValue {
		return 0, "", errors.BadRequest(
			v3.ErrorReason_EXPRESS_ORDER_INVALID_ARGUMENT.String(),
			fmt.Sprintf("保价金额 %.0f 元已超过单票最高保价 %.0f 元，无法投保保价服务", declaredValue, maxInsuredValue),
		)
	}
	var raw float64
	var formula string
	if fragile {
		// 易碎品专属费率
		switch {
		case declaredValue <= 500:
			raw = 3
			formula = "易碎品 ≤500元 固定 3 元"
		case declaredValue <= 1000:
			raw = 7
			formula = "易碎品 500<金额≤1000 固定 7 元"
		default:
			raw = declaredValue * 0.015
			formula = fmt.Sprintf("%.0f × 1.5%% = %.2f", declaredValue, raw)
		}
	} else {
		// 普通物品费率
		switch {
		case declaredValue <= 500:
			raw = 1
			formula = "普通品 ≤500元 固定 1 元"
		case declaredValue <= 1000:
			raw = 2
			formula = "普通品 500<金额≤1000 固定 2 元"
		default:
			raw = declaredValue * 0.005
			formula = fmt.Sprintf("%.0f × 0.5%% = %.2f", declaredValue, raw)
		}
	}
	return math.Round(raw), formula, nil
}

// calcPackageFee 根据选用的包装材料与数量、以及按 id 查表得到的单价计算打包总费用。
// 返回：打包总费用、逐条计费明细；材料不在单价映射中（id 不存在）时返回错误。
func calcPackageFee(items []PackageItem, rates map[int64]*materialRate) (float64, []packageLine, error) {
	var total float64
	lines := make([]packageLine, 0, len(items))
	for _, it := range items {
		// 在单价映射中按材料 id 查找单价（已通过 materialRepo 从 packaging_material 表解析）
		rate, ok := rates[it.MaterialID]
		if !ok {
			return 0, nil, errors.BadRequest(
				v3.ErrorReason_EXPRESS_ORDER_INVALID_ARGUMENT.String(),
				fmt.Sprintf("未知包装材料(id=%d)：请先从包装材料表中选择", it.MaterialID),
			)
		}
		// 数量缺省按 1 计
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		sub := rate.UnitPrice * float64(qty)
		lines = append(lines, packageLine{
			Material:  rate.Name,
			Quantity:  qty,
			UnitPrice: rate.UnitPrice,
			Subtotal:  sub,
		})
		total += sub
	}
	return total, lines, nil
}
