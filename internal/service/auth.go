package service

import (
	"context"
	"time"

	v1 "shunfeng-miniprogram/api/auth/v1"
	"shunfeng-miniprogram/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
)

// AuthService 是身份认证传输层，只负责 DTO 和 biz 命令之间的转换。
type AuthService struct {
	v1.UnimplementedAuthServiceServer

	uc *biz.AuthUsecase
}

// NewAuthService 创建身份认证服务。
func NewAuthService(uc *biz.AuthUsecase) *AuthService {
	return &AuthService{uc: uc}
}

// SendVerificationCode 发送注册等场景验证码。
func (s *AuthService) SendVerificationCode(ctx context.Context, req *v1.SendVerificationCodeRequest) (*v1.VerificationCodeReply, error) {
	vc, err := s.uc.SendVerificationCode(ctx, req.GetTarget(), req.GetTargetType(), req.GetScene())
	if err != nil {
		return nil, err
	}
	return &v1.VerificationCodeReply{
		VerificationId: vc.VerificationID,
		ExpiresIn:      int32(maxInt64(0, int64(vc.ExpiresAt.Sub(time.Now()).Seconds()))),
	}, nil
}

// Register 注册手机号或邮箱账号。
func (s *AuthService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.TokenReply, error) {
	result, err := s.uc.Register(ctx, &biz.RegisterCommand{
		Phone:      req.GetPhone(),
		Email:      req.GetEmail(),
		VerifyCode: req.GetVerifyCode(),
		Password:   req.GetPassword(),
		DeviceID:   req.GetDeviceId(),
	})
	if err != nil {
		return nil, err
	}
	return convertTokenReply(result), nil
}

// Login 使用账号密码登录并签发 Token。
func (s *AuthService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.TokenReply, error) {
	result, err := s.uc.Login(ctx, &biz.LoginCommand{
		Account:  req.GetAccount(),
		Password: req.GetPassword(),
		DeviceID: req.GetDeviceId(),
	})
	if err != nil {
		return nil, err
	}
	return convertTokenReply(result), nil
}

// RefreshToken 使用 refresh token 刷新登录态。
func (s *AuthService) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenReply, error) {
	result, err := s.uc.RefreshToken(ctx, req.GetRefreshToken(), req.GetDeviceId())
	if err != nil {
		return nil, err
	}
	return &v1.RefreshTokenReply{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// Logout 退出当前设备或全部设备。
func (s *AuthService) Logout(ctx context.Context, req *v1.LogoutRequest) (*emptypb.Empty, error) {
	if err := s.uc.Logout(ctx, req.GetUserId(), req.GetAccessToken(), req.GetRefreshToken(), req.GetDeviceId(), req.GetLogoutAll()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ChangePassword 修改密码，成功后清理全部登录会话。
func (s *AuthService) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*emptypb.Empty, error) {
	err := s.uc.ChangePassword(ctx, &biz.ChangePasswordCommand{
		UserID:      req.GetUserId(),
		OldPassword: req.GetOldPassword(),
		NewPassword: req.GetNewPassword(),
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// SubmitRealNameAuth 提交实名认证资料。
func (s *AuthService) SubmitRealNameAuth(ctx context.Context, req *v1.SubmitRealNameAuthRequest) (*v1.RealNameAuthReply, error) {
	auth, err := s.uc.SubmitRealNameAuth(ctx, &biz.RealNameAuthCommand{
		UserID:    req.GetUserId(),
		RealName:  req.GetRealName(),
		IDCardNo:  req.GetIdCardNo(),
		ImageURLs: req.GetImageUrls(),
	})
	if err != nil {
		return nil, err
	}
	return &v1.RealNameAuthReply{AuthId: auth.AuthID, Status: auth.AuthStatus}, nil
}

// GetRealNameAuth 查询实名认证状态，返回脱敏字段。
func (s *AuthService) GetRealNameAuth(ctx context.Context, req *v1.GetRealNameAuthRequest) (*v1.RealNameAuthInfo, error) {
	auth, err := s.uc.GetRealNameAuth(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &v1.RealNameAuthInfo{
		AuthId:       auth.AuthID,
		UserId:       auth.UserID,
		RealNameMask: auth.RealNameMask,
		IdCardMask:   auth.IDCardMask,
		ImageUrls:    auth.ImageURLs,
		Status:       auth.AuthStatus,
		RejectReason: auth.RejectReason,
	}, nil
}

func convertTokenReply(in *biz.TokenResult) *v1.TokenReply {
	return &v1.TokenReply{
		UserId:        in.UserID,
		AccessToken:   in.AccessToken,
		RefreshToken:  in.RefreshToken,
		Role:          in.RoleCode,
		AccountStatus: in.AccountStatus,
		ExpiresIn:     in.ExpiresIn,
	}
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
