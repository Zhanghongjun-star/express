package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v3/errors"
)

// PackagingMaterial 包装材料领域对象（打包费计费数据源）。
type PackagingMaterial struct {
	ID          int64
	Name        string  // 材料名称（唯一）
	UnitPrice   float64 // 单价（元）
	Unit        string  // 计价单位（个/卷/张/件）
	Description string  // 说明
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PackagingMaterialRepo 包装材料仓储接口（数据层实现）。
type PackagingMaterialRepo interface {
	Create(ctx context.Context, m *PackagingMaterial) (*PackagingMaterial, error)
	Update(ctx context.Context, m *PackagingMaterial, mask []string) (*PackagingMaterial, error)
	Delete(ctx context.Context, id int64) error
	GetByName(ctx context.Context, name string) (*PackagingMaterial, error)
	GetByID(ctx context.Context, id int64) (*PackagingMaterial, error)
	List(ctx context.Context, opts ...ListOption) ([]*PackagingMaterial, int32, error)
}

// 包装材料相关错误。
var (
	ErrPackagingMaterialNotFound        = errors.NotFound("PACKAGING_MATERIAL_NOT_FOUND", "包装材料不存在")
	ErrPackagingMaterialInvalidArgument = errors.BadRequest("PACKAGING_MATERIAL_INVALID_ARGUMENT", "包装材料参数不合法")
)

// PackagingMaterialUsecase 包装材料用例。
type PackagingMaterialUsecase struct {
	repo PackagingMaterialRepo
}

// NewPackagingMaterialUsecase 创建包装材料用例。
func NewPackagingMaterialUsecase(repo PackagingMaterialRepo) *PackagingMaterialUsecase {
	return &PackagingMaterialUsecase{repo: repo}
}

// CreatePackagingMaterial 新增包装材料。
func (uc *PackagingMaterialUsecase) CreatePackagingMaterial(ctx context.Context, m *PackagingMaterial) (*PackagingMaterial, error) {
	if m == nil || m.Name == "" || m.UnitPrice < 0 {
		return nil, ErrPackagingMaterialInvalidArgument
	}
	return uc.repo.Create(ctx, m)
}

// UpdatePackagingMaterial 更新包装材料。
func (uc *PackagingMaterialUsecase) UpdatePackagingMaterial(ctx context.Context, m *PackagingMaterial, mask []string) (*PackagingMaterial, error) {
	if m == nil || m.ID == 0 {
		return nil, ErrPackagingMaterialInvalidArgument
	}
	return uc.repo.Update(ctx, m, mask)
}

// DeletePackagingMaterial 删除包装材料。
func (uc *PackagingMaterialUsecase) DeletePackagingMaterial(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrPackagingMaterialInvalidArgument
	}
	return uc.repo.Delete(ctx, id)
}

// GetPackagingMaterial 按名称查询包装材料。
func (uc *PackagingMaterialUsecase) GetPackagingMaterial(ctx context.Context, name string) (*PackagingMaterial, error) {
	return uc.repo.GetByName(ctx, name)
}

// ListPackagingMaterials 列出包装材料（供前端下单时选择）。
func (uc *PackagingMaterialUsecase) ListPackagingMaterials(ctx context.Context, opts ...ListOption) ([]*PackagingMaterial, int32, error) {
	return uc.repo.List(ctx, opts...)
}

// DefaultPackagingMaterials 默认示例材料（用于初始化 packaging_material 表）。
// 注意：单价为常见顺丰包装示例价，请按官方最新标准在校验后调整。
func DefaultPackagingMaterials() []*PackagingMaterial {
	return []*PackagingMaterial{
		{Name: "文件封", UnitPrice: 1, Unit: "个", Description: "顺丰文件封"},
		{Name: "1号纸箱", UnitPrice: 2, Unit: "个", Description: "小号纸箱"},
		{Name: "2号纸箱", UnitPrice: 3, Unit: "个", Description: "中号纸箱"},
		{Name: "3号纸箱", UnitPrice: 4, Unit: "个", Description: "大号纸箱"},
		{Name: "气泡膜", UnitPrice: 5, Unit: "卷", Description: "气泡膜/卷"},
		{Name: "珍珠棉", UnitPrice: 6, Unit: "张", Description: "珍珠棉/张"},
		{Name: "防水袋", UnitPrice: 1, Unit: "个", Description: "防水袋"},
		{Name: "缠绕膜", UnitPrice: 4, Unit: "卷", Description: "缠绕膜/卷"},
		{Name: "木架", UnitPrice: 20, Unit: "件", Description: "木架加固"},
		{Name: "木箱", UnitPrice: 30, Unit: "件", Description: "木箱"},
	}
}
