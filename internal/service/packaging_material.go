package service

import (
	"context"

	v3 "shunfeng-miniprogram/api/order/v3"
	"shunfeng-miniprogram/internal/biz"

	"go.einride.tech/aip/pagination"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const defaultPackagingMaterialPageSize = 100

// PackagingMaterialService 包装材料服务。
type PackagingMaterialService struct {
	v3.UnimplementedPackagingMaterialServiceServer

	uc *biz.PackagingMaterialUsecase
}

// NewPackagingMaterialService 创建包装材料服务。
func NewPackagingMaterialService(uc *biz.PackagingMaterialUsecase) *PackagingMaterialService {
	return &PackagingMaterialService{uc: uc}
}

// ListPackagingMaterials 列出可选包装材料（前端下单时下拉选择）。
func (s *PackagingMaterialService) ListPackagingMaterials(ctx context.Context, req *v3.ListPackagingMaterialsRequest) (*v3.PackagingMaterialSet, error) {
	pageToken, err := pagination.ParsePageToken(req)
	if err != nil {
		return nil, err
	}
	if req.PageSize <= 0 {
		req.PageSize = defaultPackagingMaterialPageSize
	}
	materials, total, err := s.uc.ListPackagingMaterials(ctx,
		biz.ListLimit(int(req.GetPageSize())),
		biz.ListOffset(int(pageToken.Offset)),
	)
	if err != nil {
		return nil, err
	}
	set := &v3.PackagingMaterialSet{
		Materials:  make([]*v3.PackagingMaterial, 0, len(materials)),
		TotalCount: total,
	}
	if len(materials) >= int(req.GetPageSize()) {
		set.NextPageToken = pageToken.Next(req).String()
	}
	for _, m := range materials {
		set.Materials = append(set.Materials, convertPackagingMaterialReply(m))
	}
	return set, nil
}

// CreatePackagingMaterial 新增包装材料。
func (s *PackagingMaterialService) CreatePackagingMaterial(ctx context.Context, req *v3.CreatePackagingMaterialRequest) (*v3.PackagingMaterial, error) {
	m, err := s.uc.CreatePackagingMaterial(ctx, convertPackagingMaterial(req.GetMaterial()))
	if err != nil {
		return nil, err
	}
	return convertPackagingMaterialReply(m), nil
}

// UpdatePackagingMaterial 更新包装材料。
func (s *PackagingMaterialService) UpdatePackagingMaterial(ctx context.Context, req *v3.UpdatePackagingMaterialRequest) (*v3.PackagingMaterial, error) {
	current, err := s.uc.GetPackagingMaterial(ctx, req.GetMaterial().GetName())
	if err != nil {
		return nil, err
	}
	m := convertPackagingMaterial(req.GetMaterial())
	m.ID = current.ID
	updated, err := s.uc.UpdatePackagingMaterial(ctx, m, req.GetUpdateMask().GetPaths())
	if err != nil {
		return nil, err
	}
	return convertPackagingMaterialReply(updated), nil
}

// DeletePackagingMaterial 删除包装材料。
func (s *PackagingMaterialService) DeletePackagingMaterial(ctx context.Context, req *v3.DeletePackagingMaterialRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeletePackagingMaterial(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// convertPackagingMaterial 将 proto 转换为 biz 模型（DTO → DO）。
func convertPackagingMaterial(in *v3.PackagingMaterial) *biz.PackagingMaterial {
	if in == nil {
		return nil
	}
	return &biz.PackagingMaterial{
		ID:          in.GetId(),
		Name:        in.GetName(),
		UnitPrice:   in.GetUnitPrice(),
		Unit:        in.GetUnit(),
		Description: in.GetDescription(),
	}
}

// convertPackagingMaterialReply 将 biz 模型转换为 proto（DO → DTO）。
func convertPackagingMaterialReply(in *biz.PackagingMaterial) *v3.PackagingMaterial {
	if in == nil {
		return nil
	}
	return &v3.PackagingMaterial{
		Id:          in.ID,
		Name:        in.Name,
		UnitPrice:   in.UnitPrice,
		Unit:        in.Unit,
		Description: in.Description,
		CreatedAt:   timestamppb.New(in.CreatedAt),
		UpdatedAt:   timestamppb.New(in.UpdatedAt),
	}
}
