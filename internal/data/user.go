package data

import (
	"context"
	"sync"
	"time"

	"shunfeng-miniprogram/internal/biz"
)

// userRepo 用户仓储的内存实现。
type userRepo struct {
	data *Data

	mu         sync.RWMutex
	nextID     int64
	users      map[int64]*biz.User
	authStatus map[int64]*biz.AuthStatus
}

// NewUserRepo 创建用户仓储。
func NewUserRepo(data *Data) biz.UserRepo {
	r := &userRepo{
		data:       data,
		nextID:     1,
		users:      make(map[int64]*biz.User),
		authStatus: make(map[int64]*biz.AuthStatus),
	}
	now := time.Now()
	r.users[1] = &biz.User{
		ID:              1,
		Phone:           "138****1234",
		Email:           "user@example.com",
		AvatarURL:       "https://example.com/avatar/default.png",
		NickName:        "用户",
		RealNameAuth:    0,
		AccountStatus:   1,
		IsEnterprise:    false,
		QueryCountToday: 0,
		CreateTime:      now,
		UpdateTime:      now,
	}
	r.authStatus[1] = &biz.AuthStatus{
		UserID:       1,
		RealNameAuth: 0,
		RejectReason: "",
	}
	return r
}

func (r *userRepo) FindByID(_ context.Context, id int64) (*biz.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, biz.ErrUserNotFound
	}
	return cloneUser(user), nil
}

func (r *userRepo) UpdateAvatar(_ context.Context, id int64, avatarURL string) (*biz.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, biz.ErrUserNotFound
	}
	updated := cloneUser(user)
	updated.AvatarURL = avatarURL
	updated.UpdateTime = time.Now()
	r.users[id] = cloneUser(updated)
	return cloneUser(updated), nil
}

func (r *userRepo) UpdateNickname(_ context.Context, id int64, nickName string) (*biz.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, biz.ErrUserNotFound
	}
	updated := cloneUser(user)
	updated.NickName = nickName
	updated.UpdateTime = time.Now()
	r.users[id] = cloneUser(updated)
	return cloneUser(updated), nil
}

func (r *userRepo) GetAuthStatus(_ context.Context, userID int64) (*biz.AuthStatus, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	status, ok := r.authStatus[userID]
	if !ok {
		return nil, biz.ErrUserNotFound
	}
	cp := *status
	return &cp, nil
}

func (r *userRepo) GetAccountStatus(_ context.Context, userID int64) (*biz.AccountStatusInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[userID]
	if !ok {
		return nil, biz.ErrUserNotFound
	}
	info := &biz.AccountStatusInfo{
		UserID:        user.ID,
		AccountStatus: user.AccountStatus,
		RealNameAuth:  user.RealNameAuth,
		IsEnterprise:  user.IsEnterprise,
	}
	switch {
	case user.AccountStatus == 2:
		info.CanModifyProfile = false
		info.CanManageAddress = false
		info.CanQueryHistory = true
		info.CanExport = false
	case user.RealNameAuth == 0 || user.RealNameAuth == 3:
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

// cloneUser 深拷贝用户对象，防止数据竞态。
func cloneUser(u *biz.User) *biz.User {
	if u == nil {
		return nil
	}
	return &biz.User{
		ID:              u.ID,
		Phone:           u.Phone,
		Email:           u.Email,
		AvatarURL:       u.AvatarURL,
		NickName:        u.NickName,
		RealNameAuth:    u.RealNameAuth,
		AccountStatus:   u.AccountStatus,
		IsEnterprise:    u.IsEnterprise,
		QueryCountToday: u.QueryCountToday,
		CreateTime:      u.CreateTime,
		UpdateTime:      u.UpdateTime,
	}
}
