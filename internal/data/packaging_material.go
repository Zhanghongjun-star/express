package data

import (
	"context"
	"time"

	"shunfeng-miniprogram/internal/biz"
)

// packagingMaterialRepo 包装材料仓储的 GORM 实现。
type packagingMaterialRepo struct{}

// NewPackagingMaterialRepo 创建包装材料仓储（返回 biz 接口）。
func NewPackagingMaterialRepo() biz.PackagingMaterialRepo {
	return &packagingMaterialRepo{}
}

// Create 新增包装材料。
func (r *packagingMaterialRepo) Create(ctx context.Context, m *biz.PackagingMaterial) (*biz.PackagingMaterial, error) {
	po := newPackagingMaterialPO(m)
	po.CreatedAt = time.Now()
	po.UpdatedAt = time.Now()
	if err := DB.WithContext(ctx).Create(po).Error; err != nil {
		return nil, err
	}
	return toBizPackagingMaterial(po), nil
}

// Update 更新包装材料（全量更新）。
func (r *packagingMaterialRepo) Update(ctx context.Context, m *biz.PackagingMaterial, mask []string) (*biz.PackagingMaterial, error) {
	var po PackagingMaterial
	if err := DB.WithContext(ctx).First(&po, m.ID).Error; err != nil {
		return nil, biz.ErrPackagingMaterialNotFound
	}
	updated := newPackagingMaterialPO(m)
	updated.ID = po.ID
	updated.CreatedAt = po.CreatedAt
	updated.UpdatedAt = time.Now()
	if err := DB.WithContext(ctx).Model(&po).Updates(updated).Error; err != nil {
		return nil, err
	}
	return toBizPackagingMaterial(updated), nil
}

// Delete 删除包装材料。
func (r *packagingMaterialRepo) Delete(ctx context.Context, id int64) error {
	result := DB.WithContext(ctx).Delete(&PackagingMaterial{}, id)
	if result.RowsAffected == 0 {
		return biz.ErrPackagingMaterialNotFound
	}
	return result.Error
}

// GetByName 按名称查询包装材料。
func (r *packagingMaterialRepo) GetByName(ctx context.Context, name string) (*biz.PackagingMaterial, error) {
	var po PackagingMaterial
	if err := DB.WithContext(ctx).Where("name = ?", name).First(&po).Error; err != nil {
		return nil, biz.ErrPackagingMaterialNotFound
	}
	return toBizPackagingMaterial(&po), nil
}

// GetByID 按主键 ID 查询包装材料（下单时按所选材料 id 解析单价）。
func (r *packagingMaterialRepo) GetByID(ctx context.Context, id int64) (*biz.PackagingMaterial, error) {
	var po PackagingMaterial
	if err := DB.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, biz.ErrPackagingMaterialNotFound
	}
	return toBizPackagingMaterial(&po), nil
}

// List 列出包装材料，支持 AIP 分页选项。
func (r *packagingMaterialRepo) List(ctx context.Context, opts ...biz.ListOption) ([]*biz.PackagingMaterial, int32, error) {
	options := biz.ListOptions{Limit: 100}
	for _, opt := range opts {
		opt(&options)
	}
	var total int64
	if err := DB.WithContext(ctx).Model(&PackagingMaterial{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	tx := DB.WithContext(ctx).Order("id ASC")
	if options.Limit > 0 {
		tx = tx.Limit(options.Limit)
	}
	if options.Offset > 0 {
		tx = tx.Offset(options.Offset)
	}
	var pos []PackagingMaterial
	if err := tx.Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*biz.PackagingMaterial, len(pos))
	for i := range pos {
		out[i] = toBizPackagingMaterial(&pos[i])
	}
	return out, int32(total), nil
}

// newPackagingMaterialPO 将 biz 模型转换为 PO（DO → PO）。
func newPackagingMaterialPO(in *biz.PackagingMaterial) *PackagingMaterial {
	if in == nil {
		return nil
	}
	return &PackagingMaterial{
		ID:          in.ID,
		Name:        in.Name,
		UnitPrice:   in.UnitPrice,
		Unit:        in.Unit,
		Description: in.Description,
	}
}

// toBizPackagingMaterial 将 PO 转换为 biz 模型（PO → DO）。
func toBizPackagingMaterial(in *PackagingMaterial) *biz.PackagingMaterial {
	if in == nil {
		return nil
	}
	return &biz.PackagingMaterial{
		ID:          in.ID,
		Name:        in.Name,
		UnitPrice:   in.UnitPrice,
		Unit:        in.Unit,
		Description: in.Description,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

// seedPackagingMaterials 若包装材料表为空，则写入默认示例材料（幂等）。
func seedPackagingMaterials(ctx context.Context) {
	var count int64
	if err := DB.WithContext(ctx).Model(&PackagingMaterial{}).Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		return
	}
	now := time.Now()
	for _, m := range biz.DefaultPackagingMaterials() {
		po := newPackagingMaterialPO(m)
		po.CreatedAt = now
		po.UpdatedAt = now
		_ = DB.WithContext(ctx).Create(po).Error
	}
}
