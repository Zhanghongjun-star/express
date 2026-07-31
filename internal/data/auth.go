package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"shunfeng-miniprogram/internal/biz"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

// authRepo 实现身份业务仓储，MySQL 保存账号资料，Redis 保存短期状态。
type authRepo struct {
}

// NewAuthRepo 创建身份仓储。
func NewAuthRepo() biz.AuthRepo {
	return &authRepo{}
}

func (r *authRepo) SaveVerificationCode(ctx context.Context, vc *biz.VerificationCode) error {
	if RDB == nil {
		return biz.ErrAuthInvalidStatus
	}
	limitKey := verificationLimitKey(vc.Target, vc.TargetType, vc.Scene)
	ok, err := RDB.SetNX(ctx, limitKey, "1", time.Minute).Result()
	if err != nil {
		return err
	}
	if !ok {
		return biz.ErrAuthVerificationRateLimited
	}
	body, err := json.Marshal(vc)
	if err != nil {
		return err
	}
	if err := RDB.Set(ctx, verificationKey(vc.Target, vc.TargetType, vc.Scene), body, time.Until(vc.ExpiresAt)).Err(); err != nil {
		return err
	}
	if vc.TargetType == biz.TargetTypePhone {
		msg := fmt.Sprintf("您的验证码是：%s。请不要把验证码泄露给其他人。", vc.Code)
		if sendErr := SendSMS(ctx, vc.Target, msg); sendErr != nil {
			fmt.Println(fmt.Sprintf("sms 发送失败: %v", sendErr))
		}
	}
	return nil
}

func (r *authRepo) GetVerificationCode(ctx context.Context, target, targetType, scene string) (*biz.VerificationCode, error) {
	if RDB == nil {
		return nil, biz.ErrAuthInvalidStatus
	}
	body, err := RDB.Get(ctx, verificationKey(target, targetType, scene)).Bytes()
	if err == redis.Nil {
		return nil, biz.ErrAuthInvalidArgument
	}
	if err != nil {
		return nil, err
	}
	var vc biz.VerificationCode
	if err := json.Unmarshal(body, &vc); err != nil {
		return nil, err
	}
	return &vc, nil
}

func (r *authRepo) FindUserByAccount(ctx context.Context, account string) (*biz.IdentityAccount, error) {
	db, err := authSQLDB()
	if err != nil {
		return nil, err
	}
	const query = `
SELECT id, user_no, phone_cipher, phone_hash, email, password_hash, role_code, account_status, locked_until, created_at, updated_at
FROM identity_user
WHERE email = ? OR phone_hash = ?
LIMIT 1`
	return scanIdentityAccount(db.QueryRowContext(ctx, query, account, hashTextData(account)))
}

func (r *authRepo) FindUserByID(ctx context.Context, userID int64) (*biz.IdentityAccount, error) {
	db, err := authSQLDB()
	if err != nil {
		return nil, err
	}
	const query = `
SELECT id, user_no, phone_cipher, phone_hash, email, password_hash, role_code, account_status, locked_until, created_at, updated_at
FROM identity_user
WHERE id = ?
LIMIT 1`
	return scanIdentityAccount(db.QueryRowContext(ctx, query, userID))
}

func (r *authRepo) CreateUser(ctx context.Context, user *biz.IdentityAccount) (*biz.IdentityAccount, error) {
	db, err := authSQLDB()
	if err != nil {
		return nil, err
	}
	const stmt = `
INSERT INTO identity_user (user_no, phone_cipher, phone_hash, email, password_hash, role_code, account_status, locked_until, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := db.ExecContext(ctx, stmt,
		user.UserNo,
		user.PhoneCipher,
		nullableString(user.PhoneHash),
		nullableString(user.Email),
		user.PasswordHash,
		user.RoleCode,
		user.AccountStatus,
		user.LockedUntil,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if isDuplicateEntry(err) {
		return nil, biz.ErrAuthDuplicate
	}
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = id
	return cloneIdentityAccountData(user), nil
}

func (r *authRepo) RecordLoginFailure(ctx context.Context, userID int64, ttl time.Duration) (int, error) {
	if RDB == nil {
		return 0, biz.ErrAuthInvalidStatus
	}
	key := loginFailureKey(userID)
	count, err := RDB.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	_ = RDB.Expire(ctx, key, ttl).Err()
	return int(count), nil
}

func (r *authRepo) ClearLoginFailures(ctx context.Context, userID int64) error {
	if RDB == nil {
		return nil
	}
	return RDB.Del(ctx, loginFailureKey(userID)).Err()
}

func (r *authRepo) UpdateLoginLock(ctx context.Context, userID int64, lockedUntil *time.Time) error {
	db, err := authSQLDB()
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, `UPDATE identity_user SET locked_until = ?, updated_at = ? WHERE id = ?`, lockedUntil, time.Now(), userID)
	return err
}

func (r *authRepo) SaveTokenSession(ctx context.Context, session *biz.TokenSession) error {
	if RDB == nil {
		return biz.ErrAuthInvalidStatus
	}
	body, err := json.Marshal(session)
	if err != nil {
		return err
	}
	ttl := time.Until(session.ExpiresAt)
	accessTTL := 2 * time.Hour
	pipe := RDB.TxPipeline()
	pipe.Set(ctx, tokenSessionKey(session.RefreshToken), body, ttl)
	pipe.Set(ctx, accessTokenKey(session.AccessToken), body, accessTTL)
	pipe.SAdd(ctx, userSessionsKey(session.UserID), session.RefreshToken)
	pipe.Expire(ctx, userSessionsKey(session.UserID), ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *authRepo) ValidateAccessToken(ctx context.Context, accessToken string) (*biz.TokenSession, error) {
	if RDB == nil {
		return nil, biz.ErrAuthInvalidStatus
	}
	blacklisted, err := RDB.Exists(ctx, accessBlacklistKey(accessToken)).Result()
	if err != nil {
		return nil, err
	}
	if blacklisted > 0 {
		return nil, biz.ErrAuthUnauthenticated
	}
	body, err := RDB.Get(ctx, accessTokenKey(accessToken)).Bytes()
	if err == redis.Nil {
		return nil, biz.ErrAuthUnauthenticated
	}
	if err != nil {
		return nil, err
	}
	var session biz.TokenSession
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *authRepo) GetTokenSession(ctx context.Context, refreshToken string) (*biz.TokenSession, error) {
	if RDB == nil {
		return nil, biz.ErrAuthInvalidStatus
	}
	body, err := RDB.Get(ctx, tokenSessionKey(refreshToken)).Bytes()
	if err == redis.Nil {
		return nil, biz.ErrAuthUnauthenticated
	}
	if err != nil {
		return nil, err
	}
	var session biz.TokenSession
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *authRepo) DeleteTokenSession(ctx context.Context, refreshToken string) error {
	if RDB == nil {
		return nil
	}
	session, _ := r.GetTokenSession(ctx, refreshToken)
	pipe := RDB.TxPipeline()
	pipe.Del(ctx, tokenSessionKey(refreshToken))
	if session != nil {
		pipe.SRem(ctx, userSessionsKey(session.UserID), refreshToken)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (r *authRepo) DeleteUserTokenSessions(ctx context.Context, userID int64) error {
	if RDB == nil {
		return nil
	}
	key := userSessionsKey(userID)
	tokens, err := RDB.SMembers(ctx, key).Result()
	if err != nil {
		return err
	}
	pipe := RDB.TxPipeline()
	for _, token := range tokens {
		pipe.Del(ctx, tokenSessionKey(token))
	}
	pipe.Del(ctx, key)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *authRepo) BlacklistAccessToken(ctx context.Context, accessToken string, ttl time.Duration) error {
	if RDB == nil {
		return nil
	}
	return RDB.Set(ctx, accessBlacklistKey(accessToken), "1", ttl).Err()
}

func (r *authRepo) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	db, err := authSQLDB()
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, `UPDATE identity_user SET password_hash = ?, updated_at = ? WHERE id = ?`, passwordHash, time.Now(), userID)
	return err
}

func (r *authRepo) SaveRealNameAuth(ctx context.Context, auth *biz.RealNameAuth) (*biz.RealNameAuth, error) {
	db, err := authSQLDB()
	if err != nil {
		return nil, err
	}
	imageBody, err := json.Marshal(auth.ImageURLs)
	if err != nil {
		return nil, err
	}
	const stmt = `
INSERT INTO identity_real_name_auth (user_id, real_name_cipher, id_card_cipher, id_card_hash, image_urls, auth_status, reject_reason, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  real_name_cipher = VALUES(real_name_cipher),
  id_card_cipher = VALUES(id_card_cipher),
  id_card_hash = VALUES(id_card_hash),
  image_urls = VALUES(image_urls),
  auth_status = VALUES(auth_status),
  reject_reason = VALUES(reject_reason),
  updated_at = VALUES(updated_at)`
	_, err = db.ExecContext(ctx, stmt,
		auth.UserID,
		auth.RealNameCipher,
		auth.IDCardCipher,
		auth.IDCardHash,
		string(imageBody),
		auth.AuthStatus,
		auth.RejectReason,
		time.Now(),
		time.Now(),
	)
	if isDuplicateEntry(err) {
		return nil, biz.ErrAuthDuplicate
	}
	if err != nil {
		return nil, err
	}
	return r.GetRealNameAuth(ctx, auth.UserID)
}

func (r *authRepo) GetRealNameAuth(ctx context.Context, userID int64) (*biz.RealNameAuth, error) {
	db, err := authSQLDB()
	if err != nil {
		return nil, err
	}
	const query = `
SELECT user_id, real_name_cipher, id_card_cipher, id_card_hash, image_urls, auth_status, reject_reason
FROM identity_real_name_auth
WHERE user_id = ?
LIMIT 1`
	row := db.QueryRowContext(ctx, query, userID)
	var auth biz.RealNameAuth
	var imageBody string
	if err := row.Scan(&auth.UserID, &auth.RealNameCipher, &auth.IDCardCipher, &auth.IDCardHash, &imageBody, &auth.AuthStatus, &auth.RejectReason); err != nil {
		if err == sql.ErrNoRows {
			return nil, biz.ErrAuthNotFound
		}
		return nil, err
	}
	auth.AuthID = auth.UserID
	_ = json.Unmarshal([]byte(imageBody), &auth.ImageURLs)
	auth.RealNameMask = maskRealNameData(decodeCipherData(auth.RealNameCipher))
	auth.IDCardMask = maskIDCardData(decodeCipherData(auth.IDCardCipher))
	return &auth, nil
}

func (r *authRepo) IDCardHashExists(ctx context.Context, idCardHash string, exceptUserID int64) (bool, error) {
	db, err := authSQLDB()
	if err != nil {
		return false, err
	}
	const query = `SELECT user_id FROM identity_real_name_auth WHERE id_card_hash = ? AND user_id <> ? LIMIT 1`
	var id int64
	err = db.QueryRowContext(ctx, query, idCardHash, exceptUserID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *authRepo) AddSecurityLog(ctx context.Context, log *biz.SecurityLog) error {
	if log == nil {
		return nil
	}
	db, err := authSQLDB()
	if err != nil {
		return err
	}
	const stmt = `
INSERT INTO identity_security_log (user_id, event_type, result, ip, device_id, failure_reason, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = db.ExecContext(ctx, stmt, log.UserID, log.EventType, log.Result, log.IP, log.DeviceID, log.FailureReason, log.CreatedAt)
	return err
}

func authSQLDB() (*sql.DB, error) {
	if DB == nil {
		return nil, biz.ErrAuthInvalidStatus
	}
	db, err := DB.DB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanIdentityAccount(row scanner) (*biz.IdentityAccount, error) {
	var user biz.IdentityAccount
	var phoneHash sql.NullString
	var email sql.NullString
	var lockedUntil sql.NullTime
	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	if err := row.Scan(
		&user.ID,
		&user.UserNo,
		&user.PhoneCipher,
		&phoneHash,
		&email,
		&user.PasswordHash,
		&user.RoleCode,
		&user.AccountStatus,
		&lockedUntil,
		&createdAt,
		&updatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, biz.ErrAuthNotFound
		}
		return nil, err
	}
	user.Phone = decodeCipherData(user.PhoneCipher)
	user.PhoneHash = phoneHash.String
	user.Email = email.String
	if lockedUntil.Valid {
		user.LockedUntil = &lockedUntil.Time
	}
	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}
	return &user, nil
}

func nullableString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func isDuplicateEntry(err error) bool {
	if err == nil {
		return false
	}
	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
		return true
	}
	return strings.Contains(err.Error(), "Duplicate entry")
}

func hashTextData(text string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(text)))
	return hex.EncodeToString(sum[:])
}

func decodeCipherData(text string) string {
	body, err := hex.DecodeString(text)
	if err != nil {
		return ""
	}
	return string(body)
}

func maskRealNameData(name string) string {
	runes := []rune(strings.TrimSpace(name))
	if len(runes) <= 1 {
		return "*"
	}
	return string(runes[0]) + strings.Repeat("*", len(runes)-1)
}

func maskIDCardData(idCardNo string) string {
	id := strings.ToUpper(strings.TrimSpace(idCardNo))
	if len(id) < 8 {
		return "****"
	}
	return id[:4] + "**********" + id[len(id)-4:]
}

func cloneIdentityAccountData(in *biz.IdentityAccount) *biz.IdentityAccount {
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

func verificationKey(target, targetType, scene string) string {
	return "auth:verification:" + targetType + ":" + scene + ":" + target
}

func verificationLimitKey(target, targetType, scene string) string {
	return "auth:verification_limit:" + targetType + ":" + scene + ":" + target
}

func loginFailureKey(userID int64) string {
	return "auth:login_failures:" + hex.EncodeToString([]byte(time.Unix(userID, 0).String()))
}

func tokenSessionKey(refreshToken string) string {
	return "auth:refresh:" + refreshToken
}

func userSessionsKey(userID int64) string {
	return "auth:user_sessions:" + hex.EncodeToString([]byte(time.Unix(userID, 0).String()))
}

func accessBlacklistKey(accessToken string) string {
	return "auth:access_blacklist:" + accessToken
}

func accessTokenKey(accessToken string) string {
	return "auth:access:" + accessToken
}
