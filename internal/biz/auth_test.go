package biz

import (
	"context"
	"testing"
	"time"

	kratoserrors "github.com/go-kratos/kratos/v3/errors"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUsecaseRegisterRejectsDuplicateAccount(t *testing.T) {
	repo := newFakeAuthRepo()
	repo.verification = &VerificationCode{
		Target:     "13800138000",
		TargetType: TargetTypePhone,
		Scene:      VerifySceneRegister,
		Code:       "123456",
		ExpiresAt:  time.Now().Add(time.Minute),
	}
	repo.usersByAccount["13800138000"] = &IdentityAccount{ID: 1, Phone: "13800138000"}
	uc := NewAuthUsecase(repo)

	_, err := uc.Register(context.Background(), &RegisterCommand{
		Phone:      "13800138000",
		VerifyCode: "123456",
		Password:   "password123",
		DeviceID:   "ios-1",
	})

	if !kratoserrors.IsConflict(err) {
		t.Fatalf("Register() error = %v, want conflict", err)
	}
}

func TestAuthUsecaseRegisterRejectsShortPassword(t *testing.T) {
	uc := NewAuthUsecase(newFakeAuthRepo())

	_, err := uc.Register(context.Background(), &RegisterCommand{
		Phone:      "13800138000",
		VerifyCode: "123456",
		Password:   "short",
		DeviceID:   "ios-1",
	})

	if !kratoserrors.IsBadRequest(err) {
		t.Fatalf("Register(short password) error = %v, want bad request", err)
	}
}

func TestAuthUsecaseLoginLocksAccountAfterFiveFailures(t *testing.T) {
	repo := newFakeAuthRepo()
	repo.usersByAccount["13800138000"] = &IdentityAccount{
		ID:            1,
		Phone:         "13800138000",
		PasswordHash:  mustPasswordHash(t, "right-password"),
		RoleCode:      RoleUser,
		AccountStatus: AccountStatusNormal,
	}
	uc := NewAuthUsecase(repo)

	for i := 0; i < 5; i++ {
		_, err := uc.Login(context.Background(), &LoginCommand{
			Account:  "13800138000",
			Password: "wrong-password",
			DeviceID: "ios-1",
		})
		if !kratoserrors.IsUnauthorized(err) {
			t.Fatalf("Login wrong password #%d error = %v, want unauthorized", i+1, err)
		}
	}

	user := repo.usersByAccount["13800138000"]
	if user.LockedUntil == nil || time.Until(*user.LockedUntil) < 29*time.Minute {
		t.Fatalf("LockedUntil = %v, want about 30 minutes in the future", user.LockedUntil)
	}
}

func TestAuthUsecaseLoginReturnsNotFoundForMissingAccount(t *testing.T) {
	uc := NewAuthUsecase(newFakeAuthRepo())

	_, err := uc.Login(context.Background(), &LoginCommand{
		Account:  "13900139000",
		Password: "password123",
		DeviceID: "ios-1",
	})

	if !kratoserrors.IsNotFound(err) {
		t.Fatalf("Login(missing account) error = %v, want not found", err)
	}
}

func TestAuthUsecaseChangePasswordRejectsWrongOldPassword(t *testing.T) {
	repo := newFakeAuthRepo()
	repo.usersByID[7] = &IdentityAccount{
		ID:            7,
		PasswordHash:  mustPasswordHash(t, "old-password"),
		AccountStatus: AccountStatusNormal,
	}
	uc := NewAuthUsecase(repo)

	err := uc.ChangePassword(context.Background(), &ChangePasswordCommand{
		UserID:      7,
		OldPassword: "wrong-password",
		NewPassword: "new-password",
	})

	if !kratoserrors.IsUnauthorized(err) {
		t.Fatalf("ChangePassword() error = %v, want unauthorized", err)
	}
}

func TestAuthUsecaseSubmitRealNameAuthRejectsInvalidIDCard(t *testing.T) {
	repo := newFakeAuthRepo()
	repo.usersByID[9] = &IdentityAccount{ID: 9, AccountStatus: AccountStatusNormal}
	uc := NewAuthUsecase(repo)

	_, err := uc.SubmitRealNameAuth(context.Background(), &RealNameAuthCommand{
		UserID:   9,
		RealName: "张三",
		IDCardNo: "123",
	})

	if !kratoserrors.IsBadRequest(err) {
		t.Fatalf("SubmitRealNameAuth(invalid id card) error = %v, want bad request", err)
	}
}

func mustPasswordHash(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword() error = %v", err)
	}
	return string(hash)
}

type fakeAuthRepo struct {
	usersByID      map[int64]*IdentityAccount
	usersByAccount map[string]*IdentityAccount
	verification   *VerificationCode
	loginFailures  map[int64]int
	realNameAuth   map[int64]*RealNameAuth
}

func newFakeAuthRepo() *fakeAuthRepo {
	return &fakeAuthRepo{
		usersByID:      make(map[int64]*IdentityAccount),
		usersByAccount: make(map[string]*IdentityAccount),
		loginFailures:  make(map[int64]int),
		realNameAuth:   make(map[int64]*RealNameAuth),
	}
}

func (r *fakeAuthRepo) SaveVerificationCode(context.Context, *VerificationCode) error {
	return nil
}

func (r *fakeAuthRepo) GetVerificationCode(_ context.Context, target, targetType, scene string) (*VerificationCode, error) {
	if r.verification == nil {
		return nil, ErrAuthInvalidArgument
	}
	if r.verification.Target != target || r.verification.TargetType != targetType || r.verification.Scene != scene {
		return nil, ErrAuthInvalidArgument
	}
	cp := *r.verification
	return &cp, nil
}

func (r *fakeAuthRepo) FindUserByAccount(_ context.Context, account string) (*IdentityAccount, error) {
	if user, ok := r.usersByAccount[account]; ok {
		return cloneIdentityAccount(user), nil
	}
	return nil, ErrAuthNotFound
}

func (r *fakeAuthRepo) FindUserByID(_ context.Context, userID int64) (*IdentityAccount, error) {
	if user, ok := r.usersByID[userID]; ok {
		return cloneIdentityAccount(user), nil
	}
	return nil, ErrAuthNotFound
}

func (r *fakeAuthRepo) CreateUser(_ context.Context, user *IdentityAccount) (*IdentityAccount, error) {
	cp := cloneIdentityAccount(user)
	if cp.ID == 0 {
		cp.ID = int64(len(r.usersByID) + len(r.usersByAccount) + 1)
	}
	r.usersByID[cp.ID] = cloneIdentityAccount(cp)
	if cp.Phone != "" {
		r.usersByAccount[cp.Phone] = cloneIdentityAccount(cp)
	}
	if cp.Email != "" {
		r.usersByAccount[cp.Email] = cloneIdentityAccount(cp)
	}
	return cloneIdentityAccount(cp), nil
}

func (r *fakeAuthRepo) RecordLoginFailure(_ context.Context, userID int64, _ time.Duration) (int, error) {
	r.loginFailures[userID]++
	return r.loginFailures[userID], nil
}

func (r *fakeAuthRepo) ClearLoginFailures(_ context.Context, userID int64) error {
	delete(r.loginFailures, userID)
	return nil
}

func (r *fakeAuthRepo) UpdateLoginLock(_ context.Context, userID int64, lockedUntil *time.Time) error {
	if user, ok := r.usersByID[userID]; ok {
		user.LockedUntil = lockedUntil
	}
	for _, user := range r.usersByAccount {
		if user.ID == userID {
			user.LockedUntil = lockedUntil
		}
	}
	return nil
}

func (r *fakeAuthRepo) SaveTokenSession(context.Context, *TokenSession) error {
	return nil
}

func (r *fakeAuthRepo) GetTokenSession(context.Context, string) (*TokenSession, error) {
	return nil, ErrAuthUnauthenticated
}

func (r *fakeAuthRepo) DeleteTokenSession(context.Context, string) error {
	return nil
}

func (r *fakeAuthRepo) DeleteUserTokenSessions(context.Context, int64) error {
	return nil
}

func (r *fakeAuthRepo) BlacklistAccessToken(context.Context, string, time.Duration) error {
	return nil
}

func (r *fakeAuthRepo) ValidateAccessToken(_ context.Context, token string) (*TokenSession, error) {
	if token == "" {
		return nil, ErrAuthUnauthenticated
	}
	sessions := make(map[int64]*TokenSession)
	for _, user := range r.usersByID {
		sessions[user.ID] = &TokenSession{
			UserID:   user.ID,
			RoleCode: user.RoleCode,
		}
	}
	for _, s := range sessions {
		return s, nil
	}
	return nil, ErrAuthUnauthenticated
}

func (r *fakeAuthRepo) UpdatePassword(_ context.Context, userID int64, passwordHash string) error {
	if user, ok := r.usersByID[userID]; ok {
		user.PasswordHash = passwordHash
	}
	return nil
}

func (r *fakeAuthRepo) SaveRealNameAuth(_ context.Context, auth *RealNameAuth) (*RealNameAuth, error) {
	cp := *auth
	cp.AuthID = int64(len(r.realNameAuth) + 1)
	r.realNameAuth[auth.UserID] = &cp
	return &cp, nil
}

func (r *fakeAuthRepo) GetRealNameAuth(_ context.Context, userID int64) (*RealNameAuth, error) {
	if auth, ok := r.realNameAuth[userID]; ok {
		cp := *auth
		return &cp, nil
	}
	return nil, ErrAuthNotFound
}

func (r *fakeAuthRepo) IDCardHashExists(_ context.Context, idCardHash string, exceptUserID int64) (bool, error) {
	for userID, auth := range r.realNameAuth {
		if userID != exceptUserID && auth.IDCardHash == idCardHash {
			return true, nil
		}
	}
	return false, nil
}

func (r *fakeAuthRepo) AddSecurityLog(context.Context, *SecurityLog) error {
	return nil
}
