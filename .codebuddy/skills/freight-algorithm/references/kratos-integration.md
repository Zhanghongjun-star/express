# Kratos 架构集成指南 (运费算法)

## 项目上下文

此 skill 设计用于 `shunfeng-miniprogram` (顺丰快递小程序后端)，基于 Kratos v3 框架。

## 当前状态

Order v2 运费预估 (`EstimateFreight`) 的 Proto 定义和生成代码已就绪，但缺少三层业务实现：

```
已就绪:
api/order/v2/order.proto       ✓ Proto 定义
api/order/v2/error_reason.proto ✓ 错误枚举
api/order/v2/*.pb.go            ✓ 生成代码

待实现:
internal/biz/freight.go         ✗ 领域对象 + 用例 + 仓库接口
internal/data/freight.go        ✗ 仓库实现
internal/service/freight.go     ✗ 传输适配层
internal/server/http.go          ✗ 注册 v2 服务
internal/server/grpc.go          ✗ 注册 v2 服务
```

## 实现分层指引

### 1. Biz 层 (`internal/biz/freight.go`)

#### 领域对象 (DO)

```go
package biz

import (
    "context"
    v2 "shunfeng-miniprogram/api/order/v2"
    "github.com/go-kratos/kratos/v3/errors"
)

// FreightCharge 运费计算结果（领域对象，不引入 proto）。
type FreightCharge struct {
    BaseFreight float64 // 基础运费
    InsureFee   float64 // 保价费
    TotalPrice  float64 // 总费用
    CalcWeight  float64 // 计费重量
    Tips        string  // 提示信息
}

// FreightRequest 运费计算请求（领域对象）。
type FreightRequest struct {
    SenderProvince   string
    SenderCity       string
    SenderArea       string
    ReceiverProvince string
    ReceiverCity     string
    ReceiverArea     string
    Weight           float64
    Length           int32
    Width            int32
    Height           int32
    ExpressType      int32
    InsureValue      float64
}
```

#### 错误定义

```go
var (
    ErrFreightCalcFailed     = errors.InternalServer(v2.ErrorReason_ORDER_FREIGHT_CALC_FAILED.String(), "freight calculation failed")
    ErrFreightInvalidArg     = errors.BadRequest(v2.ErrorReason_ORDER_FREIGHT_INVALID_ARGUMENT.String(), "invalid freight argument")
)
```

#### 仓库接口

```go
// FreightRepo 运费仓库接口 — 用于获取定价规则数据。
type FreightRepo interface {
    // GetPricingRule 获取指定快递类型和区域等级的定价规则。
    GetPricingRule(ctx context.Context, expressType int32, zoneLevel int) (*PricingRule, error)
}

// PricingRule 定价规则领域对象。
type PricingRule struct {
    ExpressType   int32   // 快递类型
    ZoneLevel     int     // 区域等级：0同城 1省内 2经济圈 3省外一类 4省外二类 5省外三类 6港澳台
    FirstWeight   float64 // 首重(kg)
    FirstPrice    float64 // 首重价格(元)
    ContinuedUnit float64 // 续重单位(kg)
    ContinuedPrice float64 // 续重单价(元/单位)
    VolumeFactor  int32   // 体积系数
    RateMultiplier float64 // 类型倍率(标快=1.0, 特快=1.4, 特惠=0.7)
    FuelSurchargeRate float64 // 燃油附加费率
    RemoteSurchargeRate float64 // 偏远附加费率
}
```

#### 用例 (Usecase)

用例包含核心运费计算逻辑所在。实现 `Estimate` 方法：
- 验证输入参数
- 获取定价规则
- 计算计费重量（体积重量 vs 实际重量取大值）
- 计算基础运费（首重 + 续重）
- 计算附加费（保价费、偏远附加费）
- 汇总并返回 `FreightCharge`

#### 区域等级判断逻辑

区域等级根据寄收双方地址判定：

```
若 sender_city == receiver_city              → ZoneLevel 0 (同城)
若 sender_province == receiver_province       → ZoneLevel 1 (省内)
若 同经济圈                                   → ZoneLevel 2 (经济圈)
若 邻近省份                                   → ZoneLevel 3 (省外一类)
若 远距离省份                                 → ZoneLevel 4 (省外二类)
若 偏远地区(藏/疆/青/蒙/宁/甘)                → ZoneLevel 5 (省外三类)
若 港澳台                                     → ZoneLevel 6 (港澳台)
```

### 2. Data 层 (`internal/data/freight.go`)

实现 `biz.FreightRepo`，提供定价数据访问。可以用以下方式：

**方案一：本地配置表（推荐起步方案）**

```go
package data

import "shunfeng-miniprogram/internal/biz"

func NewFreightRepo(d *Data) biz.FreightRepo {
    return &freightRepo{data: d}
}
```

定价规则存储为内嵌数据表，按 (expressType, zoneLevel) 组合索引。

**方案二：数据库表（生产方案）**

建表：

```sql
CREATE TABLE sf_freight_pricing (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    express_type TINYINT NOT NULL COMMENT '快递类型:1标快2特快3特惠',
    zone_level TINYINT NOT NULL COMMENT '区域等级:0同城1省内2经济圈3省外一类4省外二类5省外三类6港澳台',
    first_weight DOUBLE NOT NULL COMMENT '首重(kg)',
    first_price DOUBLE NOT NULL COMMENT '首重价格(元)',
    continued_unit DOUBLE NOT NULL COMMENT '续重单位(kg)',
    continued_price DOUBLE NOT NULL COMMENT '续重单价(元/单位)',
    volume_factor INT NOT NULL DEFAULT 6000 COMMENT '体积系数',
    rate_multiplier DOUBLE NOT NULL DEFAULT 1.0 COMMENT '快递类型倍率',
    fuel_surcharge_rate DOUBLE NOT NULL DEFAULT 0.0 COMMENT '燃油附加费率',
    remote_surcharge_rate DOUBLE NOT NULL DEFAULT 0.0 COMMENT '偏远附加费率',
    UNIQUE KEY uk_express_zone (express_type, zone_level)
);
```

### 3. Service 层 (`internal/service/freight.go`)

将 proto DTO 转为 biz DO，调用 usecase，将结果转为 proto reply：

```go
package service

import (
    "context"
    v2 "shunfeng-miniprogram/api/order/v2"
    "shunfeng-miniprogram/internal/biz"
)

type FreightService struct {
    v2.UnimplementedOrderServiceServer
    uc *biz.FreightUsecase
}

func NewFreightService(uc *biz.FreightUsecase) *FreightService {
    return &FreightService{uc: uc}
}

func (s *FreightService) EstimateFreight(ctx context.Context, req *v2.EstimateFreightRequest) (*v2.EstimateFreightReply, error) {
    // 参数验证
    if req.Weight <= 0 || req.ExpressType < 1 || req.ExpressType > 3 {
        return nil, biz.ErrFreightInvalidArg
    }
    // 转换并调用用例
    charge, err := s.uc.Estimate(ctx, &biz.FreightRequest{...})
    if err != nil {
        return nil, err
    }
    return &v2.EstimateFreightReply{...}, nil
}
```

### 4. 依赖注入 (Wire)

在以下文件的 ProviderSet 中添加：

- `internal/data/data.go` → `data.ProviderSet` 添加 `NewFreightRepo`
- `internal/biz/biz.go` → `biz.ProviderSet` 添加 `NewFreightUsecase`
- `internal/service/service.go` → `service.ProviderSet` 添加 `NewFreightService`
- `internal/server/http.go` → 注册 `v2.RegisterOrderServiceHTTPServer`
- `internal/server/grpc.go` → 注册 `v2.RegisterOrderServiceServer`

### 5. 验证运行

```bash
cd d:/code/express
make all    # 或手动：go generate ./... && go build ./cmd/server/
```

---

## 注意事项

1. **分层约束**：service 不直接访问 data；biz 不引入 proto struct；data 不暴露存储驱动类型
2. **错误处理**：参数校验在 service 层做快速失败检查，核心校验在 biz 层
3. **体积重量**：当 weight > 0 但 length/width/height 未传时，直接使用实际重量作为计费重量
4. **保价费**：insure_value=0 表示不保价，保价费直接返回 0
5. **Decimal 精度**：金额使用 float64，展示时由前端保留两位小数
