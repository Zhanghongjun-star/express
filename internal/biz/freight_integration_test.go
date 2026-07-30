package biz

import (
	"context"
	"testing"
)

// TestFreightEstimateIntegration 集成测试运费计算全流程（独立于 data 包）
func TestFreightEstimateIntegration(t *testing.T) {
	repo := &freightRepoWithPricing{
		pricing: getFullPricingTable(),
	}
	uc := NewFreightUsecase(repo, nil)

	tests := []struct {
		name                                   string
		sendProv, sendCity, recvProv, recvCity string
		weight                                 float64
		length, width, height                  int32
		expressType                            ExpressType
		insureValue                            float64
		wantMinBase                            float64
		wantMaxBase                            float64
		wantCalcWeight                         float64
		wantInsureMin                          float64
	}{
		{name: "同城-深圳到深圳-标快-2.5kg", sendProv: "广东省", sendCity: "深圳市", recvProv: "广东省", recvCity: "深圳市", weight: 2.5, length: 30, width: 20, height: 15, expressType: ExpressTypeStandard, insureValue: 0, wantMinBase: 2, wantMaxBase: 10, wantCalcWeight: 2.5, wantInsureMin: 0},
		{name: "跨省-深圳到北京-标快-2.5kg", sendProv: "广东省", sendCity: "深圳市", recvProv: "北京市", recvCity: "北京市", weight: 2.5, length: 30, width: 20, height: 15, expressType: ExpressTypeStandard, insureValue: 0, wantMinBase: 10, wantMaxBase: 30, wantCalcWeight: 2.5, wantInsureMin: 0},
		{name: "体积重大于实际重", sendProv: "广东省", sendCity: "深圳市", recvProv: "北京市", recvCity: "北京市", weight: 2, length: 60, width: 50, height: 40, expressType: ExpressTypeStandard, insureValue: 0, wantMinBase: 50, wantMaxBase: 120, wantCalcWeight: 20, wantInsureMin: 0},
		{name: "特惠件-深圳到北京-2.5kg", sendProv: "广东省", sendCity: "深圳市", recvProv: "北京市", recvCity: "北京市", weight: 2.5, length: 30, width: 20, height: 15, expressType: ExpressTypeEconomy, insureValue: 0, wantMinBase: 5, wantMaxBase: 20, wantCalcWeight: 2.5, wantInsureMin: 0},
		{name: "偏远-深圳到新疆-标快-2.5kg", sendProv: "广东省", sendCity: "深圳市", recvProv: "新疆维吾尔自治区", recvCity: "乌鲁木齐市", weight: 2.5, length: 30, width: 20, height: 15, expressType: ExpressTypeStandard, insureValue: 0, wantMinBase: 20, wantMaxBase: 50, wantCalcWeight: 2.5, wantInsureMin: 0},
		{name: "带保价-深圳到北京-2000元", sendProv: "广东省", sendCity: "深圳市", recvProv: "北京市", recvCity: "北京市", weight: 1, length: 20, width: 15, height: 10, expressType: ExpressTypeStandard, insureValue: 2000, wantMinBase: 10, wantMaxBase: 30, wantCalcWeight: 1, wantInsureMin: 1},
		{name: "同省不同市-广州到深圳-2.5kg", sendProv: "广东省", sendCity: "广州市", recvProv: "广东省", recvCity: "深圳市", weight: 2.5, length: 30, width: 20, height: 15, expressType: ExpressTypeStandard, insureValue: 0, wantMinBase: 3, wantMaxBase: 15, wantCalcWeight: 2.5, wantInsureMin: 0},
		{name: "同经济圈-广州到南宁-2.5kg", sendProv: "广东省", sendCity: "广州市", recvProv: "广西壮族自治区", recvCity: "南宁市", weight: 2.5, length: 30, width: 20, height: 15, expressType: ExpressTypeStandard, insureValue: 0, wantMinBase: 5, wantMaxBase: 20, wantCalcWeight: 2.5, wantInsureMin: 0},
		{name: "港澳台-深圳到香港-特惠", sendProv: "广东省", sendCity: "深圳市", recvProv: "香港特别行政区", recvCity: "香港", weight: 1, length: 20, width: 15, height: 10, expressType: ExpressTypeEconomy, insureValue: 0, wantMinBase: 10, wantMaxBase: 40, wantCalcWeight: 1, wantInsureMin: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &FreightRequest{
				SenderProvince:   tt.sendProv,
				SenderCity:       tt.sendCity,
				ReceiverProvince: tt.recvProv,
				ReceiverCity:     tt.recvCity,
				Weight:           tt.weight,
				Length:           tt.length,
				Width:            tt.width,
				Height:           tt.height,
				ExpressType:      tt.expressType,
				InsureValue:      tt.insureValue,
			}

			result, err := uc.Estimate(context.Background(), req)
			if err != nil {
				t.Fatalf("Estimate() error = %v", err)
			}

			if result.CalcWeight != tt.wantCalcWeight {
				t.Errorf("CalcWeight = %v, want %v", result.CalcWeight, tt.wantCalcWeight)
			}
			if result.BaseFreight < tt.wantMinBase {
				t.Errorf("BaseFreight = %v, want >= %v (too low)", result.BaseFreight, tt.wantMinBase)
			}
			if result.BaseFreight > tt.wantMaxBase {
				t.Errorf("BaseFreight = %v, want <= %v (too high)", result.BaseFreight, tt.wantMaxBase)
			}
			if result.InsureFee < tt.wantInsureMin {
				t.Errorf("InsureFee = %v, want >= %v", result.InsureFee, tt.wantInsureMin)
			}

			t.Logf("OK | CalcWeight=%.2f BaseFreight=%.2f InsureFee=%.2f Total=%.2f Tips=%s",
				result.CalcWeight, result.BaseFreight, result.InsureFee, result.TotalPrice, result.Tips)
		})
	}
}

// freightRepoWithPricing 测试用仓库实现，持有完整定价表。
type freightRepoWithPricing struct {
	pricing map[ExpressType]map[Zone][]PricingTier
}

func (r *freightRepoWithPricing) GetPricingTiers(_ context.Context, expressType ExpressType, zone Zone) ([]PricingTier, error) {
	if tiers, ok := r.pricing[expressType][zone]; ok {
		return tiers, nil
	}
	return r.pricing[expressType][ZoneDistant], nil
}

// getFullPricingTable 返回与 data 包一致的完整阶梯定价表副本。
func getFullPricingTable() map[ExpressType]map[Zone][]PricingTier {
	return map[ExpressType]map[Zone][]PricingTier{
		ExpressTypeStandard: {
			ZoneSameCity:     {{WeightMax: 1, UnitPrice: 12}, {WeightMax: 3, UnitPrice: 2}, {WeightMax: 5, UnitPrice: 1.5}, {WeightMax: 10, UnitPrice: 1}, {WeightMax: 0, UnitPrice: 0.8}},
			ZoneSameProvince: {{WeightMax: 1, UnitPrice: 15}, {WeightMax: 3, UnitPrice: 3}, {WeightMax: 5, UnitPrice: 2.5}, {WeightMax: 10, UnitPrice: 2}, {WeightMax: 0, UnitPrice: 1.5}},
			ZoneSameCircle:   {{WeightMax: 1, UnitPrice: 18}, {WeightMax: 3, UnitPrice: 5}, {WeightMax: 5, UnitPrice: 3.5}, {WeightMax: 10, UnitPrice: 2.5}, {WeightMax: 0, UnitPrice: 2}},
			ZoneNearby:       {{WeightMax: 1, UnitPrice: 23}, {WeightMax: 3, UnitPrice: 7}, {WeightMax: 5, UnitPrice: 5}, {WeightMax: 10, UnitPrice: 4}, {WeightMax: 0, UnitPrice: 3}},
			ZoneDistant:      {{WeightMax: 1, UnitPrice: 25}, {WeightMax: 3, UnitPrice: 10}, {WeightMax: 5, UnitPrice: 7}, {WeightMax: 10, UnitPrice: 5}, {WeightMax: 0, UnitPrice: 4}},
			ZoneRemote:       {{WeightMax: 1, UnitPrice: 30}, {WeightMax: 3, UnitPrice: 15}, {WeightMax: 5, UnitPrice: 12}, {WeightMax: 10, UnitPrice: 10}, {WeightMax: 0, UnitPrice: 8}},
			ZoneSpecial:      {{WeightMax: 1, UnitPrice: 36}, {WeightMax: 3, UnitPrice: 20}, {WeightMax: 5, UnitPrice: 16}, {WeightMax: 10, UnitPrice: 14}, {WeightMax: 0, UnitPrice: 12}},
		},
		ExpressTypeSpecial: {
			ZoneSameCity:     {{WeightMax: 1, UnitPrice: 18}, {WeightMax: 3, UnitPrice: 3}, {WeightMax: 5, UnitPrice: 2.25}, {WeightMax: 10, UnitPrice: 1.5}, {WeightMax: 0, UnitPrice: 1.2}},
			ZoneSameProvince: {{WeightMax: 1, UnitPrice: 22.5}, {WeightMax: 3, UnitPrice: 4.5}, {WeightMax: 5, UnitPrice: 3.75}, {WeightMax: 10, UnitPrice: 3}, {WeightMax: 0, UnitPrice: 2.25}},
			ZoneSameCircle:   {{WeightMax: 1, UnitPrice: 27}, {WeightMax: 3, UnitPrice: 7.5}, {WeightMax: 5, UnitPrice: 5.25}, {WeightMax: 10, UnitPrice: 3.75}, {WeightMax: 0, UnitPrice: 3}},
			ZoneNearby:       {{WeightMax: 1, UnitPrice: 34.5}, {WeightMax: 3, UnitPrice: 10.5}, {WeightMax: 5, UnitPrice: 7.5}, {WeightMax: 10, UnitPrice: 6}, {WeightMax: 0, UnitPrice: 4.5}},
			ZoneDistant:      {{WeightMax: 1, UnitPrice: 37.5}, {WeightMax: 3, UnitPrice: 15}, {WeightMax: 5, UnitPrice: 10.5}, {WeightMax: 10, UnitPrice: 7.5}, {WeightMax: 0, UnitPrice: 6}},
			ZoneRemote:       {{WeightMax: 1, UnitPrice: 45}, {WeightMax: 3, UnitPrice: 22.5}, {WeightMax: 5, UnitPrice: 18}, {WeightMax: 10, UnitPrice: 15}, {WeightMax: 0, UnitPrice: 12}},
			ZoneSpecial:      {{WeightMax: 1, UnitPrice: 54}, {WeightMax: 3, UnitPrice: 30}, {WeightMax: 5, UnitPrice: 24}, {WeightMax: 10, UnitPrice: 21}, {WeightMax: 0, UnitPrice: 18}},
		},
		ExpressTypeEconomy: {
			ZoneSameCity:     {{WeightMax: 1, UnitPrice: 9.6}, {WeightMax: 3, UnitPrice: 1.6}, {WeightMax: 5, UnitPrice: 1.2}, {WeightMax: 10, UnitPrice: 0.8}, {WeightMax: 0, UnitPrice: 0.64}},
			ZoneSameProvince: {{WeightMax: 1, UnitPrice: 12}, {WeightMax: 3, UnitPrice: 2.4}, {WeightMax: 5, UnitPrice: 2}, {WeightMax: 10, UnitPrice: 1.6}, {WeightMax: 0, UnitPrice: 1.2}},
			ZoneSameCircle:   {{WeightMax: 1, UnitPrice: 14.4}, {WeightMax: 3, UnitPrice: 4}, {WeightMax: 5, UnitPrice: 2.8}, {WeightMax: 10, UnitPrice: 2}, {WeightMax: 0, UnitPrice: 1.6}},
			ZoneNearby:       {{WeightMax: 1, UnitPrice: 18.4}, {WeightMax: 3, UnitPrice: 5.6}, {WeightMax: 5, UnitPrice: 4}, {WeightMax: 10, UnitPrice: 3.2}, {WeightMax: 0, UnitPrice: 2.4}},
			ZoneDistant:      {{WeightMax: 1, UnitPrice: 20}, {WeightMax: 3, UnitPrice: 8}, {WeightMax: 5, UnitPrice: 5.6}, {WeightMax: 10, UnitPrice: 4}, {WeightMax: 0, UnitPrice: 3.2}},
			ZoneRemote:       {{WeightMax: 1, UnitPrice: 24}, {WeightMax: 3, UnitPrice: 12}, {WeightMax: 5, UnitPrice: 9.6}, {WeightMax: 10, UnitPrice: 8}, {WeightMax: 0, UnitPrice: 6.4}},
			ZoneSpecial:      {{WeightMax: 1, UnitPrice: 28.8}, {WeightMax: 3, UnitPrice: 16}, {WeightMax: 5, UnitPrice: 12.8}, {WeightMax: 10, UnitPrice: 11.2}, {WeightMax: 0, UnitPrice: 9.6}},
		},
	}
}
