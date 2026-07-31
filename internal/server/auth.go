package server

import (
	"context"
	"strings"

	"shunfeng-miniprogram/internal/biz"

	"github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/transport"

	gerrors "errors"
)

type contextKey string

const (
	CtxKeyUserID   contextKey = "x-user-id"
	CtxKeyRoleCode contextKey = "x-role-code"
)

var publicPaths = map[string]bool{
	"/auth.v1.AuthService/SendVerificationCode":    true,
	"/auth.v1.AuthService/Register":                true,
	"/auth.v1.AuthService/Login":                   true,
	"/auth.v1.AuthService/RefreshToken":            true,
	"/auth.v1.AuthService/Logout":                  true,
	"/order.v3.OrderService/HandlePaymentCallback": true,
}

func isPublicPath(path string) bool {
	if publicPaths[path] {
		return true
	}
	return strings.HasPrefix(path, "/api/v1/auth/")
}

func AuthMiddleware(uc *biz.AuthUsecase) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.Unauthorized("AUTH", "no transport")
			}
			path := tr.Operation()
			if isPublicPath(path) {
				return handler(ctx, req)
			}
			authHeader := tr.RequestHeader().Get("Authorization")
			if authHeader == "" {
				return nil, errors.Unauthorized("AUTH", "missing authorization header")
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return nil, errors.Unauthorized("AUTH", "invalid authorization format")
			}
			session, err := uc.ValidateAccessToken(ctx, parts[1])
			if err != nil {
				if gerrors.Is(err, biz.ErrAuthAccountDisabled) {
					return nil, err
				}
				return nil, errors.Unauthorized("AUTH", "invalid or expired token")
			}
			ctx = context.WithValue(ctx, CtxKeyUserID, session.UserID)
			ctx = context.WithValue(ctx, CtxKeyRoleCode, session.RoleCode)
			return handler(ctx, req)
		}
	}
}

func GetUserID(ctx context.Context) int64 {
	if v, ok := ctx.Value(CtxKeyUserID).(int64); ok {
		return v
	}
	return 0
}

func GetRoleCode(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyRoleCode).(string); ok {
		return v
	}
	return ""
}
