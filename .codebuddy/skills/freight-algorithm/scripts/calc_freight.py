#!/usr/bin/env python3
"""
快递运费计算脚本

用法:
    python calc_freight.py --weight 2.5 --length 30 --width 20 --height 15 \
                           --sender "广东,深圳" --receiver "北京,北京" \
                           --express-type 1 --insure-value 1000

    python calc_freight.py --weight 5.0 --sender-province 广东 --sender-city 深圳 \
                           --receiver-province 上海 --receiver-city 上海 \
                           --express-type 1

支持三种快递类型:
    1 = 标快 (Standard Express)
    2 = 特快 (Express Plus)
    3 = 特惠 (Economy)
"""

import argparse
import json
import math
import sys
from dataclasses import dataclass
from typing import Optional


# ============================================================
# 定价规则数据表
# ============================================================

@dataclass
class PricingRule:
    express_type_desc: str
    zone_level: int
    zone_desc: str
    first_weight: float
    first_price: float
    continued_unit: float
    continued_price: float
    volume_factor: int
    rate_multiplier: float
    fuel_surcharge_rate: float
    remote_surcharge_rate: float

# 定价表: (express_type, zone_level) → PricingRule
PRICING_TABLE = {
    # ===== 标快 (type=1) =====
    (1, 0): PricingRule("标快", 0, "同城",   1.0, 8.0,  0.5, 1.0,  6000, 1.0, 0.0, 0.0),
    (1, 1): PricingRule("标快", 1, "省内",   1.0, 10.0, 0.5, 1.5,  6000, 1.0, 0.0, 0.0),
    (1, 2): PricingRule("标快", 2, "经济圈", 1.0, 12.0, 0.5, 2.0,  6000, 1.0, 0.0, 0.0),
    (1, 3): PricingRule("标快", 3, "省外一类",1.0, 15.0, 0.5, 3.0,  6000, 1.0, 0.0, 0.0),
    (1, 4): PricingRule("标快", 4, "省外二类",1.0, 18.0, 0.5, 4.0,  6000, 1.0, 0.0, 0.0),
    (1, 5): PricingRule("标快", 5, "偏远",   1.0, 22.0, 0.5, 6.0,  6000, 1.0, 0.0, 0.15),
    (1, 6): PricingRule("标快", 6, "港澳台", 1.0, 30.0, 0.5, 10.0, 6000, 1.0, 0.0, 0.0),
    # ===== 特快 (type=2) ===== 倍率 1.4
    (2, 0): PricingRule("特快", 0, "同城",   1.0, 11.2, 0.5, 1.4,  6000, 1.4, 0.0, 0.0),
    (2, 1): PricingRule("特快", 1, "省内",   1.0, 14.0, 0.5, 2.1,  6000, 1.4, 0.0, 0.0),
    (2, 2): PricingRule("特快", 2, "经济圈", 1.0, 16.8, 0.5, 2.8,  6000, 1.4, 0.0, 0.0),
    (2, 3): PricingRule("特快", 3, "省外一类",1.0, 21.0, 0.5, 4.2,  6000, 1.4, 0.0, 0.0),
    (2, 4): PricingRule("特快", 4, "省外二类",1.0, 25.2, 0.5, 5.6,  6000, 1.4, 0.0, 0.0),
    (2, 5): PricingRule("特快", 5, "偏远",   1.0, 30.8, 0.5, 8.4,  6000, 1.4, 0.0, 0.15),
    (2, 6): PricingRule("特快", 6, "港澳台", 1.0, 42.0, 0.5, 14.0, 6000, 1.4, 0.0, 0.0),
    # ===== 特惠 (type=3) ===== 倍率 0.7, 体积系数 12000
    (3, 0): PricingRule("特惠", 0, "同城",   1.0, 5.6,  1.0, 1.4,  12000, 0.7, 0.0, 0.0),
    (3, 1): PricingRule("特惠", 1, "省内",   1.0, 7.0,  1.0, 2.1,  12000, 0.7, 0.0, 0.0),
    (3, 2): PricingRule("特惠", 2, "经济圈", 1.0, 8.4,  1.0, 2.8,  12000, 0.7, 0.0, 0.0),
    (3, 3): PricingRule("特惠", 3, "省外一类",1.0, 10.5, 1.0, 4.2,  12000, 0.7, 0.0, 0.0),
    (3, 4): PricingRule("特惠", 4, "省外二类",1.0, 12.6, 1.0, 5.6,  12000, 0.7, 0.0, 0.0),
    (3, 5): PricingRule("特惠", 5, "偏远",   1.0, 15.4, 1.0, 8.4,  12000, 0.7, 0.0, 0.15),
    (3, 6): PricingRule("特惠", 6, "港澳台", 1.0, 21.0, 1.0, 14.0, 12000, 0.7, 0.0, 0.0),
}

# 经济圈定义
ECONOMIC_CIRCLES = {
    "长三角": {"上海", "江苏", "浙江", "安徽"},
    "珠三角": {"广东", "广西", "海南"},
    "京津冀": {"北京", "天津", "河北"},
    "成渝":   {"四川", "重庆"},
}

# 省份分类
PROVINCE_ZONES = {
    **{p: 3 for p in ["福建", "山东", "湖南", "湖北", "河南", "江西", "陕西", "山西"]},
    **{p: 4 for p in ["云南", "贵州", "黑龙江", "吉林", "辽宁"]},
    **{p: 5 for p in ["新疆", "西藏", "青海", "内蒙古", "宁夏", "甘肃"]},
}

REMOTE_PROVINCES = {"新疆", "西藏", "青海", "内蒙古", "宁夏", "甘肃"}
HK_MO_TW = {"香港", "澳门", "台湾"}


def determine_zone_level(sender_province: str, sender_city: str,
                         receiver_province: str, receiver_city: str) -> int:
    """根据寄收双方地址判定区域等级"""
    # 港澳台
    if receiver_province in HK_MO_TW:
        return 6
    if sender_province in HK_MO_TW:
        return 6

    # 同城
    if sender_city == receiver_city:
        return 0

    # 省内
    if sender_province == receiver_province:
        return 1

    # 经济圈
    for circle_name, provinces in ECONOMIC_CIRCLES.items():
        if sender_province in provinces and receiver_province in provinces:
            return 2

    # 偏远地区
    if receiver_province in REMOTE_PROVINCES or sender_province in REMOTE_PROVINCES:
        return 5

    # 省外判定
    zone = PROVINCE_ZONES.get(receiver_province, 3)
    sender_zone = PROVINCE_ZONES.get(sender_province, 3)
    return max(zone, sender_zone)


def calc_chargeable_weight(weight: float, length: int, width: int, height: int,
                           volume_factor: int) -> float:
    """计算计费重量（实际重量 vs 体积重量取大值）"""
    if length > 0 and width > 0 and height > 0:
        volume_weight = (length * width * height) / volume_factor
    else:
        volume_weight = 0
    return max(weight, volume_weight)


def calc_base_freight(calc_weight: float, rule: PricingRule) -> float:
    """计算基础运费（首重 + 续重）"""
    if calc_weight <= rule.first_weight:
        return rule.first_price

    continued_weight = calc_weight - rule.first_weight
    continued_units = math.ceil(continued_weight / rule.continued_unit)
    return rule.first_price + continued_units * rule.continued_price


def calc_insure_fee(insure_value: float) -> float:
    """计算保价费：保价金额 × 0.5%，最低 1 元"""
    if insure_value <= 0:
        return 0.0
    fee = insure_value * 0.005
    return max(fee, 1.0)


def estimate_freight(
    weight: float,
    express_type: int,
    sender_province: str,
    sender_city: str,
    receiver_province: str,
    receiver_city: str,
    sender_area: str = "",
    receiver_area: str = "",
    length: int = 0,
    width: int = 0,
    height: int = 0,
    insure_value: float = 0.0,
) -> dict:
    """
    预估运费主函数。

    返回:
        {
            "base_freight": float,   # 基础运费
            "insure_fee": float,     # 保价费
            "surcharge": float,      # 附加费
            "total_price": float,    # 总费用
            "calc_weight": float,    # 计费重量
            "zone_level": int,       # 区域等级
            "zone_desc": str,        # 区域描述
            "express_type_desc": str, # 快递类型描述
            "tips": str,             # 提示
        }
    """
    # 参数校验
    if weight <= 0:
        return {"error": "包裹重量必须大于 0"}
    if express_type not in (1, 2, 3):
        return {"error": "快递类型无效，有效值: 1(标快) 2(特快) 3(特惠)"}

    # 确定区域等级
    zone_level = determine_zone_level(
        sender_province, sender_city, receiver_province, receiver_city
    )

    # 获取定价规则
    key = (express_type, zone_level)
    if key not in PRICING_TABLE:
        return {"error": f"未找到定价规则: express_type={express_type}, zone={zone_level}"}

    rule = PRICING_TABLE[key]

    # 计算计费重量
    calc_weight = calc_chargeable_weight(weight, length, width, height, rule.volume_factor)
    calc_weight_rounded = math.ceil(calc_weight) if calc_weight > 0 else calc_weight

    # 计算基础运费
    base_freight = calc_base_freight(calc_weight_rounded, rule)

    # 附加费
    surcharge = base_freight * (rule.fuel_surcharge_rate + rule.remote_surcharge_rate)

    # 保价费
    insure_fee = calc_insure_fee(insure_value)

    # 总价
    total_price = base_freight + surcharge + insure_fee

    # 生成 tips
    tips_parts = []
    if length > 0 and width > 0 and height > 0:
        volume_weight = (length * width * height) / rule.volume_factor
        if volume_weight > weight:
            tips_parts.append(f"体积重量({volume_weight:.1f}kg)超过实际重量({weight}kg)，按体积重量计费")
        else:
            tips_parts.append(f"按实际重量({weight}kg)计费")
    else:
        tips_parts.append(f"按实际重量({weight}kg)计费，未提供体积信息")

    if insure_value > 0:
        tips_parts.append(f"保价{insure_value}元，保价费{insure_fee:.2f}元")
    else:
        tips_parts.append("未选择保价服务")

    if rule.remote_surcharge_rate > 0:
        tips_parts.append("偏远地区附加费已包含")

    return {
        "base_freight": round(base_freight, 2),
        "insure_fee": round(insure_fee, 2),
        "surcharge": round(surcharge, 2),
        "total_price": round(total_price, 2),
        "calc_weight": round(calc_weight, 2),
        "zone_level": zone_level,
        "zone_desc": rule.zone_desc,
        "express_type_desc": rule.express_type_desc,
        "tips": "；".join(tips_parts),
    }


def main():
    parser = argparse.ArgumentParser(
        description="快递运费计算工具",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  python calc_freight.py --weight 2.5 --sender-province 广东 --sender-city 深圳 \\
                         --receiver-province 北京 --receiver-city 北京 \\
                         --express-type 1

  python calc_freight.py --weight 5.0 --length 30 --width 20 --height 15 \\
                         --sender-province 上海 --sender-city 上海 \\
                         --receiver-province 西藏 --receiver-city 拉萨 \\
                         --express-type 2 --insure-value 2000
        """,
    )
    parser.add_argument("--weight", type=float, required=True, help="包裹重量(kg)")
    parser.add_argument("--length", type=int, default=0, help="包裹长(cm)")
    parser.add_argument("--width", type=int, default=0, help="包裹宽(cm)")
    parser.add_argument("--height", type=int, default=0, help="包裹高(cm)")
    parser.add_argument("--sender-province", required=True, help="寄件省份")
    parser.add_argument("--sender-city", required=True, help="寄件城市")
    parser.add_argument("--sender-area", default="", help="寄件区县")
    parser.add_argument("--receiver-province", required=True, help="收件省份")
    parser.add_argument("--receiver-city", required=True, help="收件城市")
    parser.add_argument("--receiver-area", default="", help="收件区县")
    parser.add_argument("--express-type", type=int, required=True,
                        choices=[1, 2, 3], help="快递类型: 1标快 2特快 3特惠")
    parser.add_argument("--insure-value", type=float, default=0.0, help="保价金额(元)")
    parser.add_argument("--json", action="store_true", help="以 JSON 格式输出")

    args = parser.parse_args()

    result = estimate_freight(
        weight=args.weight,
        express_type=args.express_type,
        sender_province=args.sender_province,
        sender_city=args.sender_city,
        receiver_province=args.receiver_province,
        receiver_city=args.receiver_city,
        sender_area=args.sender_area,
        receiver_area=args.receiver_area,
        length=args.length,
        width=args.width,
        height=args.height,
        insure_value=args.insure_value,
    )

    if "error" in result:
        print(f"❌ 错误: {result['error']}", file=sys.stderr)
        sys.exit(1)

    if args.json:
        print(json.dumps(result, ensure_ascii=False, indent=2))
    else:
        print("=" * 50)
        print("  快递运费预估结果")
        print("=" * 50)
        print(f"  快递类型:   {result['express_type_desc']}")
        print(f"  区域等级:   {result['zone_desc']}(Level {result['zone_level']})")
        print(f"  计费重量:   {result['calc_weight']} kg")
        print(f"  基础运费:   ¥{result['base_freight']:.2f}")
        if result['surcharge'] > 0:
            print(f"  附加费:     ¥{result['surcharge']:.2f}")
        if result['insure_fee'] > 0:
            print(f"  保价费:     ¥{result['insure_fee']:.2f}")
        print(f"  {'─' * 20}")
        print(f"  总费用:     ¥{result['total_price']:.2f}")
        print(f"  提示:       {result['tips']}")
        print("=" * 50)


if __name__ == "__main__":
    main()
