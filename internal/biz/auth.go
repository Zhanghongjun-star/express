package biz

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
	"time"

	v1 "shunfeng-miniprogram/api/auth/v1"

	"github.com/go-kratos/kratos/v3/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	// TargetTypePhone 表示验证码目标是手机号。
	TargetTypePhone = "phone"
	// TargetTypeEmail 表示验证码目标是邮箱。
	TargetTypeEmail = "email"
	// VerifySceneRegister 表示注册场景验证码。
	VerifySceneRegister = "register"

	// RoleUser 是注册成功后的默认角色。
	RoleUser = "user"
	// AccountStatusNormal 表示账号可正常登录和访问。
	AccountStatusNormal = "normal"
	// AccountStatusDisabled 表示账号已封禁。
	AccountStatusDisabled = "disabled"

	// RealNameStatusNotSubmitted 表示实名认证未提交。
	RealNameStatusNotSubmitted = "not_submitted"
	// RealNameStatusPending 表示实名认证审核中。
	RealNameStatusPending = "pending"

	accessTokenTTL  = 2 * time.Hour
	refreshTokenTTL = 30 * 24 * time.Hour
	loginLockTTL    = 30 * time.Minute
)

var (
	// ErrAuthInvalidArgument 表示请求参数不符合业务要求。
	ErrAuthInvalidArgument = errors.BadRequest(v1.ErrorReason_AUTH_INVALID_ARGUMENT.String(), "invalid auth argument")
	// ErrAuthUnauthenticated 表示账号、密码或 Token 无效。
	ErrAuthUnauthenticated = errors.Unauthorized(v1.ErrorReason_AUTH_UNAUTHENTICATED.String(), "auth unauthenticated")
	// ErrAuthPermissionDenied 表示当前账号无权访问目标资源。
	ErrAuthPermissionDenied = errors.Forbidden(v1.ErrorReason_AUTH_PERMISSION_DENIED.String(), "permission denied")
	// ErrAuthAccountDisabled 表示账号已封禁。
	ErrAuthAccountDisabled = errors.Forbidden(v1.ErrorReason_AUTH_ACCOUNT_DISABLED.String(), "account disabled")
	// ErrAuthNotFound 表示身份资源不存在。
	ErrAuthNotFound = errors.NotFound(v1.ErrorReason_AUTH_RESOURCE_NOT_FOUND.String(), "auth resource not found")
	// ErrAuthDuplicate 表示手机号、邮箱或身份证已存在。
	ErrAuthDuplicate = errors.Conflict(v1.ErrorReason_AUTH_DUPLICATE_REQUEST.String(), "auth resource duplicated")
	// ErrAuthVerificationRateLimited 表示验证码发送过于频繁。
	ErrAuthVerificationRateLimited = errors.TooManyRequests(v1.ErrorReason_AUTH_DUPLICATE_REQUEST.String(), "verification code sent too frequently")
	// ErrAuthInvalidStatus 表示账号或认证状态不允许当前操作。
	ErrAuthInvalidStatus = errors.Conflict(v1.ErrorReason_AUTH_INVALID_STATUS.String(), "invalid auth status")
)

// IdentityAccount 是身份账号领域对象，不带数据库标签，供 biz/service 使用。
type IdentityAccount struct {
	ID            int64
	UserNo        string
	Phone         string
	PhoneCipher   string
	PhoneHash     string
	Email         string
	PasswordHash  string
	RoleCode      string
	AccountStatus string
	LockedUntil   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// VerificationCode 是验证码领域对象，Redis 或短信服务细节留在 data 层。
type VerificationCode struct {
	VerificationID string
	Target         string
	TargetType     string
	Scene          string
	Code           string
	ExpiresAt      time.Time
}

// TokenSession 是登录会话领域对象，refresh token 与设备绑定。
type TokenSession struct {
	UserID       int64
	DeviceID     string
	AccessToken  string
	RefreshToken string
	RoleCode     string
	ExpiresAt    time.Time
}

// RealNameAuth 是实名认证领域对象，敏感原文不向 service 返回。
type RealNameAuth struct {
	AuthID         int64
	UserID         int64
	RealNameCipher string
	IDCardCipher   string
	IDCardHash     string
	RealNameMask   string
	IDCardMask     string
	ImageURLs      []string
	AuthStatus     string
	RejectReason   string
}

// SecurityLog 记录登录、改密、退出等安全事件。
type SecurityLog struct {
	UserID        int64
	EventType     string
	Result        string
	IP            string
	DeviceID      string
	FailureReason string
	CreatedAt     time.Time
}

// RegisterCommand 是注册输入。
type RegisterCommand struct {
	Phone      string
	Email      string
	VerifyCode string
	Password   string
	DeviceID   string
}

// LoginCommand 是账号密码登录输入。
type LoginCommand struct {
	Account  string
	Password string
	DeviceID string
	IP       string
}

// ChangePasswordCommand 是修改密码输入。
type ChangePasswordCommand struct {
	UserID      int64
	OldPassword string
	NewPassword string
}

// RealNameAuthCommand 是提交实名认证输入。
type RealNameAuthCommand struct {
	UserID    int64
	RealName  string
	IDCardNo  string
	ImageURLs []string
}

// TokenResult 是注册、登录、刷新后的 Token 输出。
type TokenResult struct {
	UserID        int64
	AccessToken   string
	RefreshToken  string
	RoleCode      string
	AccountStatus string
	ExpiresIn     int64
}

// AuthRepo 是身份业务仓储接口，由 data 层用 MySQL/Redis 实现。
type AuthRepo interface {
	SaveVerificationCode(context.Context, *VerificationCode) error
	GetVerificationCode(context.Context, string, string, string) (*VerificationCode, error)
	FindUserByAccount(context.Context, string) (*IdentityAccount, error)
	FindUserByID(context.Context, int64) (*IdentityAccount, error)
	CreateUser(context.Context, *IdentityAccount) (*IdentityAccount, error)
	RecordLoginFailure(context.Context, int64, time.Duration) (int, error)
	ClearLoginFailures(context.Context, int64) error
	UpdateLoginLock(context.Context, int64, *time.Time) error
	SaveTokenSession(context.Context, *TokenSession) error
	GetTokenSession(context.Context, string) (*TokenSession, error)
	DeleteTokenSession(context.Context, string) error
	DeleteUserTokenSessions(context.Context, int64) error
	BlacklistAccessToken(context.Context, string, time.Duration) error
	ValidateAccessToken(context.Context, string) (*TokenSession, error)
	UpdatePassword(context.Context, int64, string) error
	SaveRealNameAuth(context.Context, *RealNameAuth) (*RealNameAuth, error)
	GetRealNameAuth(context.Context, int64) (*RealNameAuth, error)
	IDCardHashExists(context.Context, string, int64) (bool, error)
	AddSecurityLog(context.Context, *SecurityLog) error
}

// AuthUsecase 聚合身份、登录、Token、改密和实名认证规则。
type AuthUsecase struct {
	repo AuthRepo
}

// NewAuthUsecase 创建身份用例。
func NewAuthUsecase(repo AuthRepo) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

// SendVerificationCode 生成 5 分钟有效验证码，60 秒限频由 data/Redis 保证。
func (uc *AuthUsecase) SendVerificationCode(ctx context.Context, target, targetType, scene string) (*VerificationCode, error) {
	target = strings.TrimSpace(target)
	if target == "" || !validTargetType(targetType) || strings.TrimSpace(scene) == "" {
		return nil, ErrAuthInvalidArgument
	}
	code, err := randomDigits(6)
	if err != nil {
		return nil, err
	}
	vc := &VerificationCode{
		VerificationID: newToken(),
		Target:         target,
		TargetType:     targetType,
		Scene:          scene,
		Code:           code,
		ExpiresAt:      time.Now().Add(5 * time.Minute),
	}
	if err := uc.repo.SaveVerificationCode(ctx, vc); err != nil {
		return nil, err
	}
	return vc, nil
}

// Register 校验验证码和账号唯一性，创建默认 user 角色账号并签发 Token。
func (uc *AuthUsecase) Register(ctx context.Context, cmd *RegisterCommand) (*TokenResult, error) {
	if cmd == nil || !validPassword(cmd.Password) || strings.TrimSpace(cmd.DeviceID) == "" {
		return nil, ErrAuthInvalidArgument
	}
	target, targetType, err := registerTarget(cmd.Phone, cmd.Email)
	if err != nil {
		return nil, err
	}
	if existing, err := uc.repo.FindUserByAccount(ctx, target); err == nil && existing != nil {
		return nil, ErrAuthDuplicate
	} else if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	vc, err := uc.repo.GetVerificationCode(ctx, target, targetType, VerifySceneRegister)
	if err != nil {
		return nil, err
	}
	if vc.Code != cmd.VerifyCode || time.Now().After(vc.ExpiresAt) {
		return nil, ErrAuthInvalidArgument
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user := &IdentityAccount{
		UserNo:        "U" + newShortToken(),
		Phone:         strings.TrimSpace(cmd.Phone),
		PhoneCipher:   simpleCipher(strings.TrimSpace(cmd.Phone)),
		Email:         strings.TrimSpace(cmd.Email),
		PasswordHash:  string(passwordHash),
		RoleCode:      RoleUser,
		AccountStatus: AccountStatusNormal,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if user.Phone != "" {
		user.PhoneHash = hashText(user.Phone)
	}
	created, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return uc.issueToken(ctx, created, cmd.DeviceID)
}

// Login 校验账号密码；连续 5 次错误后锁定 30 分钟。
func (uc *AuthUsecase) Login(ctx context.Context, cmd *LoginCommand) (*TokenResult, error) {
	if cmd == nil || strings.TrimSpace(cmd.Account) == "" || strings.TrimSpace(cmd.Password) == "" || strings.TrimSpace(cmd.DeviceID) == "" {
		return nil, ErrAuthInvalidArgument
	}
	user, err := uc.repo.FindUserByAccount(ctx, strings.TrimSpace(cmd.Account))
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, ErrAuthNotFound
		}
		return nil, err
	}
	if user.AccountStatus == AccountStatusDisabled {
		return nil, ErrAuthAccountDisabled
	}
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return nil, ErrAuthUnauthenticated
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cmd.Password)) != nil {
		count, _ := uc.repo.RecordLoginFailure(ctx, user.ID, loginLockTTL)
		if count >= 5 {
			lockedUntil := time.Now().Add(loginLockTTL)
			_ = uc.repo.UpdateLoginLock(ctx, user.ID, &lockedUntil)
		}
		_ = uc.repo.AddSecurityLog(ctx, &SecurityLog{UserID: user.ID, EventType: "login", Result: "failed", IP: cmd.IP, DeviceID: cmd.DeviceID, FailureReason: "bad_password", CreatedAt: time.Now()})
		return nil, ErrAuthUnauthenticated
	}
	_ = uc.repo.ClearLoginFailures(ctx, user.ID)
	_ = uc.repo.UpdateLoginLock(ctx, user.ID, nil)
	_ = uc.repo.AddSecurityLog(ctx, &SecurityLog{UserID: user.ID, EventType: "login", Result: "success", IP: cmd.IP, DeviceID: cmd.DeviceID, CreatedAt: time.Now()})
	return uc.issueToken(ctx, user, cmd.DeviceID)
}

// RefreshToken 使用有效 refresh token 换取新 Token。
func (uc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken, deviceID string) (*TokenResult, error) {
	if strings.TrimSpace(refreshToken) == "" || strings.TrimSpace(deviceID) == "" {
		return nil, ErrAuthInvalidArgument
	}
	session, err := uc.repo.GetTokenSession(ctx, refreshToken)
	if err != nil || session.DeviceID != deviceID || time.Now().After(session.ExpiresAt) {
		return nil, ErrAuthUnauthenticated
	}
	user, err := uc.repo.FindUserByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	_ = uc.repo.DeleteTokenSession(ctx, refreshToken)
	return uc.issueToken(ctx, user, deviceID)
}

// ValidateAccessToken 校验 access token 有效性，返回会话信息。
// 同时检查账号状态，封禁账号拒绝所有请求。
func (uc *AuthUsecase) ValidateAccessToken(ctx context.Context, accessToken string) (*TokenSession, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, ErrAuthInvalidArgument
	}
	session, err := uc.repo.ValidateAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, ErrAuthUnauthenticated
	}
	user, err := uc.repo.FindUserByID(ctx, session.UserID)
	if err != nil {
		return nil, ErrAuthUnauthenticated
	}
	if user.AccountStatus == AccountStatusDisabled {
		return nil, ErrAuthAccountDisabled
	}
	return session, nil
}

// Logout 退出当前设备或全部设备，并将 access token 加入黑名单。
func (uc *AuthUsecase) Logout(ctx context.Context, userID int64, accessToken, refreshToken, deviceID string, logoutAll bool) error {
	if userID <= 0 || strings.TrimSpace(accessToken) == "" || strings.TrimSpace(deviceID) == "" {
		return ErrAuthInvalidArgument
	}
	if logoutAll {
		if err := uc.repo.DeleteUserTokenSessions(ctx, userID); err != nil {
			return err
		}
	} else if strings.TrimSpace(refreshToken) != "" {
		if err := uc.repo.DeleteTokenSession(ctx, refreshToken); err != nil {
			return err
		}
	}
	if err := uc.repo.BlacklistAccessToken(ctx, accessToken, accessTokenTTL); err != nil {
		return err
	}
	return uc.repo.AddSecurityLog(ctx, &SecurityLog{UserID: userID, EventType: "logout", Result: "success", DeviceID: deviceID, CreatedAt: time.Now()})
}

// ChangePassword 校验旧密码，新旧密码不能相同，成功后清理全部会话。
func (uc *AuthUsecase) ChangePassword(ctx context.Context, cmd *ChangePasswordCommand) error {
	if cmd == nil || cmd.UserID <= 0 || !validPassword(cmd.NewPassword) || strings.TrimSpace(cmd.OldPassword) == "" {
		return ErrAuthInvalidArgument
	}
	if cmd.OldPassword == cmd.NewPassword {
		return ErrAuthInvalidArgument
	}
	user, err := uc.repo.FindUserByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cmd.OldPassword)) != nil {
		return ErrAuthUnauthenticated
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cmd.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := uc.repo.UpdatePassword(ctx, cmd.UserID, string(passwordHash)); err != nil {
		return err
	}
	if err := uc.repo.DeleteUserTokenSessions(ctx, cmd.UserID); err != nil {
		return err
	}
	return uc.repo.AddSecurityLog(ctx, &SecurityLog{UserID: cmd.UserID, EventType: "change_password", Result: "success", CreatedAt: time.Now()})
}

// SubmitRealNameAuth 保存脱敏、Hash 和密文字段，重复身份证不允许提交。
func (uc *AuthUsecase) SubmitRealNameAuth(ctx context.Context, cmd *RealNameAuthCommand) (*RealNameAuth, error) {
	if cmd == nil || cmd.UserID <= 0 || strings.TrimSpace(cmd.RealName) == "" || !validIDCard(cmd.IDCardNo) {
		return nil, ErrAuthInvalidArgument
	}
	if _, err := uc.repo.FindUserByID(ctx, cmd.UserID); err != nil {
		return nil, err
	}
	idHash := hashText(strings.ToUpper(strings.TrimSpace(cmd.IDCardNo)))
	exists, err := uc.repo.IDCardHashExists(ctx, idHash, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrAuthDuplicate
	}
	auth := &RealNameAuth{
		UserID:         cmd.UserID,
		RealNameCipher: simpleCipher(strings.TrimSpace(cmd.RealName)),
		IDCardCipher:   simpleCipher(strings.ToUpper(strings.TrimSpace(cmd.IDCardNo))),
		IDCardHash:     idHash,
		RealNameMask:   maskRealName(cmd.RealName),
		IDCardMask:     maskIDCard(cmd.IDCardNo),
		ImageURLs:      append([]string(nil), cmd.ImageURLs...),
		AuthStatus:     RealNameStatusPending,
	}
	return uc.repo.SaveRealNameAuth(ctx, auth)
}

// GetRealNameAuth 返回实名认证状态；未提交时返回默认状态。
func (uc *AuthUsecase) GetRealNameAuth(ctx context.Context, userID int64) (*RealNameAuth, error) {
	if userID <= 0 {
		return nil, ErrAuthInvalidArgument
	}
	auth, err := uc.repo.GetRealNameAuth(ctx, userID)
	if errors.IsNotFound(err) {
		return &RealNameAuth{UserID: userID, AuthStatus: RealNameStatusNotSubmitted}, nil
	}
	return auth, err
}

func (uc *AuthUsecase) issueToken(ctx context.Context, user *IdentityAccount, deviceID string) (*TokenResult, error) {
	accessToken := newToken()
	refreshToken := newToken()
	session := &TokenSession{
		UserID:       user.ID,
		DeviceID:     deviceID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RoleCode:     user.RoleCode,
		ExpiresAt:    time.Now().Add(refreshTokenTTL),
	}
	if err := uc.repo.SaveTokenSession(ctx, session); err != nil {
		return nil, err
	}
	return &TokenResult{
		UserID:        user.ID,
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		RoleCode:      user.RoleCode,
		AccountStatus: user.AccountStatus,
		ExpiresIn:     int64(accessTokenTTL.Seconds()),
	}, nil
}

func registerTarget(phone, email string) (string, string, error) {
	phone = strings.TrimSpace(phone)
	email = strings.TrimSpace(email)
	switch {
	case phone != "":
		return phone, TargetTypePhone, nil
	case email != "":
		return email, TargetTypeEmail, nil
	default:
		return "", "", ErrAuthInvalidArgument
	}
}

func validPassword(password string) bool {
	n := len([]rune(password))
	return n >= 8 && n <= 32
}

func validTargetType(targetType string) bool {
	return targetType == TargetTypePhone || targetType == TargetTypeEmail
}

var idCardPattern = regexp.MustCompile(`^[0-9]{17}[0-9Xx]$`)

func validIDCard(idCardNo string) bool {
	return idCardPattern.MatchString(strings.TrimSpace(idCardNo))
}

func hashText(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}

func simpleCipher(text string) string {
	return hex.EncodeToString([]byte(text))
}

func maskRealName(name string) string {
	runes := []rune(strings.TrimSpace(name))
	if len(runes) <= 1 {
		return "*"
	}
	return string(runes[0]) + strings.Repeat("*", len(runes)-1)
}

func maskIDCard(idCardNo string) string {
	id := strings.ToUpper(strings.TrimSpace(idCardNo))
	if len(id) < 8 {
		return "****"
	}
	return id[:4] + "**********" + id[len(id)-4:]
}

func randomDigits(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = byte('0' + b[i]%10)
	}
	return string(b), nil
}

func newToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().String()))
	}
	return hex.EncodeToString(b)
}

func newShortToken() string {
	token := newToken()
	if len(token) > 12 {
		return token[:12]
	}
	return token
}

func cloneIdentityAccount(in *IdentityAccount) *IdentityAccount {
	if in == nil {
		return nil
	}
	cp := *in
	if in.LockedUntil != nil {
		t := *in.LockedUntil
		cp.LockedUntil = &t
	}
	return &cp
}
