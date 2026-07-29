package server

import (
	addressv1 "shunfeng-miniprogram/api/address/v1"
	authv1 "shunfeng-miniprogram/api/auth/v1"
	orderv1 "shunfeng-miniprogram/api/order/v1"
	todov1 "shunfeng-miniprogram/api/todo/v1"
	userv1 "shunfeng-miniprogram/api/user/v1"
	"shunfeng-miniprogram/internal/conf"
	"shunfeng-miniprogram/internal/service"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, todo *service.TodoService, user *service.UserService, address *service.AddressService, order *service.OrderService, auth *service.AuthService) *grpc.Server {

	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	todov1.RegisterTodoServiceServer(srv, todo)
	userv1.RegisterUserServiceServer(srv, user)
	orderv1.RegisterOrderServiceServer(srv, order)
	authv1.RegisterAuthServiceServer(srv, auth)
	if address != nil {
		addressv1.RegisterAddressServiceServer(srv, address)
	}
	return srv
}
