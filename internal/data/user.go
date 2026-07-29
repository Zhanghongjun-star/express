package data

import (
	"context"
	"time"

	"shunfeng-miniprogram/internal/biz"
)

// userRepo 用户仓储的 GORM 实现。
type userRepo struct{}

// NewUserRepo 创建用户仓储。
func NewUserRepo() biz.UserRepo {
	return &userRepo{}
}

func (r *userRepo) FindByID(ctx context.Context, id int64) (*biz.User, error) {
	var po User
	if err := DB.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, biz.ErrUserNotFound
	}
	return toUserBiz(&po), nil
}

func (r *userRepo) UpdateAvatar(ctx context.Context, id int64, avatarURL string) (*biz.User, error) {
	var po User
	if err := DB.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, biz.ErrUserNotFound
	}
	now := time.Now()
	if err := DB.WithContext(ctx).Model(&po).Updates(map[string]any{
		"avatar_url":  avatarURL,
		"update_time": now,
	}).Error; err != nil {
		return nil, err
	}
	po.AvatarURL = avatarURL
	po.UpdateTime = now
	return toUserBiz(&po), nil
}

func (r *userRepo) UpdateNickname(ctx context.Context, id int64, nickName string) (*biz.User, error) {
	var po User
	if err := DB.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, biz.ErrUserNotFound
	}
	now := time.Now()
	if err := DB.WithContext(ctx).Model(&po).Updates(map[string]any{
		"nick_name":   nickName,
		"update_time": now,
	}).Error; err != nil {
		return nil, err
	}
	po.NickName = nickName
	po.UpdateTime = now
	return toUserBiz(&po), nil
}

func (r *userRepo) GetAuthStatus(ctx context.Context, userID int64) (*biz.AuthStatus, error) {
	var po User
	if err := DB.WithContext(ctx).First(&po, userID).Error; err != nil {
		return nil, biz.ErrUserNotFound
	}
	status := &biz.AuthStatus{
		UserID:       po.ID,
		RealNameAuth: po.RealNameAuth,
	}
	var auth IdentityRealNameAuth
	if err := DB.WithContext(ctx).Where("user_id = ?", userID).First(&auth).Error; err == nil {
		status.RejectReason = auth.RejectReason
	}
	return status, nil
}

func (r *userRepo) GetAccountStatus(ctx context.Context, userID int64) (*biz.AccountStatusInfo, error) {
	var po User
	if err := DB.WithContext(ctx).First(&po, userID).Error; err != nil {
		return nil, biz.ErrUserNotFound
	}
	info := &biz.AccountStatusInfo{
		UserID:        po.ID,
		AccountStatus: po.AccountStatus,
		RealNameAuth:  po.RealNameAuth,
		IsEnterprise:  po.IsEnterprise,
	}
	switch {
	case po.AccountStatus == 2:
		info.CanModifyProfile = false
		info.CanManageAddress = false
		info.CanQueryHistory = true
		info.CanExport = false
	case po.RealNameAuth == 0 || po.RealNameAuth == 3:
		info.CanModifyProfile = true
		info.CanManageAddress = true
		info.CanQueryHistory = true
		info.CanExport = false
	default:
		info.CanModifyProfile = true
		info.CanManageAddress = true
		info.CanQueryHistory = true
		info.CanExport = true
	}
	return info, nil
}

// toUserBiz 将 User 转换为 biz.User（PO → DO）。
func toUserBiz(in *User) *biz.User {
	if in == nil {
		return nil
	}
	return &biz.User{
		ID:              in.ID,
		Phone:           in.Phone,
		Email:           in.Email,
		AvatarURL:       in.AvatarURL,
		NickName:        in.NickName,
		RealNameAuth:    in.RealNameAuth,
		AccountStatus:   in.AccountStatus,
		IsEnterprise:    in.IsEnterprise,
		QueryCountToday: in.QueryCountToday,
		CreateTime:      in.CreateTime,
		UpdateTime:      in.UpdateTime,
	}
}
