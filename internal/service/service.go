package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewTodoService,
	NewAddressService,
	NewUserService,
	NewOrderService,
	// 查快递模块服务
	NewExpressService,
)
var ProviderSet = wire.NewSet(NewTodoService, NewAddressService, NewUserService, NewOrderService, NewShippingService, NewFreightService, NewAuthService)
