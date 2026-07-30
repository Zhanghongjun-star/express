package server

import (
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/middleware/validate"
	"github.com/go-kratos/kratos/v3/transport/http"
	addressv1 "shunfeng-miniprogram/api/address/v1"
	expressv1 "shunfeng-miniprogram/api/express/v1" // 查快递模块 HTTP 路由注册
	authv1 "shunfeng-miniprogram/api/auth/v1"
	orderv1 "shunfeng-miniprogram/api/order/v1"
	orderv2 "shunfeng-miniprogram/api/order/v2"
	shippingv1 "shunfeng-miniprogram/api/shipping/v1"
	todov1 "shunfeng-miniprogram/api/todo/v1"
	userv1 "shunfeng-miniprogram/api/user/v1"
	"shunfeng-miniprogram/internal/conf"
	"shunfeng-miniprogram/internal/service"


	"go.einride.tech/aip/fieldbehavior"
	"google.golang.org/protobuf/proto"
)

// NewHTTPServer new an HTTP server。
func NewHTTPServer(c *conf.Server, todo *service.TodoService, user *service.UserService, address *service.AddressService, order *service.OrderService, express *service.ExpressService, shipping *service.ShippingService, freight *service.FreightService, auth *service.AuthService) *http.Server {

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(func(req any) error {
				if msg, ok := req.(proto.Message); ok {
					if err := fieldbehavior.ValidateRequiredFields(msg); err != nil {
						return err
					}
				}
				return nil
			}),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	todov1.RegisterTodoServiceHTTPServer(srv, todo)
	userv1.RegisterUserServiceHTTPServer(srv, user)
	orderv1.RegisterOrderServiceHTTPServer(srv, order)
	shippingv1.RegisterShippingServiceHTTPServer(srv, shipping)
	orderv2.RegisterOrderServiceHTTPServer(srv, freight)
	authv1.RegisterAuthServiceHTTPServer(srv, auth)
	if address != nil {
		addressv1.RegisterAddressServiceHTTPServer(srv, address)
	}
	// 查快递模块: 注册 12 个 HTTP 接口（搜索/关注/标签/分享/消息）
	if express != nil {
		expressv1.RegisterExpressServiceHTTPServer(srv, express)
	}
	return srv
}
