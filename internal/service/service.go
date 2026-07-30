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
	// 寄件模块服务
	NewShippingService,
	NewFreightService,
	NewAuthService,
)
