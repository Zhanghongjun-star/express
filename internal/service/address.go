package service

import (
	"context"

	v1 "shunfeng-miniprogram/api/address/v1"
	"shunfeng-miniprogram/internal/biz"

	"go.einride.tech/aip/pagination"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const defaultAddressPageSize = 20

// AddressService 地址服务。
type AddressService struct {
	v1.UnimplementedAddressServiceServer

	uc *biz.AddressUsecase
}

// NewAddressService 创建地址服务。
func NewAddressService(uc *biz.AddressUsecase) *AddressService {
	return &AddressService{uc: uc}
}

// ListAddresses 获取地址列表。
func (s *AddressService) ListAddresses(ctx context.Context, req *v1.ListAddressesRequest) (*v1.AddressSet, error) {
	pageToken, err := pagination.ParsePageToken(req)
	if err != nil {
		return nil, err
	}
	if req.PageSize <= 0 {
		req.PageSize = defaultAddressPageSize
	}
	addresses, err := s.uc.ListAddresses(ctx, req.GetUserId(),
		biz.AddressListType(int32(req.GetAddrType())),
		biz.AddressListLimit(int(req.GetPageSize())),
		biz.AddressListOffset(int(pageToken.Offset)),
	)
	if err != nil {
		return nil, err
	}
	set := &v1.AddressSet{
		Addresses: make([]*v1.Address, 0, len(addresses)),
	}
	if len(addresses) >= int(req.PageSize) {
		set.NextPageToken = pageToken.Next(req).String()
	}
	for _, address := range addresses {
		set.Addresses = append(set.Addresses, convertAddressReply(address))
	}
	return set, nil
}

// CreateAddress 创建地址。
func (s *AddressService) CreateAddress(ctx context.Context, req *v1.CreateAddressRequest) (*v1.Address, error) {
	address, err := s.uc.CreateAddress(ctx, convertAddress(req.GetAddress()))
	if err != nil {
		return nil, err
	}
	return convertAddressReply(address), nil
}

// UpdateAddress 更新地址。
func (s *AddressService) UpdateAddress(ctx context.Context, req *v1.UpdateAddressRequest) (*v1.Address, error) {
	current, err := s.uc.GetAddress(ctx, req.GetAddress().GetUserId(), req.GetAddress().GetId())
	if err != nil {
		return nil, err
	}
	applyAddressMask(current, req.GetAddress(), req.GetUpdateMask().GetPaths())
	address, err := s.uc.UpdateAddress(ctx, current)
	if err != nil {
		return nil, err
	}
	return convertAddressReply(address), nil
}

// DeleteAddress 删除地址。
func (s *AddressService) DeleteAddress(ctx context.Context, req *v1.DeleteAddressRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeleteAddress(ctx, req.GetUserId(), req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// BatchDeleteAddresses 批量删除地址。
func (s *AddressService) BatchDeleteAddresses(ctx context.Context, req *v1.BatchDeleteAddressesRequest) (*v1.BatchDeleteAddressesReply, error) {
	deleted, err := s.uc.BatchDeleteAddresses(ctx, req.GetUserId(), req.GetIds())
	if err != nil {
		return nil, err
	}
	return &v1.BatchDeleteAddressesReply{DeletedCount: deleted}, nil
}

// ParseAddress 解析文本为地址。
func (s *AddressService) ParseAddress(ctx context.Context, req *v1.ParseAddressRequest) (*v1.ParseAddressReply, error) {
	address, err := s.uc.ParseAddress(ctx, req.GetContent())
	if err != nil {
		return nil, err
	}
	return &v1.ParseAddressReply{Address: convertAddressReply(address)}, nil
}

// convertAddress 将 v1.Address 转换为 biz.Address（DTO → DO）。
func convertAddress(in *v1.Address) *biz.Address {
	if in == nil {
		return nil
	}
	return &biz.Address{
		ID:            in.GetId(),
		UserID:        in.GetUserId(),
		AddrType:      int32(in.GetAddrType()),
		ReceiverName:  in.GetReceiverName(),
		ReceiverPhone: in.GetReceiverPhone(),
		Province:      in.GetProvince(),
		City:          in.GetCity(),
		District:      in.GetDistrict(),
		DetailAddr:    in.GetDetailAddr(),
		IsDefault:     in.GetIsDefault(),
		DelFlag:       in.GetDelFlag(),
	}
}

// convertAddressReply 将 biz.Address 转换为 v1.Address（DO → DTO）。
func convertAddressReply(in *biz.Address) *v1.Address {
	if in == nil {
		return nil
	}
	return &v1.Address{
		Id:            in.ID,
		UserId:        in.UserID,
		AddrType:      v1.AddressType(in.AddrType),
		ReceiverName:  in.ReceiverName,
		ReceiverPhone: in.ReceiverPhone,
		Province:      in.Province,
		City:          in.City,
		District:      in.District,
		DetailAddr:    in.DetailAddr,
		IsDefault:     in.IsDefault,
		CreateTime:    timestamppb.New(in.CreateTime),
		UpdateTime:    timestamppb.New(in.UpdateTime),
		DelFlag:       in.DelFlag,
	}
}

// applyAddressMask 根据字段掩码更新地址字段。
func applyAddressMask(dst *biz.Address, src *v1.Address, paths []string) {
	if dst == nil || src == nil {
		return
	}
	for _, path := range paths {
		switch path {
		case "*":
			*dst = *convertAddress(src)
			return
		case "addr_type":
			dst.AddrType = int32(src.GetAddrType())
		case "receiver_name":
			dst.ReceiverName = src.GetReceiverName()
		case "receiver_phone":
			dst.ReceiverPhone = src.GetReceiverPhone()
		case "province":
			dst.Province = src.GetProvince()
		case "city":
			dst.City = src.GetCity()
		case "district":
			dst.District = src.GetDistrict()
		case "detail_addr":
			dst.DetailAddr = src.GetDetailAddr()
		case "is_default":
			dst.IsDefault = src.GetIsDefault()
		}
	}
}
