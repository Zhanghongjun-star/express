package service

import (
	"context"
	"testing"
	"time"

	v1 "shunfeng-miniprogram/api/auth/v1"
	"shunfeng-miniprogram/internal/biz"

	"golang.org/x/crypto/bcrypt"
)

func TestAuthServiceRegisterLoginAndRealNameAuth(t *testing.T) {
	repo := newServiceAuthRepo(t)
	svc := NewAuthService(biz.NewAuthUsecase(repo))
	ctx := context.Background()

	registered, err := svc.Register(ctx, &v1.RegisterRequest{
		Phone:      "13800138000",
		VerifyCode: "123456",
		Password:   "password123",
		DeviceId:   "ios-1",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if registered.GetUserId() == 0 || registered.GetAccessToken() == "" || registered.GetRefreshToken() == "" {
		t.Fatalf("Register() = %+v, want user id and tokens", registered)
	}

	loggedIn, err := svc.Login(ctx, &v1.LoginRequest{
		Account:  "13800138000",
		Password: "password123",
		DeviceId: "ios-1",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if loggedIn.GetRole() != biz.RoleUser || loggedIn.GetAccountStatus() != biz.AccountStatusNormal {
		t.Fatalf("Login() = %+v, want role/account status", loggedIn)
	}

	submitted, err := svc.SubmitRealNameAuth(ctx, &v1.SubmitRealNameAuthRequest{
		UserId:    registered.GetUserId(),
		RealName:  "张三",
		IdCardNo:  "11010119900307001X",
		ImageUrls: []string{"https://qiniu.example/id-front.jpg"},
	})
	if err != nil {
		t.Fatalf("SubmitRealNameAuth() error = %v", err)
	}
	if submitted.GetStatus() != biz.RealNameStatusPending {
		t.Fatalf("SubmitRealNameAuth() status = %q, want pending", submitted.GetStatus())
	}

	info, err := svc.GetRealNameAuth(ctx, &v1.GetRealNameAuthRequest{UserId: registered.GetUserId()})
	if err != nil {
		t.Fatalf("GetRealNameAuth() error = %v", err)
	}
	if info.GetRealNameMask() != "张*" || info.GetIdCardMask() != "1101**********001X" {
		t.Fatalf("GetRealNameAuth() = %+v, want masked sensitive fields", info)
	}
}

type serviceAuthRepo struct {
	t              *testing.T
	nextID         int64
	usersByID      map[int64]*biz.IdentityAccount
	usersByAccount map[string]*biz.IdentityAccount
	verification   *biz.VerificationCode
	sessions       map[string]*biz.TokenSession
	realNameAuth   map[int64]*biz.RealNameAuth
}

func newServiceAuthRepo(t *testing.T) *serviceAuthRepo {
	return &serviceAuthRepo{
		t:              t,
		nextID:         1,
		usersByID:      make(map[int64]*biz.IdentityAccount),
		usersByAccount: make(map[string]*biz.IdentityAccount),
		verification: &biz.VerificationCode{
			Target:     "13800138000",
			TargetType: biz.TargetTypePhone,
			Scene:      biz.VerifySceneRegister,
			Code:       "123456",
			ExpiresAt:  time.Now().Add(time.Minute),
		},
		sessions:     make(map[string]*biz.TokenSession),
		realNameAuth: make(map[int64]*biz.RealNameAuth),
	}
}

func (r *serviceAuthRepo) SaveVerificationCode(_ context.Context, vc *biz.VerificationCode) error {
	r.verification = vc
	return nil
}

func (r *serviceAuthRepo) GetVerificationCode(_ context.Context, target, targetType, scene string) (*biz.VerificationCode, error) {
	if r.verification == nil || r.verification.Target != target || r.verification.TargetType != targetType || r.verification.Scene != scene {
		return nil, biz.ErrAuthInvalidArgument
	}
	cp := *r.verification
	return &cp, nil
}

func (r *serviceAuthRepo) FindUserByAccount(_ context.Context, account string) (*biz.IdentityAccount, error) {
	if user, ok := r.usersByAccount[account]; ok {
		return cloneServiceAccount(user), nil
	}
	return nil, biz.ErrAuthNotFound
}

func (r *serviceAuthRepo) FindUserByID(_ context.Context, userID int64) (*biz.IdentityAccount, error) {
	if user, ok := r.usersByID[userID]; ok {
		return cloneServiceAccount(user), nil
	}
	return nil, biz.ErrAuthNotFound
}

func (r *serviceAuthRepo) CreateUser(_ context.Context, user *biz.IdentityAccount) (*biz.IdentityAccount, error) {
	cp := cloneServiceAccount(user)
	cp.ID = r.nextID
	r.nextID++
	r.usersByID[cp.ID] = cloneServiceAccount(cp)
	if cp.Phone != "" {
		r.usersByAccount[cp.Phone] = cloneServiceAccount(cp)
	}
	if cp.Email != "" {
		r.usersByAccount[cp.Email] = cloneServiceAccount(cp)
	}
	return cloneServiceAccount(cp), nil
}

func (r *serviceAuthRepo) RecordLoginFailure(context.Context, int64, time.Duration) (int, error) {
	return 1, nil
}

func (r *serviceAuthRepo) ClearLoginFailures(context.Context, int64) error {
	return nil
}

func (r *serviceAuthRepo) UpdateLoginLock(context.Context, int64, *time.Time) error {
	return nil
}

func (r *serviceAuthRepo) SaveTokenSession(_ context.Context, session *biz.TokenSession) error {
	r.sessions[session.RefreshToken] = session
	return nil
}

func (r *serviceAuthRepo) GetTokenSession(_ context.Context, refreshToken string) (*biz.TokenSession, error) {
	if session, ok := r.sessions[refreshToken]; ok {
		cp := *session
		return &cp, nil
	}
	return nil, biz.ErrAuthUnauthenticated
}

func (r *serviceAuthRepo) DeleteTokenSession(_ context.Context, refreshToken string) error {
	delete(r.sessions, refreshToken)
	return nil
}

func (r *serviceAuthRepo) DeleteUserTokenSessions(_ context.Context, userID int64) error {
	for token, session := range r.sessions {
		if session.UserID == userID {
			delete(r.sessions, token)
		}
	}
	return nil
}

func (r *serviceAuthRepo) BlacklistAccessToken(context.Context, string, time.Duration) error {
	return nil
}

func (r *serviceAuthRepo) ValidateAccessToken(_ context.Context, accessToken string) (*biz.TokenSession, error) {
	for _, session := range r.sessions {
		if session.AccessToken == accessToken {
			cp := *session
			return &cp, nil
		}
	}
	return nil, biz.ErrAuthUnauthenticated
}

func (r *serviceAuthRepo) UpdatePassword(_ context.Context, userID int64, passwordHash string) error {
	if user, ok := r.usersByID[userID]; ok {
		user.PasswordHash = passwordHash
	}
	return nil
}

func (r *serviceAuthRepo) SaveRealNameAuth(_ context.Context, auth *biz.RealNameAuth) (*biz.RealNameAuth, error) {
	cp := *auth
	cp.AuthID = int64(len(r.realNameAuth) + 1)
	r.realNameAuth[auth.UserID] = &cp
	return &cp, nil
}

func (r *serviceAuthRepo) GetRealNameAuth(_ context.Context, userID int64) (*biz.RealNameAuth, error) {
	if auth, ok := r.realNameAuth[userID]; ok {
		cp := *auth
		return &cp, nil
	}
	return nil, biz.ErrAuthNotFound
}

func (r *serviceAuthRepo) IDCardHashExists(context.Context, string, int64) (bool, error) {
	return false, nil
}

func (r *serviceAuthRepo) AddSecurityLog(context.Context, *biz.SecurityLog) error {
	return nil
}

func cloneServiceAccount(in *biz.IdentityAccount) *biz.IdentityAccount {
	if in == nil {
		return nil
	}
	cp := *in
	return &cp
}

func servicePasswordHash(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword() error = %v", err)
	}
	return string(hash)
}
