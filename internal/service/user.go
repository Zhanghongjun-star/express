package service

import (
	"context"

	v1 "shunfeng-miniprogram/api/user/v1"
	"shunfeng-miniprogram/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserService 用户服务。
type UserService struct {
	v1.UnimplementedUserServiceServer

	uc *biz.UserUsecase
}

// NewUserService 创建用户服务。
func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

// GetProfile 获取个人资料。
func (s *UserService) GetProfile(ctx context.Context, req *v1.GetProfileRequest) (*v1.Profile, error) {
	user, err := s.uc.GetProfile(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return convertProfileReply(user), nil
}

// UpdateAvatar 修改头像。
func (s *UserService) UpdateAvatar(ctx context.Context, req *v1.UpdateAvatarRequest) (*v1.Profile, error) {
	user, err := s.uc.UpdateAvatar(ctx, req.GetUserId(), req.GetAvatarUrl())
	if err != nil {
		return nil, err
	}
	return convertProfileReply(user), nil
}

// UpdateNickname 修改昵称。
func (s *UserService) UpdateNickname(ctx context.Context, req *v1.UpdateNicknameRequest) (*v1.Profile, error) {
	user, err := s.uc.UpdateNickname(ctx, req.GetUserId(), req.GetNickName())
	if err != nil {
		return nil, err
	}
	return convertProfileReply(user), nil
}

// GetAccountStatus 查询账号状态和权限信息。
func (s *UserService) GetAccountStatus(ctx context.Context, req *v1.GetAccountStatusRequest) (*v1.AccountStatusInfo, error) {
	info, err := s.uc.GetAccountStatus(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &v1.AccountStatusInfo{
		UserId:           info.UserID,
		AccountStatus:    info.AccountStatus,
		RealNameAuth:     info.RealNameAuth,
		IsEnterprise:     info.IsEnterprise,
		CanModifyProfile: info.CanModifyProfile,
		CanManageAddress: info.CanManageAddress,
		CanQueryHistory:  info.CanQueryHistory,
		CanExport:        info.CanExport,
	}, nil
}

// GetAuthStatus 查询实名认证状态。
func (s *UserService) GetAuthStatus(ctx context.Context, req *v1.GetAuthStatusRequest) (*v1.AuthStatus, error) {
	status, err := s.uc.GetAuthStatus(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &v1.AuthStatus{
		UserId:       status.UserID,
		RealNameAuth: status.RealNameAuth,
		RejectReason: status.RejectReason,
	}, nil
}

// convertProfileReply 将 biz.User 转换为 v1.Profile。
func convertProfileReply(in *biz.User) *v1.Profile {
	if in == nil {
		return nil
	}
	return &v1.Profile{
		UserId:          in.ID,
		AvatarUrl:       in.AvatarURL,
		NickName:        in.NickName,
		Phone:           in.Phone,
		RealNameAuth:    in.RealNameAuth,
		AccountStatus:   in.AccountStatus,
		IsEnterprise:    in.IsEnterprise,
		QueryCountToday: in.QueryCountToday,
		CreateTime:      timestamppb.New(in.CreateTime),
		UpdateTime:      timestamppb.New(in.UpdateTime),
	}
}
