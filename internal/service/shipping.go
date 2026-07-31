package service

import (
	"context"

	v1 "shunfeng-miniprogram/api/shipping/v1"
	"shunfeng-miniprogram/internal/biz"
)

type ShippingService struct {
	v1.UnimplementedShippingServiceServer

	channelUc       *biz.ChannelUsecase
	lockerUc        *biz.LockerUsecase
	servicePointUc  *biz.ServicePointUsecase
	pickupUc        *biz.PickupUsecase
}

func NewShippingService(channelUc *biz.ChannelUsecase, lockerUc *biz.LockerUsecase, servicePointUc *biz.ServicePointUsecase, pickupUc *biz.PickupUsecase) *ShippingService {
	return &ShippingService{
		channelUc:      channelUc,
		lockerUc:       lockerUc,
		servicePointUc: servicePointUc,
		pickupUc:       pickupUc,
	}
}

func (s *ShippingService) ListChannels(ctx context.Context, req *v1.ListChannelsRequest) (*v1.ListChannelsResponse, error) {
	channels, areas, err := s.channelUc.ListChannels(ctx, req.GetDistrictCode())
	if err != nil {
		return nil, err
	}

	areaMap := make(map[string]*biz.ChannelArea)
	for _, a := range areas {
		areaMap[a.ChannelCode] = a
	}

	resp := &v1.ListChannelsResponse{
		Channels: make([]*v1.Channel, 0, len(channels)),
	}
	for _, ch := range channels {
		area, hasArea := areaMap[ch.ChannelCode]
		available := hasArea && area.Status == 1

		v := &v1.Channel{
			ChannelCode: toChannelCode(ch.ChannelCode),
			Name:        ch.ChannelName,
			Available:   available,
			ExtraFee:    ch.BaseFee,
		}
		if hasArea {
			v.ExtraFee = ch.BaseFee + area.ExtraFee
			if !available {
				v.Tip = area.UnavailableReason
			}
		} else {
			v.Tip = "该地区暂不支持此渠道"
		}
		resp.Channels = append(resp.Channels, v)
	}
	return resp, nil
}

func (s *ShippingService) ListLockers(ctx context.Context, req *v1.ListLockersRequest) (*v1.ListLockersResponse, error) {
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 50
	}
	offset := 0

	lockers, total, err := s.lockerUc.ListLockers(ctx, req.GetDistrictCode(), req.GetLatitude(), req.GetLongitude(), req.GetSearchRadius(), pageSize, offset)
	if err != nil {
		return nil, err
	}

	resp := &v1.ListLockersResponse{
		Lockers: make([]*v1.Locker, 0, len(lockers)),
	}
	for _, l := range lockers {
		resp.Lockers = append(resp.Lockers, &v1.Locker{
			Id:        l.ID,
			Name:      l.LockerName,
			Address:   l.DetailAddress,
			Latitude:  l.Latitude,
			Longitude: l.Longitude,
			Status:    l.Status,
		})
	}
	if len(lockers) >= pageSize && (offset+len(lockers)) < total {
		resp.NextPageToken = ""
	}
	return resp, nil
}

func (s *ShippingService) ListBoxTypes(ctx context.Context, req *v1.ListBoxTypesRequest) (*v1.ListBoxTypesResponse, error) {
	boxes, err := s.lockerUc.ListBoxTypes(ctx, req.GetLockerId())
	if err != nil {
		return nil, err
	}

	boxTypeMap := make(map[string]*biz.LockerBox)
	for _, b := range boxes {
		boxTypeMap[b.BoxType] = b
	}

	availableCounts := make(map[string]int)
	for _, b := range boxes {
		if b.Status == 1 {
			availableCounts[b.BoxType]++
		}
	}

	resp := &v1.ListBoxTypesResponse{
		BoxTypes: make([]*v1.BoxTypeInfo, 0),
	}
	seen := make(map[string]bool)
	for _, b := range boxes {
		if seen[b.BoxType] {
			continue
		}
		seen[b.BoxType] = true
		minFee := b.BoxFee
		for _, other := range boxes {
			if other.BoxType == b.BoxType && other.BoxFee < minFee {
				minFee = other.BoxFee
			}
		}
		resp.BoxTypes = append(resp.BoxTypes, &v1.BoxTypeInfo{
			BoxType:        toBoxType(b.BoxType),
			Length:         b.BoxLength,
			Width:          b.BoxWidth,
			Height:         b.BoxHeight,
			MaxWeight:      b.MaxWeight,
			BoxFee:         minFee,
			AvailableCount: int32(availableCounts[b.BoxType]),
		})
	}
	return resp, nil
}

func (s *ShippingService) ListServicePoints(ctx context.Context, req *v1.ListServicePointsRequest) (*v1.ListServicePointsResponse, error) {
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 50
	}
	offset := 0

	pointType := int32(0)
	if req.GetPointType() != v1.ServicePointType_SERVICE_POINT_TYPE_UNSPECIFIED {
		pointType = int32(req.GetPointType())
	}

	points, total, err := s.servicePointUc.ListServicePoints(ctx, pointType, req.GetDistrictCode(), req.GetLatitude(), req.GetLongitude(), pageSize, offset)
	if err != nil {
		return nil, err
	}

	resp := &v1.ListServicePointsResponse{
		ServicePoints: make([]*v1.ServicePoint, 0, len(points)),
	}
	for _, p := range points {
		resp.ServicePoints = append(resp.ServicePoints, &v1.ServicePoint{
			Id:        p.ID,
			Name:      p.PointName,
			Address:   p.DetailAddress,
			Latitude:  p.Latitude,
			Longitude: p.Longitude,
			PointType: toServicePointType(p.PointType),
			Status:    p.Status,
		})
	}
	if len(points) >= pageSize && (offset+len(points)) < total {
		resp.NextPageToken = ""
	}
	return resp, nil
}

func (s *ShippingService) ListPickupDates(ctx context.Context, req *v1.ListPickupDatesRequest) (*v1.ListPickupDatesResponse, error) {
	dates, err := s.pickupUc.ListPickupDates(ctx, req.GetDistrictCode())
	if err != nil {
		return nil, err
	}

	resp := &v1.ListPickupDatesResponse{
		PickupDates: make([]*v1.PickupDate, 0, len(dates)),
	}
	for _, d := range dates {
		resp.PickupDates = append(resp.PickupDates, &v1.PickupDate{
			Date:      d,
			Available: true,
		})
	}
	return resp, nil
}

func (s *ShippingService) ListPickupTimeSlots(ctx context.Context, req *v1.ListPickupTimeSlotsRequest) (*v1.ListPickupTimeSlotsResponse, error) {
	slots, err := s.pickupUc.ListPickupTimeSlots(ctx, req.GetDistrictCode(), req.GetPickupDate())
	if err != nil {
		return nil, err
	}

	resp := &v1.ListPickupTimeSlotsResponse{
		TimeSlots: make([]*v1.PickupTimeSlot, 0, len(slots)),
	}
	for _, sl := range slots {
		available := sl.ReservedCount < sl.Capacity
		resp.TimeSlots = append(resp.TimeSlots, &v1.PickupTimeSlot{
			DistrictCode:  sl.DistrictCode,
			PickupDate:    sl.PickupDate,
			SlotCode:      sl.SlotCode,
			StartTime:     sl.StartTime,
			EndTime:       sl.EndTime,
			Capacity:      sl.Capacity,
			ReservedCount: sl.ReservedCount,
			Available:     available,
			Version:       sl.Version,
		})
	}
	return resp, nil
}

func toChannelCode(code string) v1.ChannelCode {
	switch code {
	case "DOOR_PICKUP":
		return v1.ChannelCode_DOOR_PICKUP
	case "LOCKER":
		return v1.ChannelCode_LOCKER
	case "SF_STATION":
		return v1.ChannelCode_SF_STATION
	case "PARTNER_STORE":
		return v1.ChannelCode_PARTNER_STORE
	default:
		return v1.ChannelCode_CHANNEL_UNSPECIFIED
	}
}

func toBoxType(boxType string) v1.BoxType {
	switch boxType {
	case "SMALL":
		return v1.BoxType_SMALL
	case "MEDIUM":
		return v1.BoxType_MEDIUM
	case "LARGE":
		return v1.BoxType_LARGE
	default:
		return v1.BoxType_BOX_TYPE_UNSPECIFIED
	}
}

func toServicePointType(pt int32) v1.ServicePointType {
	switch pt {
	case 1:
		return v1.ServicePointType_POINT_SF_STATION
	case 2:
		return v1.ServicePointType_POINT_PARTNER_STORE
	default:
		return v1.ServicePointType_SERVICE_POINT_TYPE_UNSPECIFIED
	}
}
