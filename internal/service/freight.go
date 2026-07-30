package service

import (
	"context"

	v2 "shunfeng-miniprogram/api/order/v2"
	"shunfeng-miniprogram/internal/biz"
)

// =============================================================================
// FreightService 运费预估服务
// =============================================================================

// FreightService 运费预估传输层适配器，负责 DTO ↔ DO 转换。
// 实现 v2.OrderServiceServer 接口。
type FreightService struct {
	v2.UnimplementedOrderServiceServer

	uc *biz.FreightUsecase
}

// NewFreightService 创建运费预估服务。
func NewFreightService(uc *biz.FreightUsecase) *FreightService {
	return &FreightService{uc: uc}
}

// EstimateFreight 运费预估接口。
// 将 proto 请求转换为领域对象，调用用例层计算，再将结果回传为 proto 响应。
func (s *FreightService) EstimateFreight(ctx context.Context, req *v2.EstimateFreightRequest) (*v2.EstimateFreightReply, error) {
	// 将 proto DTO 转换为领域对象
	freightReq := &biz.FreightRequest{
		SenderProvince:   req.GetSenderProvince(),
		SenderCity:       req.GetSenderCity(),
		SenderArea:       req.GetSenderArea(),
		ReceiverProvince: req.GetReceiverProvince(),
		ReceiverCity:     req.GetReceiverCity(),
		ReceiverArea:     req.GetReceiverArea(),
		Weight:           req.GetWeight(),
		Length:           req.GetLength(),
		Width:            req.GetWidth(),
		Height:           req.GetHeight(),
		ExpressType:      convertExpressType(req.GetExpressType()),
		InsureValue:      req.GetInsureValue(),
	}

	// 调用用例层计算运费
	result, err := s.uc.Estimate(ctx, freightReq)
	if err != nil {
		return nil, err
	}

	// 将领域对象转换为 proto 响应
	return &v2.EstimateFreightReply{
		BaseFreight: result.BaseFreight,
		InsureFee:   result.InsureFee,
		TotalPrice:  result.TotalPrice,
		CalcWeight:  result.CalcWeight,
		Tips:        result.Tips,
	}, nil
}

// convertExpressType 将 proto 中的快递类型 int32 转换为 biz 层的 ExpressType。
func convertExpressType(t int32) biz.ExpressType {
	return biz.ExpressType(t)
}
