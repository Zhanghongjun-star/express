package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewTodoUsecase,
	NewAddressUsecase,
	NewUserUsecase,
	NewOrderUsecase,
	// 查快递模块：关注·标签·分享·消息
	NewOrderFollowUsecase,
	NewOrderLabelUsecase,
	NewOrderShareUsecase,
	NewUserMessageUsecase,
)
