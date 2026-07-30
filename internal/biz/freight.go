package biz

import (
	"context"
	"maps"
	"math"
	"slices"

	v2 "shunfeng-miniprogram/api/order/v2"

	"github.com/go-kratos/kratos/v3/errors"
)

// =============================================================================
// 错误定义
// =============================================================================

var (
	// ErrFreightCalcFailed 运费计算失败。
	ErrFreightCalcFailed = errors.InternalServer(v2.ErrorReason_ORDER_FREIGHT_CALC_FAILED.String(), "freight calculation failed")
	// ErrFreightInvalidArgument 运费预估请求参数无效。
	ErrFreightInvalidArgument = errors.BadRequest(v2.ErrorReason_ORDER_FREIGHT_INVALID_ARGUMENT.String(), "invalid freight argument")
)

// =============================================================================
// 计算常量
// =============================================================================

const (
	// VolDivisor 标准体积重量除数（cm³/kg），长×宽×高÷6000。
	VolDivisor = 6000
	// VolDivisorEconomy 特惠件体积重量除数（cm³/kg），长×宽×高÷7000，更宽松。
	VolDivisorEconomy = 7000
	// MinWeight 最小计费重量（kg），不足 1kg 按 1kg 计。
	MinWeight = 1.0
	// InsureRate 保价费率（千分之五）。
	InsureRate = 0.005
)

// =============================================================================
// 快递类型
// =============================================================================

// ExpressType 快递类型枚举。
type ExpressType int

const (
	// ExpressTypeStandard 标快（标准快递）。
	ExpressTypeStandard ExpressType = 1
	// ExpressTypeSpecial 特快（时效优先）。
	ExpressTypeSpecial ExpressType = 2
	// ExpressTypeEconomy 特惠（经济快递）。
	ExpressTypeEconomy ExpressType = 3
)

// =============================================================================
// 区域等级
// =============================================================================

// Zone 收寄件双方距离区域等级（0-6），数字越大距离越远。
type Zone int

const (
	// ZoneSameCity 同城。
	ZoneSameCity Zone = 0
	// ZoneSameProvince 同省不同市。
	ZoneSameProvince Zone = 1
	// ZoneSameCircle 同经济圈（长三角/珠三角/京津冀/成渝）。
	ZoneSameCircle Zone = 2
	// ZoneNearby 邻近省份。
	ZoneNearby Zone = 3
	// ZoneDistant 远距离省份。
	ZoneDistant Zone = 4
	// ZoneRemote 偏远地区（新疆、西藏、青海、内蒙古、宁夏、甘肃）。
	ZoneRemote Zone = 5
	// ZoneSpecial 港澳台。
	ZoneSpecial Zone = 6
)

// =============================================================================
// 领域对象（DO）
// =============================================================================

// FreightRequest 运费预估请求（领域对象）。
type FreightRequest struct {
	SenderProvince   string      // 寄件省份
	SenderCity       string      // 寄件城市
	SenderArea       string      // 寄件区/县（可选）
	ReceiverProvince string      // 收件省份
	ReceiverCity     string      // 收件城市
	ReceiverArea     string      // 收件区/县（可选）
	Weight           float64     // 实际重量（kg）
	Length           int32       // 长（cm）
	Width            int32       // 宽（cm）
	Height           int32       // 高（cm）
	ExpressType      ExpressType // 快递类型
	InsureValue      float64     // 保价金额（0 表示不保价）
}

// FreightResult 运费预估结果（领域对象）。
type FreightResult struct {
	BaseFreight float64 // 基础运费（元）
	InsureFee   float64 // 保价费（元）
	TotalPrice  float64 // 总价（元）
	CalcWeight  float64 // 计费重量（kg）
	Tips        string  // 提示信息
}

// PricingTier 阶梯定价档位。
// 包含该重量区间的单价，WeightMax=0 表示无上限。
type PricingTier struct {
	WeightMax float64 // 重量区间上限（kg），0 表示无上限
	UnitPrice float64 // 单价（元/kg）
}

// =============================================================================
// 仓储接口（Repo）
// =============================================================================

// FreightRepo 运费预估仓储接口，负责提供定价数据。
type FreightRepo interface {
	// GetPricingTiers 根据快递类型和区域等级查询阶梯定价表。
	GetPricingTiers(ctx context.Context, expressType ExpressType, zone Zone) ([]PricingTier, error)
}

// =============================================================================
// 用例（Usecase）
// =============================================================================

// FreightUsecase 运费预估用例，包含核心计价算法。
type FreightUsecase struct {
	repo FreightRepo
}

// NewFreightUsecase 创建运费预估用例。
func NewFreightUsecase(repo FreightRepo) *FreightUsecase {
	return &FreightUsecase{repo: repo}
}

// Estimate 执行运费预估。
// 步骤：
//   1. 参数校验
//   2. 确定区域等级
//   3. 计算体积重量与计费重量
//   4. 查询阶梯定价并计算基础运费
//   5. 计算保价费
//   6. 汇总总价
func (uc *FreightUsecase) Estimate(ctx context.Context, req *FreightRequest) (*FreightResult, error) {
	// 1. 参数校验
	if err := uc.validateRequest(req); err != nil {
		return nil, err
	}

	// 2. 确定收发件双方区域等级
	zone := uc.determineZone(req.SenderProvince, req.SenderCity, req.ReceiverProvince, req.ReceiverCity)

	// 3. 计算体积重量与计费重量
	calcWeight := uc.calcVolumetricWeight(req.Weight, float64(req.Length), float64(req.Width), float64(req.Height), req.ExpressType)

	// 4. 查询阶梯定价表并计算基础运费
	tiers, err := uc.repo.GetPricingTiers(ctx, req.ExpressType, zone)
	if err != nil {
		return nil, ErrFreightCalcFailed
	}
	baseFreight := uc.calcBaseFreight(calcWeight, tiers)

	// 5. 计算保价费
	var insureFee float64
	if req.InsureValue > 0 {
		insureFee = math.Max(req.InsureValue*InsureRate, 1.0) // 保价费最低 1 元
	}

	// 6. 汇总总价
	totalPrice := math.Round((baseFreight+insureFee)*100) / 100

	// 兜底最低收费
	if totalPrice < 1.0 {
		totalPrice = 1.0
	}

	return &FreightResult{
		BaseFreight: math.Round(baseFreight*100) / 100,
		InsureFee:   math.Round(insureFee*100) / 100,
		TotalPrice:  totalPrice,
		CalcWeight:  math.Round(calcWeight*100) / 100,
		Tips:        uc.buildTips(zone, req.ExpressType),
	}, nil
}

// =============================================================================
// 私有方法：区域判定
// =============================================================================

// 经济圈定义：同一经济圈内的省份视为较近（Zone 2）。
var economicCircles = map[string]int{
	// 长三角：上海、江苏、浙江、安徽
	"上海": 1, "上海市": 1,
	"江苏": 1, "江苏省": 1,
	"浙江": 1, "浙江省": 1,
	"安徽": 1, "安徽省": 1,
	// 珠三角：广东、广西、海南
	"广东": 2, "广东省": 2,
	"广西": 2, "广西壮族自治区": 2,
	"海南": 2, "海南省": 2,
	// 京津冀：北京、天津、河北
	"北京": 3, "北京市": 3,
	"天津": 3, "天津市": 3,
	"河北": 3, "河北省": 3,
	// 成渝：四川、重庆
	"四川": 4, "四川省": 4,
	"重庆": 4, "重庆市": 4,
}

// provinceZones 省份距离分组，用于不属于同一经济圈的省份之间判断。
// Zone 3：较近省份；Zone 4：较远省份；Zone 5：偏远省份。
var provinceZones = map[int][]string{
	3: {"福建", "福建省", "山东", "山东省", "湖南", "湖南省", "湖北", "湖北省",
		"河南", "河南省", "江西", "江西省", "陕西", "陕西省", "山西", "山西省"},
	4: {"云南", "云南省", "贵州", "贵州省", "黑龙江", "黑龙江省",
		"吉林", "吉林省", "辽宁", "辽宁省"},
	5: {"新疆", "新疆维吾尔自治区", "西藏", "西藏自治区", "青海", "青海省",
		"内蒙古", "内蒙古自治区", "宁夏", "宁夏回族自治区", "甘肃", "甘肃省"},
}

// specialRegions 港澳台地区（Zone 6）。
var specialRegions = map[string]bool{
	"香港": true, "香港特别行政区": true,
	"澳门": true, "澳门特别行政区": true,
	"台湾": true, "台湾省": true,
}

// determineZone 根据收发件省份和城市确定区域等级。
//
// 判定逻辑（优先级从高到低）：
//   1. 任意一方是港澳台 → Zone 6
//   2. 任意一方是偏远省份 → Zone 5
//   3. 同城 → Zone 0
//   4. 同省 → Zone 1
//   5. 同经济圈 → Zone 2
//   6. 较近省份 → Zone 3
//   7. 较远省份 → Zone 4
//   8. 默认 → Zone 4
func (uc *FreightUsecase) determineZone(sendProvince, sendCity, recvProvince, recvCity string) Zone {
	// 港澳台判定
	if specialRegions[sendProvince] || specialRegions[recvProvince] {
		return ZoneSpecial
	}

	// 偏远省份判定
	sendRemote := uc.isRemoteProvince(sendProvince)
	recvRemote := uc.isRemoteProvince(recvProvince)
	if sendRemote || recvRemote {
		return ZoneRemote
	}

	// 同城判定
	if uc.isSameCity(sendCity, recvCity) {
		return ZoneSameCity
	}

	// 同省判定
	if uc.normalizeProvince(sendProvince) == uc.normalizeProvince(recvProvince) {
		return ZoneSameProvince
	}

	// 同经济圈判定
	if uc.isSameEconomicCircle(sendProvince, recvProvince) {
		return ZoneSameCircle
	}

	// 省份距离判定
	return uc.determineProvinceZone(sendProvince, recvProvince)
}

// normalizeProvince 标准化省份名称，去除"省"、"市"、"自治区"等后缀。
// 使用标签 switch 将同一省份的多种写法合并为一个标准简称。
func (uc *FreightUsecase) normalizeProvince(name string) string {
	switch name {
	case "上海市", "上海":
		return "上海"
	case "北京市", "北京":
		return "北京"
	case "天津市", "天津":
		return "天津"
	case "重庆市", "重庆":
		return "重庆"
	case "广东省", "广东":
		return "广东"
	case "江苏省", "江苏":
		return "江苏"
	case "浙江省", "浙江":
		return "浙江"
	case "安徽省", "安徽":
		return "安徽"
	case "广西壮族自治区", "广西":
		return "广西"
	case "海南省", "海南":
		return "海南"
	case "河北省", "河北":
		return "河北"
	case "四川省", "四川":
		return "四川"
	case "福建省", "福建":
		return "福建"
	case "山东省", "山东":
		return "山东"
	case "湖南省", "湖南":
		return "湖南"
	case "湖北省", "湖北":
		return "湖北"
	case "河南省", "河南":
		return "河南"
	case "江西省", "江西":
		return "江西"
	case "陕西省", "陕西":
		return "陕西"
	case "山西省", "山西":
		return "山西"
	case "云南省", "云南":
		return "云南"
	case "贵州省", "贵州":
		return "贵州"
	case "黑龙江省", "黑龙江":
		return "黑龙江"
	case "吉林省", "吉林":
		return "吉林"
	case "辽宁省", "辽宁":
		return "辽宁"
	case "新疆维吾尔自治区", "新疆":
		return "新疆"
	case "西藏自治区", "西藏":
		return "西藏"
	case "青海省", "青海":
		return "青海"
	case "内蒙古自治区", "内蒙古":
		return "内蒙古"
	case "宁夏回族自治区", "宁夏":
		return "宁夏"
	case "甘肃省", "甘肃":
		return "甘肃"
	case "香港特别行政区", "香港":
		return "香港"
	case "澳门特别行政区", "澳门":
		return "澳门"
	case "台湾省", "台湾":
		return "台湾"
	default:
		return name
	}
}

// isSameCity 判断是否为同一城市。
func (uc *FreightUsecase) isSameCity(city1, city2 string) bool {
	// 去除"市"后缀后比较
	n1 := city1
	n2 := city2
	if len(n1) > 0 && n1[len(n1)-3:] == "市" {
		n1 = n1[:len(n1)-3]
	}
	if len(n2) > 0 && n2[len(n2)-3:] == "市" {
		n2 = n2[:len(n2)-3]
	}
	return n1 == n2
}

// isRemoteProvince 判断是否为偏远省份（Zone 5）。
func (uc *FreightUsecase) isRemoteProvince(province string) bool {
	if slices.Contains(provinceZones[5], province) {
		return true
	}
	// 同时检查标准化名称
	normalized := uc.normalizeProvince(province)
	for _, p := range provinceZones[5] {
		if normalized == uc.normalizeProvince(p) {
			return true
		}
	}
	return false
}

// isSameEconomicCircle 判断两个省份是否属于同一经济圈。
func (uc *FreightUsecase) isSameEconomicCircle(prov1, prov2 string) bool {
	p1 := uc.normalizeProvince(prov1)
	p2 := uc.normalizeProvince(prov2)

	// 查找省份对应的经济圈编号
	circle1, ok1 := economyCircleMap[p1]
	circle2, ok2 := economyCircleMap[p2]

	if ok1 && ok2 && circle1 == circle2 {
		return true
	}
	return false
}

// economyCircleMap 标准化省份 → 经济圈编号 映射表。
var economyCircleMap = buildCircleMap()

// buildCircleMap 构建标准化省份名称到经济圈编号的映射。
func buildCircleMap() map[string]int {
	result := make(map[string]int, len(economicCircles))
	maps.Copy(result, economicCircles)
	return result
}

// determineProvinceZone 非经济圈的省份之间，根据距离分组表判定 Zone。
func (uc *FreightUsecase) determineProvinceZone(p1, p2 string) Zone {
	zone1 := uc.getProvinceZone(p1)
	zone2 := uc.getProvinceZone(p2)
	// 取较远的等级
	if zone1 > zone2 {
		return Zone(zone1)
	}
	return Zone(zone2)
}

// getProvinceZone 获取省份的距离分组（3/4），非分组省份返回 4（默认较远）。
func (uc *FreightUsecase) getProvinceZone(province string) int {
	normalized := uc.normalizeProvince(province)
	// 检查省份是否落入近/远分组表
	for zone, provinces := range provinceZones {
		for _, p := range provinces {
			if uc.normalizeProvince(p) == normalized {
				return zone
			}
		}
	}
	// 经济圈省份（上海/广东/北京/四川等）未落入上述近/远分组。
	// 跨经济圈通信默认按远距离处理（Zone 4），避免把相距很远的两地误判为邻近。
	return 4 // 默认较远
}

// =============================================================================
// 私有方法：重量与运费计算
// =============================================================================

// calcVolumetricWeight 计算计费重量。
//
// 计费重量 = Max(实际重量, 体积重量, 最小重量)
//   体积重量 = 长 × 宽 × 高 ÷ 除数
//   标快/特快除数 = 6000，特惠除数 = 7000
//   不足 1kg 按 1kg 计
func (uc *FreightUsecase) calcVolumetricWeight(actualWeight float64, length, width, height float64, expressType ExpressType) float64 {
	divisor := float64(VolDivisor)
	if expressType == ExpressTypeEconomy {
		divisor = float64(VolDivisorEconomy)
	}

	// 仅在长宽高均有效时计算体积重量
	var volWeight float64
	if length > 0 && width > 0 && height > 0 {
		volWeight = (length * width * height) / divisor
	}

	// 取实际重量与体积重量中的较大值
	calcWeight := math.Max(actualWeight, volWeight)
	// 不足最小重量则取最小
	calcWeight = math.Max(calcWeight, MinWeight)

	return calcWeight
}

// calcBaseFreight 根据计费重量和阶梯定价表计算基础运费。
//
// 阶梯定价逻辑：
//   遍历定价档位，找到计费重量所在的区间，用该区间的单价 × 计费重量。
//   WeightMax=0 表示最后一档（无上限）。
func (uc *FreightUsecase) calcBaseFreight(calcWeight float64, tiers []PricingTier) float64 {
	if len(tiers) == 0 {
		return 0
	}
	for _, tier := range tiers {
		if calcWeight <= tier.WeightMax || tier.WeightMax == 0 {
			return calcWeight * tier.UnitPrice
		}
	}
	// 兜底：使用最后一档的单价
	last := tiers[len(tiers)-1]
	return calcWeight * last.UnitPrice
}

// validateRequest 校验运费预估请求参数。
func (uc *FreightUsecase) validateRequest(req *FreightRequest) error {
	if req.SenderProvince == "" || req.ReceiverProvince == "" {
		return ErrFreightInvalidArgument
	}
	if req.SenderCity == "" || req.ReceiverCity == "" {
		return ErrFreightInvalidArgument
	}
	if req.Weight <= 0 {
		return ErrFreightInvalidArgument
	}
	if req.Length < 0 || req.Width < 0 || req.Height < 0 {
		return ErrFreightInvalidArgument
	}
	if req.ExpressType < ExpressTypeStandard || req.ExpressType > ExpressTypeEconomy {
		return ErrFreightInvalidArgument
	}
	if req.InsureValue < 0 {
		return ErrFreightInvalidArgument
	}
	return nil
}

// buildTips 生成提示信息。
func (uc *FreightUsecase) buildTips(zone Zone, expressType ExpressType) string {
	typeNames := map[ExpressType]string{
		ExpressTypeStandard: "标快",
		ExpressTypeSpecial:  "特快",
		ExpressTypeEconomy:  "特惠",
	}
	typeName, ok := typeNames[expressType]
	if !ok {
		typeName = "未知"
	}

	zoneNames := map[Zone]string{
		ZoneSameCity:     "同城寄递",
		ZoneSameProvince: "同省寄递",
		ZoneSameCircle:   "同经济圈寄递",
		ZoneNearby:       "邻近省份寄递",
		ZoneDistant:      "远距离寄递",
		ZoneRemote:       "偏远地区寄递",
		ZoneSpecial:      "港澳台寄递",
	}
	zoneName, ok := zoneNames[zone]
	if !ok {
		zoneName = "普通寄递"
	}

	return "本次为" + zoneName + "，" + typeName + "计价，最终价格以收件员确认为准"
}
