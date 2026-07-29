---
name: freight-algorithm
description: >
  This skill provides expert knowledge of express shipping freight cost calculation
  algorithms, pricing models, and Kratos framework integration patterns. It should
  be used when implementing freight estimation, shipping cost calculation,
  express delivery pricing logic, or any task involving "运费计算", "运费预估",
  "快递计价", "freight estimate", "shipping cost", or "express pricing".
  It is specifically designed for the "shunfeng-miniprogram" Kratos v3 project
  but the algorithms apply to any express/logistics system.
---

# Freight Algorithm (运费算法)

## Overview

This skill enables implementation of express delivery freight cost calculation in a Kratos v3 microservice.
It covers the complete pricing lifecycle: weight calculation (actual vs volumetric), zone-based tiered pricing,
first-weight / continued-weight models, insurance fee computation, and surcharge rules across three express
types (Standard Express / Express Plus / Economy).

## When to Use This Skill

Trigger this skill when the task involves:

- Implementing or modifying the `EstimateFreight` RPC in `api/order/v2/`
- Adding shipping cost calculation to any express/logistics project
- Designing freight pricing tables or rate cards
- Debugging freight calculation discrepancies
- Extending the pricing model with new surcharges or express types

## Core Calculation Flow

Every freight estimate follows this sequence:

```
Input Parameters → Zone Determination → Pricing Lookup → Weight Calculation → Base Freight → Surcharges → Total
```

### Step 1: Zone Determination

Determine the shipping zone level by comparing sender and receiver addresses:

| Condition | Zone Level | Description |
|-----------|------------|-------------|
| Same city | 0 | 同城 |
| Same province | 1 | 省内 |
| Same economic circle | 2 | 经济圈 (e.g. 长三角, 珠三角, 京津冀) |
| Adjacent provinces | 3 | 省外一类 |
| Distant provinces | 4 | 省外二类 |
| Remote areas (新疆/西藏/青海/内蒙古/宁夏/甘肃) | 5 | 偏远地区 |
| Hong Kong/Macau/Taiwan | 6 | 港澳台 |

Load the reference file `references/express-pricing.md` for the complete province classification table
and economic circle definitions.

### Step 2: Retrieve Pricing Rules

Pricing rules are keyed by `(express_type, zone_level)`. Each rule specifies:
- `first_weight` / `first_price` — first weight threshold and price
- `continued_unit` / `continued_price` — continued weight unit and price per unit
- `volume_factor` — divisor for volumetric weight (6000 for air, 12000 for ground)
- `rate_multiplier` — express type premium (1.0 standard, 1.4 express plus, 0.7 economy)
- `fuel_surcharge_rate` / `remote_surcharge_rate` — additional surcharge rates

Complete pricing tables are in `references/express-pricing.md`.

### Step 3: Calculate Chargeable Weight

```
volume_weight = (length × width × height) / volume_factor   (when dimensions provided)
chargeable_weight = MAX(actual_weight, volume_weight)
chargeable_weight_rounded = CEIL(chargeable_weight)          (round up to nearest kg)
```

When `length`, `width`, or `height` is 0, skip volumetric weight and use actual weight directly.

### Step 4: Calculate Base Freight

```
if chargeable_weight ≤ first_weight:
    base_freight = first_price
else:
    continued = chargeable_weight - first_weight
    continued_units = CEIL(continued / continued_unit)
    base_freight = first_price + continued_units × continued_price
```

### Step 5: Calculate Surcharges & Insurance

```
insurance_fee = IF insure_value > 0 THEN MAX(insure_value × 0.005, 1.0) ELSE 0
fuel_surcharge = base_freight × fuel_surcharge_rate
remote_surcharge = base_freight × remote_surcharge_rate
```

Default insurance rate is 0.5% with a minimum of ¥1.00.

### Step 6: Total

```
total_price = base_freight + insurance_fee + fuel_surcharge + remote_surcharge
```

## Express Type Rules

| Type | Name | Volume Factor | Continued Unit | Rate Multiplier |
|------|------|---------------|----------------|-----------------|
| 1 | 标快 (Standard) | 6000 | 0.5 kg | 1.0× |
| 2 | 特快 (Express Plus) | 6000 | 0.5 kg | 1.4× |
| 3 | 特惠 (Economy) | 12000 | 1.0 kg | 0.7× |

## Kratos Architecture Integration

When implementing freight estimation in the `shunfeng-miniprogram` project,
follow the standard Kratos three-layer pattern. Detailed integration guidance
is in `references/kratos-integration.md`.

### Files to create:

| Layer | File | Responsibility |
|-------|------|----------------|
| biz | `internal/biz/freight.go` | DO definitions, Repo interface, Usecase with core calculation |
| data | `internal/data/freight.go` | Repo implementation, pricing data access |
| service | `internal/service/freight.go` | DTO ↔ DO conversion, request validation |

### Quick implementation checklist:

1. Create `internal/biz/freight.go` — define `FreightRequest`, `FreightCharge`, `PricingRule` DOs;
   declare `FreightRepo` interface; implement `FreightUsecase.Estimate()` with the calculation
   logic from this SKILL.md.
2. Create `internal/data/freight.go` — implement `FreightRepo` backed by an in-memory
   pricing table or database; register in `data.ProviderSet`.
3. Create `internal/service/freight.go` — adapt `EstimateFreightRequest` proto DTO →
   `FreightRequest` DO; call usecase; return `EstimateFreightReply`.
4. Wire everything: add constructors to `biz.ProviderSet`, `service.ProviderSet`;
   register `order.v2.OrderService` in `internal/server/http.go` and `grpc.go`.
5. Run `make all` to regenerate Wire and verify compilation.

### Validation rules (service layer):

- `weight` must be > 0
- `express_type` must be 1, 2, or 3
- `sender_province` and `receiver_province` are required
- `sender_city` and `receiver_city` are required
- `insure_value` must be ≥ 0 (0 = no insurance)
- Return `ErrFreightInvalidArg` on validation failure

## Using the Calculation Script

The bundled `scripts/calc_freight.py` can be used to test and verify pricing rules independently:

```bash
# Basic calculation
python scripts/calc_freight.py \
  --weight 2.5 \
  --sender-province 广东 --sender-city 深圳 \
  --receiver-province 北京 --receiver-city 北京 \
  --express-type 1

# With dimensions and insurance
python scripts/calc_freight.py \
  --weight 5.0 --length 30 --width 20 --height 15 \
  --sender-province 上海 --sender-city 上海 \
  --receiver-province 西藏 --receiver-city 拉萨 \
  --express-type 2 --insure-value 2000

# JSON output
python scripts/calc_freight.py --weight 1.0 \
  --sender-province 广东 --sender-city 广州 \
  --receiver-province 广东 --receiver-city 广州 \
  --express-type 1 --json
```

The script serves as a reference implementation that can be ported to Go for the actual biz layer logic.

## Resources

### references/express-pricing.md
Complete reference for express pricing models, zone mappings, surcharge rules, and the full
freight calculation formula. Load this when detailed pricing tables or algorithm edge cases
are needed.

### references/kratos-integration.md
Detailed integration guide for the `shunfeng-miniprogram` project, covering the exact files
to create, code patterns for each layer, Wire wiring instructions, and the database schema
for production deployment.

### scripts/calc_freight.py
Standalone Python script that implements the full freight calculation algorithm. Use it to
test pricing rules independently or as a reference when porting the logic to Go. The script
accepts command-line arguments matching the proto `EstimateFreightRequest` fields.
