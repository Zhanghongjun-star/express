package biz

import (
	"context"
	"encoding/json"
	"testing"
)

func TestCalcInsuranceFee(t *testing.T) {
	cases := []struct {
		name    string
		value   float64
		fragile bool
		want    float64
		wantErr bool
	}{
		// 普通物品
		{"普通-0元", 0, false, 0, false},
		{"普通-500元-固定1", 500, false, 1, false},
		{"普通-600元-固定2", 600, false, 2, false},
		{"普通-1000元-固定2", 1000, false, 2, false},
		{"普通-1200元-0.5%", 1200, false, 6, false},
		{"普通-6000元-0.5%", 6000, false, 30, false},
		// 易碎物品
		{"易碎-500元-固定3", 500, true, 3, false},
		{"易碎-1000元-固定7", 1000, true, 7, false},
		{"易碎-1200元-1.5%", 1200, true, 18, false},
		{"易碎-6000元-1.5%", 6000, true, 90, false},
		// 四舍五入
		{"普通-1234元-四舍五入", 1234, false, 6, false},
		// 超过上限
		{"超过2万上限", 20001, false, 0, true},
		{"恰好2万", 20000, false, 100, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// 新签名返回 (费用, 公式, 错误)
			got, _, err := calcInsuranceFee(c.value, c.fragile)
			if c.wantErr {
				if err == nil {
					t.Fatalf("期望返回错误，但得到 nil (value=%.0f)", c.value)
				}
				return
			}
			if err != nil {
				t.Fatalf("不期望错误: %v", err)
			}
			if got != c.want {
				t.Fatalf("calcInsuranceFee(%.0f, %v) = %.0f, want %.0f", c.value, c.fragile, got, c.want)
			}
		})
	}
}

// sampleRates 测试用示例费率表（id -> 材料单价），对应默认种子材料。
var sampleRates = map[int64]*materialRate{
	1:  {"文件封", 1},
	2:  {"1号纸箱", 2},
	3:  {"2号纸箱", 3},
	4:  {"3号纸箱", 4},
	5:  {"气泡膜", 5},
	6:  {"珍珠棉", 6},
	7:  {"防水袋", 1},
	8:  {"缠绕膜", 4},
	9:  {"木架", 20},
	10: {"木箱", 30},
}

func TestCalcPackageFee(t *testing.T) {
	cases := []struct {
		name    string
		items   []PackageItem
		want    float64
		wantErr bool
	}{
		{"空列表", nil, 0, false},
		{"单个气泡膜(id=5)x2", []PackageItem{{MaterialID: 5, Quantity: 2}}, 10, false},
		{"混合：1号纸箱(id=2)x1 + 木架(id=9)x1", []PackageItem{{MaterialID: 2, Quantity: 1}, {MaterialID: 9, Quantity: 1}}, 22, false},
		{"数量缺省按1", []PackageItem{{MaterialID: 1}}, 1, false},
		{"未知材料(id=999)报错", []PackageItem{{MaterialID: 999}}, 0, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// 传入单价表（来自 packaging_material 表，此处用示例费率）
			got, _, err := calcPackageFee(c.items, sampleRates)
			if c.wantErr {
				if err == nil {
					t.Fatalf("期望返回错误，但得到 nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("不期望错误: %v", err)
			}
			if got != c.want {
				t.Fatalf("calcPackageFee(%v) = %.0f, want %.0f", c.items, got, c.want)
			}
		})
	}
}

// TestPackageFeeDemo 演示：列出可选包装材料种类，并打印若干计算实例与过程明细。
func TestPackageFeeDemo(t *testing.T) {
	// 1) 统计可选材料种类
	t.Logf("可选包装材料共 %d 种：", len(sampleRates))
	for id, r := range sampleRates {
		t.Logf("  - id=%d %s：%.0f 元/单位", id, r.Name, r.UnitPrice)
	}
	if len(sampleRates) != 10 {
		t.Fatalf("期望可选 10 种，实际 %d 种", len(sampleRates))
	}

	// 2) 演示若干组合的计算结果
	demos := []struct {
		name  string
		items []PackageItem
	}{
		{"仅气泡膜 x3", []PackageItem{{MaterialID: 5, Quantity: 3}}},
		{"1号纸箱 x1 + 珍珠棉 x2", []PackageItem{{MaterialID: 2, Quantity: 1}, {MaterialID: 6, Quantity: 2}}},
		{"木架 x1 + 木箱 x1 + 防水袋 x2", []PackageItem{{MaterialID: 9, Quantity: 1}, {MaterialID: 10, Quantity: 1}, {MaterialID: 7, Quantity: 2}}},
	}
	for _, d := range demos {
		fee, lines, err := calcPackageFee(d.items, sampleRates)
		if err != nil {
			t.Fatalf("%s 计算失败: %v", d.name, err)
		}
		detail, _ := json.Marshal(feeBreakdown{Packages: lines, PackageFee: fee})
		t.Logf("[%s] 打包费 = %.0f 元 | 明细: %s", d.name, fee, string(detail))
	}
}

// ==================== CreateExpressOrder 主流程测试 ====================

// fakeOrderRepo 是 ExpressOrderRepo 的内存假实现，记录最后一次创建的订单。
type fakeOrderRepo struct {
	last   *ExpressOrder
	byNo   map[string]*ExpressOrder
}

func newFakeOrderRepo() *fakeOrderRepo {
	return &fakeOrderRepo{byNo: map[string]*ExpressOrder{}}
}

func (r *fakeOrderRepo) Create(_ context.Context, o *ExpressOrder) (*ExpressOrder, error) {
	r.last = o
	r.byNo[o.OrderNo] = o
	return o, nil
}
func (r *fakeOrderRepo) GetByNo(_ context.Context, orderNo string) (*ExpressOrder, error) {
	if o, ok := r.byNo[orderNo]; ok {
		return o, nil
	}
	return nil, ErrExpressOrderNotFound
}
func (r *fakeOrderRepo) List(_ context.Context, _, _ string, _ ...ListOption) ([]*ExpressOrder, int32, error) {
	out := make([]*ExpressOrder, 0, len(r.byNo))
	for _, o := range r.byNo {
		out = append(out, o)
	}
	return out, int32(len(out)), nil
}
func (r *fakeOrderRepo) Update(_ context.Context, o *ExpressOrder, _ []string) (*ExpressOrder, error) {
	r.byNo[o.OrderNo] = o
	return o, nil
}
func (r *fakeOrderRepo) Delete(_ context.Context, orderNo string) error {
	delete(r.byNo, orderNo)
	return nil
}

// fakeMaterialRepo 是 PackagingMaterialRepo 的内存假实现（按 id/名称查材料）。
type fakeMaterialRepo struct {
	byID   map[int64]*PackagingMaterial
	byName map[string]*PackagingMaterial
}

func newFakeMaterialRepo() *fakeMaterialRepo {
	ms := []*PackagingMaterial{
		{ID: 1, Name: "文件封", UnitPrice: 1, Unit: "个"},
		{ID: 2, Name: "1号纸箱", UnitPrice: 2, Unit: "个"},
		{ID: 5, Name: "气泡膜", UnitPrice: 5, Unit: "卷"},
		{ID: 9, Name: "木架", UnitPrice: 20, Unit: "件"},
	}
	byID := map[int64]*PackagingMaterial{}
	byName := map[string]*PackagingMaterial{}
	for _, m := range ms {
		byID[m.ID] = m
		byName[m.Name] = m
	}
	return &fakeMaterialRepo{byID: byID, byName: byName}
}

func (r *fakeMaterialRepo) Create(_ context.Context, m *PackagingMaterial) (*PackagingMaterial, error) {
	return m, nil
}
func (r *fakeMaterialRepo) GetByName(_ context.Context, name string) (*PackagingMaterial, error) {
	if m, ok := r.byName[name]; ok {
		return m, nil
	}
	return nil, ErrPackagingMaterialNotFound
}
func (r *fakeMaterialRepo) GetByID(_ context.Context, id int64) (*PackagingMaterial, error) {
	if m, ok := r.byID[id]; ok {
		return m, nil
	}
	return nil, ErrPackagingMaterialNotFound
}
func (r *fakeMaterialRepo) List(_ context.Context, _ ...ListOption) ([]*PackagingMaterial, int32, error) {
	out := make([]*PackagingMaterial, 0, len(r.byID))
	for _, m := range r.byID {
		out = append(out, m)
	}
	return out, int32(len(out)), nil
}
func (r *fakeMaterialRepo) Update(_ context.Context, m *PackagingMaterial, _ []string) (*PackagingMaterial, error) {
	return m, nil
}
func (r *fakeMaterialRepo) Delete(_ context.Context, _ int64) error {
	return nil
}

// TestCreateExpressOrder 验证下单主流程：费用汇总、过程明细、默认值与错误分支。
func TestCreateExpressOrder(t *testing.T) {
	t.Run("正常下单并正确汇总费用", func(t *testing.T) {
		uc := NewExpressOrderUsecase(newFakeOrderRepo(), newFakeMaterialRepo())

		// 预期：运费 10 + 打包费(气泡膜x2=10 + 木架x1=20=30) + 保价(普通品 1000元 固定2) = 42
		o := &ExpressOrder{
			FreightFee:   10,
			DeclaredValue: 1000,
			IsFragile:    false,
			Packages: []PackageItem{
				{MaterialID: 5, Quantity: 2},
				{MaterialID: 9, Quantity: 1},
			},
		}
		got, err := uc.CreateExpressOrder(context.Background(), o)
		if err != nil {
			t.Fatalf("不期望错误: %v", err)
		}

		// 费用汇总
		if got.PackageFee != 30 {
			t.Errorf("PackageFee = %.0f, want 30", got.PackageFee)
		}
		if got.InsuranceFee != 2 {
			t.Errorf("InsuranceFee = %.0f, want 2", got.InsuranceFee)
		}
		if got.TotalFee != 42 {
			t.Errorf("TotalFee = %.0f, want 42", got.TotalFee)
		}
		// 过程明细应已生成
		if got.PackageDetail == "" {
			t.Error("PackageDetail 不应为空")
		}
		// 默认值填充
		if got.OrderNo == "" {
			t.Error("OrderNo 不应为空")
		}
		if got.OrderStatus != "pending" {
			t.Errorf("OrderStatus = %q, want pending", got.OrderStatus)
		}
		if got.PayStatus != "unpaid" {
			t.Errorf("PayStatus = %q, want unpaid", got.PayStatus)
		}
	})

	t.Run("易碎品保价费率生效", func(t *testing.T) {
		uc := NewExpressOrderUsecase(newFakeOrderRepo(), newFakeMaterialRepo())
		o := &ExpressOrder{
			FreightFee:    10,
			DeclaredValue: 1000,
			IsFragile:     true, // 易碎品 1000 元固定 7 元
			Packages:      []PackageItem{{MaterialID: 1, Quantity: 1}},
		}
		got, err := uc.CreateExpressOrder(context.Background(), o)
		if err != nil {
			t.Fatalf("不期望错误: %v", err)
		}
		if got.InsuranceFee != 7 {
			t.Errorf("易碎品 InsuranceFee = %.0f, want 7", got.InsuranceFee)
		}
		if got.TotalFee != 10+1+7 { // 运费10 + 文件封1 + 保价7
			t.Errorf("TotalFee = %.0f, want 18", got.TotalFee)
		}
	})

	t.Run("未知材料id报错误", func(t *testing.T) {
		uc := NewExpressOrderUsecase(newFakeOrderRepo(), newFakeMaterialRepo())
		o := &ExpressOrder{
			FreightFee: 10,
			Packages:   []PackageItem{{MaterialID: 999, Quantity: 1}},
		}
		_, err := uc.CreateExpressOrder(context.Background(), o)
		if err == nil {
			t.Fatal("期望因未知材料 id 返回错误，但得到 nil")
		}
		t.Logf("正确拦截未知材料: %v", err)
	})

	t.Run("保价超额报错误", func(t *testing.T) {
		uc := NewExpressOrderUsecase(newFakeOrderRepo(), newFakeMaterialRepo())
		o := &ExpressOrder{
			FreightFee:    10,
			DeclaredValue: 25000, // 超过单票上限 20000
			Packages:      []PackageItem{{MaterialID: 1, Quantity: 1}},
		}
		_, err := uc.CreateExpressOrder(context.Background(), o)
		if err == nil {
			t.Fatal("期望因保价超额返回错误，但得到 nil")
		}
		t.Logf("正确拦截保价超额: %v", err)
	})

	t.Run("nil入参报错误", func(t *testing.T) {
		uc := NewExpressOrderUsecase(newFakeOrderRepo(), newFakeMaterialRepo())
		if _, err := uc.CreateExpressOrder(context.Background(), nil); err == nil {
			t.Fatal("期望因 nil 入参返回错误，但得到 nil")
		}
	})
}
