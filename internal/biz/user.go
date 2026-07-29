package biz

import (
	"context"
	"strings"
	"time"

	v1 "shunfeng-miniprogram/api/user/v1"

	"github.com/go-kratos/kratos/v3/errors"
)

var (
	// ErrUserNotFound 用户不存在时返回。
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
	// ErrUserInvalidArgument 请求参数无效时返回。
	ErrUserInvalidArgument = errors.BadRequest(v1.ErrorReason_USER_INVALID_ARGUMENT.String(), "invalid user argument")
	// ErrUserUnauthorized 用户未登录时返回。
	ErrUserUnauthorized = errors.Unauthorized(v1.ErrorReason_USER_UNAUTHORIZED.String(), "user unauthorized")
	// ErrUserForbidden 账号被封禁时返回。
	ErrUserForbidden = errors.Forbidden(v1.ErrorReason_USER_FORBIDDEN.String(), "user account is banned")
)

// User 用户个人资料领域对象。
type User struct {
	ID              int64
	Phone           string
	Email           string
	AvatarURL       string
	NickName        string
	RealNameAuth    int32
	AccountStatus   int32
	IsEnterprise    bool
	QueryCountToday int32
	CreateTime      time.Time
	UpdateTime      time.Time
}

// AuthStatus 实名认证状态领域对象。
type AuthStatus struct {
	UserID       int64
	RealNameAuth int32
	RejectReason string
}

// AccountStatusInfo 账号状态与权限领域对象。
type AccountStatusInfo struct {
	UserID          int64
	AccountStatus   int32
	RealNameAuth    int32
	IsEnterprise    bool
	CanModifyProfile bool
	CanManageAddress bool
	CanQueryHistory  bool
	CanExport        bool
}

// UserRepo 用户仓储接口。
type UserRepo interface {
	FindByID(context.Context, int64) (*User, error)
	UpdateAvatar(context.Context, int64, string) (*User, error)
	UpdateNickname(context.Context, int64, string) (*User, error)
	GetAuthStatus(context.Context, int64) (*AuthStatus, error)
	GetAccountStatus(context.Context, int64) (*AccountStatusInfo, error)
}

// UserUsecase 用户用例。
type UserUsecase struct {
	repo UserRepo
}

// NewUserUsecase 创建用户用例。
func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// GetProfile 获取用户个人资料。
func (uc *UserUsecase) GetProfile(ctx context.Context, userID int64) (*User, error) {
	if userID <= 0 {
		return nil, ErrUserInvalidArgument
	}
	return uc.repo.FindByID(ctx, userID)
}

// UpdateAvatar 修改用户头像。
func (uc *UserUsecase) UpdateAvatar(ctx context.Context, userID int64, avatarURL string) (*User, error) {
	if userID <= 0 || strings.TrimSpace(avatarURL) == "" {
		return nil, ErrUserInvalidArgument
	}
	user, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.AccountStatus == 2 {
		return nil, ErrUserForbidden
	}
	return uc.repo.UpdateAvatar(ctx, userID, avatarURL)
}

// UpdateNickname 修改用户昵称。
func (uc *UserUsecase) UpdateNickname(ctx context.Context, userID int64, nickName string) (*User, error) {
	if userID <= 0 || strings.TrimSpace(nickName) == "" {
		return nil, ErrUserInvalidArgument
	}
	if len([]rune(nickName)) > 32 {
		return nil, ErrUserInvalidArgument
	}
	user, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.AccountStatus == 2 {
		return nil, ErrUserForbidden
	}
	return uc.repo.UpdateNickname(ctx, userID, nickName)
}

// GetAuthStatus 查询实名认证状态。
func (uc *UserUsecase) GetAuthStatus(ctx context.Context, userID int64) (*AuthStatus, error) {
	if userID <= 0 {
		return nil, ErrUserInvalidArgument
	}
	return uc.repo.GetAuthStatus(ctx, userID)
}

// GetAccountStatus 查询账号状态和权限信息。
func (uc *UserUsecase) GetAccountStatus(ctx context.Context, userID int64) (*AccountStatusInfo, error) {
	if userID <= 0 {
		return nil, ErrUserInvalidArgument
	}
	return uc.repo.GetAccountStatus(ctx, userID)
}
