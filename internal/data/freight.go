package data

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"shunfeng-miniprogram/internal/biz"
	"shunfeng-miniprogram/internal/pkg/amap"
)

// =============================================================================
// FreightRepo 内存实现
// =============================================================================

// freightRepo 运费预估仓储的内存实现，内嵌一份阶梯定价表。
type freightRepo struct {
	mu sync.RWMutex
}

// NewFreightRepo 创建运费预估仓储，返回接口类型。
func NewFreightRepo() biz.FreightRepo {
	return &freightRepo{}
}

// GetPricingTiers 根据快递类型和区域等级返回阶梯定价档位。
func (r *freightRepo) GetPricingTiers(_ context.Context, expressType biz.ExpressType, zone biz.Zone) ([]biz.PricingTier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tiers, ok := pricingTable[expressType][zone]
	if !ok {
		// 未匹配到的 Zone 使用默认（同类型 Zone 4）价格
		tiers = pricingTable[expressType][biz.ZoneDistant]
	}
	return tiers, nil
}

// =============================================================================
// 阶梯定价表（内存数据）
//
// 定价逻辑：
//   每行 (WeightMax, UnitPrice) 表示：
//     重量 ≤ WeightMax → 基础运费 = 计费重量 × UnitPrice
//   WeightMax=0 表示该档无上限。
//
// 参考顺丰官网公开计价规则建模，实际生产应替换为实时定价接口。
// =============================================================================

// pricingTable 按 [快递类型][区域等级] 组织的阶梯定价表。
//
// 快递类型：
//   1 = 标快（ExpressTypeStandard）
//   2 = 特快（ExpressTypeSpecial）
//   3 = 特惠（ExpressTypeEconomy）
//
// 区域等级：
//   0 = 同城  1 = 同省  2 = 同经济圈
//   3 = 邻近省  4 = 远距离省  5 = 偏远  6 = 港澳台
var pricingTable = map[biz.ExpressType]map[biz.Zone][]biz.PricingTier{
	// =========================================================================
	// 标快（标准快递）
	// =========================================================================
	biz.ExpressTypeStandard: {
		// 同城
		biz.ZoneSameCity: {
			{WeightMax: 1, UnitPrice: 12},
			{WeightMax: 3, UnitPrice: 2},
			{WeightMax: 5, UnitPrice: 1.5},
			{WeightMax: 10, UnitPrice: 1},
			{WeightMax: 0, UnitPrice: 0.8},
		},
		// 同省不同市
		biz.ZoneSameProvince: {
			{WeightMax: 1, UnitPrice: 15},
			{WeightMax: 3, UnitPrice: 3},
			{WeightMax: 5, UnitPrice: 2.5},
			{WeightMax: 10, UnitPrice: 2},
			{WeightMax: 0, UnitPrice: 1.5},
		},
		// 同经济圈
		biz.ZoneSameCircle: {
			{WeightMax: 1, UnitPrice: 18},
			{WeightMax: 3, UnitPrice: 5},
			{WeightMax: 5, UnitPrice: 3.5},
			{WeightMax: 10, UnitPrice: 2.5},
			{WeightMax: 0, UnitPrice: 2},
		},
		// 邻近省份
		biz.ZoneNearby: {
			{WeightMax: 1, UnitPrice: 23},
			{WeightMax: 3, UnitPrice: 7},
			{WeightMax: 5, UnitPrice: 5},
			{WeightMax: 10, UnitPrice: 4},
			{WeightMax: 0, UnitPrice: 3},
		},
		// 远距离省份
		biz.ZoneDistant: {
			{WeightMax: 1, UnitPrice: 25},
			{WeightMax: 3, UnitPrice: 10},
			{WeightMax: 5, UnitPrice: 7},
			{WeightMax: 10, UnitPrice: 5},
			{WeightMax: 0, UnitPrice: 4},
		},
		// 偏远地区
		biz.ZoneRemote: {
			{WeightMax: 1, UnitPrice: 30},
			{WeightMax: 3, UnitPrice: 15},
			{WeightMax: 5, UnitPrice: 12},
			{WeightMax: 10, UnitPrice: 10},
			{WeightMax: 0, UnitPrice: 8},
		},
		// 港澳台
		biz.ZoneSpecial: {
			{WeightMax: 1, UnitPrice: 36},
			{WeightMax: 3, UnitPrice: 20},
			{WeightMax: 5, UnitPrice: 16},
			{WeightMax: 10, UnitPrice: 14},
			{WeightMax: 0, UnitPrice: 12},
		},
	},
	// =========================================================================
	// 特快（时效优先，价格为标快的 1.5 倍）
	// =========================================================================
	biz.ExpressTypeSpecial: {
		biz.ZoneSameCity: {
			{WeightMax: 1, UnitPrice: 18},
			{WeightMax: 3, UnitPrice: 3},
			{WeightMax: 5, UnitPrice: 2.25},
			{WeightMax: 10, UnitPrice: 1.5},
			{WeightMax: 0, UnitPrice: 1.2},
		},
		biz.ZoneSameProvince: {
			{WeightMax: 1, UnitPrice: 22.5},
			{WeightMax: 3, UnitPrice: 4.5},
			{WeightMax: 5, UnitPrice: 3.75},
			{WeightMax: 10, UnitPrice: 3},
			{WeightMax: 0, UnitPrice: 2.25},
		},
		biz.ZoneSameCircle: {
			{WeightMax: 1, UnitPrice: 27},
			{WeightMax: 3, UnitPrice: 7.5},
			{WeightMax: 5, UnitPrice: 5.25},
			{WeightMax: 10, UnitPrice: 3.75},
			{WeightMax: 0, UnitPrice: 3},
		},
		biz.ZoneNearby: {
			{WeightMax: 1, UnitPrice: 34.5},
			{WeightMax: 3, UnitPrice: 10.5},
			{WeightMax: 5, UnitPrice: 7.5},
			{WeightMax: 10, UnitPrice: 6},
			{WeightMax: 0, UnitPrice: 4.5},
		},
		biz.ZoneDistant: {
			{WeightMax: 1, UnitPrice: 37.5},
			{WeightMax: 3, UnitPrice: 15},
			{WeightMax: 5, UnitPrice: 10.5},
			{WeightMax: 10, UnitPrice: 7.5},
			{WeightMax: 0, UnitPrice: 6},
		},
		biz.ZoneRemote: {
			{WeightMax: 1, UnitPrice: 45},
			{WeightMax: 3, UnitPrice: 22.5},
			{WeightMax: 5, UnitPrice: 18},
			{WeightMax: 10, UnitPrice: 15},
			{WeightMax: 0, UnitPrice: 12},
		},
		biz.ZoneSpecial: {
			{WeightMax: 1, UnitPrice: 54},
			{WeightMax: 3, UnitPrice: 30},
			{WeightMax: 5, UnitPrice: 24},
			{WeightMax: 10, UnitPrice: 21},
			{WeightMax: 0, UnitPrice: 18},
		},
	},
	// =========================================================================
	// 特惠（经济快递，价格为标快的 0.8 倍）
	// =========================================================================
	biz.ExpressTypeEconomy: {
		biz.ZoneSameCity: {
			{WeightMax: 1, UnitPrice: 9.6},
			{WeightMax: 3, UnitPrice: 1.6},
			{WeightMax: 5, UnitPrice: 1.2},
			{WeightMax: 10, UnitPrice: 0.8},
			{WeightMax: 0, UnitPrice: 0.64},
		},
		biz.ZoneSameProvince: {
			{WeightMax: 1, UnitPrice: 12},
			{WeightMax: 3, UnitPrice: 2.4},
			{WeightMax: 5, UnitPrice: 2},
			{WeightMax: 10, UnitPrice: 1.6},
			{WeightMax: 0, UnitPrice: 1.2},
		},
		biz.ZoneSameCircle: {
			{WeightMax: 1, UnitPrice: 14.4},
			{WeightMax: 3, UnitPrice: 4},
			{WeightMax: 5, UnitPrice: 2.8},
			{WeightMax: 10, UnitPrice: 2},
			{WeightMax: 0, UnitPrice: 1.6},
		},
		biz.ZoneNearby: {
			{WeightMax: 1, UnitPrice: 18.4},
			{WeightMax: 3, UnitPrice: 5.6},
			{WeightMax: 5, UnitPrice: 4},
			{WeightMax: 10, UnitPrice: 3.2},
			{WeightMax: 0, UnitPrice: 2.4},
		},
		biz.ZoneDistant: {
			{WeightMax: 1, UnitPrice: 20},
			{WeightMax: 3, UnitPrice: 8},
			{WeightMax: 5, UnitPrice: 5.6},
			{WeightMax: 10, UnitPrice: 4},
			{WeightMax: 0, UnitPrice: 3.2},
		},
		biz.ZoneRemote: {
			{WeightMax: 1, UnitPrice: 24},
			{WeightMax: 3, UnitPrice: 12},
			{WeightMax: 5, UnitPrice: 9.6},
			{WeightMax: 10, UnitPrice: 8},
			{WeightMax: 0, UnitPrice: 6.4},
		},
		biz.ZoneSpecial: {
			{WeightMax: 1, UnitPrice: 28.8},
			{WeightMax: 3, UnitPrice: 16},
			{WeightMax: 5, UnitPrice: 12.8},
			{WeightMax: 10, UnitPrice: 11.2},
			{WeightMax: 0, UnitPrice: 9.6},
		},
	},
}

// =============================================================================
// DistanceCalculator 高德地图实现
//
// 在运费区域判定中引入真实驾车距离，缓解"纯省份判定"带来的偏差：
//   - 跨省相邻城市（如徐州→淮北）原本按 Zone3 计价，实际仅 ~60km，应更便宜；
//   - 同省巨远（如湛江→韶关）原本按 Zone1，实际 ~600km，应更贵。
// 距离计算失败时（未配置 Key、地理编码/路径规划失败）返回错误，
// 由 biz 层回退到省份分区逻辑，保证服务可用性。
// =============================================================================

// amapDistanceCalculator 基于高德地图计算两城市间驾车距离。
type amapDistanceCalculator struct {
	geo   *amapGeocoder
	cache *distanceCache
}

// NewAmapDistanceCalculator 基于环境变量 AMAP_WEB_KEY 构造距离计算器。
// 未配置 Key 时返回 noop 实现，运费计算将回退到省份分区逻辑。
func NewAmapDistanceCalculator() biz.DistanceCalculator {
	key := os.Getenv("AMAP_WEB_KEY")
	if key == "" {
		return noopDistanceCalculator{}
	}
	return &amapDistanceCalculator{
		geo:   &amapGeocoder{client: amap.NewClient(key)},
		cache: newDistanceCache(),
	}
}

// CityDistanceKM 先对收发城市做地理编码，再调用高德距离测算 API 返回驾车距离（公里）。
func (d *amapDistanceCalculator) CityDistanceKM(ctx context.Context, fromProvince, fromCity, toProvince, toCity string) (float64, error) {
	fromLat, fromLng, err := d.geo.GeocodeAddress(ctx, fromProvince, fromCity, "", "")
	if err != nil {
		return 0, fmt.Errorf("geocode from city: %w", err)
	}
	toLat, toLng, err := d.geo.GeocodeAddress(ctx, toProvince, toCity, "", "")
	if err != nil {
		return 0, fmt.Errorf("geocode to city: %w", err)
	}

	key := fmt.Sprintf("%f,%f|%f,%f", fromLng, fromLat, toLng, toLat)
	if km, ok := d.cache.get(key); ok {
		return km, nil
	}

	resp, err := d.geo.client.Distance(ctx, fmt.Sprintf("%f,%f", fromLng, fromLat), fmt.Sprintf("%f,%f", toLng, toLat), 1) // 1=驾车
	if err != nil {
		return 0, fmt.Errorf("amap distance: %w", err)
	}
	if len(resp.Results) == 0 {
		return 0, fmt.Errorf("empty distance result")
	}
	meters, err := strconv.ParseFloat(resp.Results[0].Distance, 64)
	if err != nil {
		return 0, fmt.Errorf("parse distance %q: %w", resp.Results[0].Distance, err)
	}
	km := meters / 1000.0
	if km <= 0 {
		return 0, fmt.Errorf("invalid distance: %s", resp.Results[0].Distance)
	}
	d.cache.set(key, km)
	return km, nil
}

// noopDistanceCalculator 未配置高德 Key 时的兜底实现。
type noopDistanceCalculator struct{}

func (noopDistanceCalculator) CityDistanceKM(ctx context.Context, fromProvince, fromCity, toProvince, toCity string) (float64, error) {
	return 0, fmt.Errorf("distance calculator not configured (AMAP_WEB_KEY missing)")
}

// distanceCache 进程内距离缓存，避免对同一城市对重复调用地图 API。
type distanceCache struct {
	mu sync.RWMutex
	m  map[string]float64
}

func newDistanceCache() *distanceCache {
	return &distanceCache{m: make(map[string]float64)}
}

func (c *distanceCache) get(key string) (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[key]
	return v, ok
}

func (c *distanceCache) set(key string, km float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = km
}
