package amap

import (
	"context"
	"net/url"
)

// DirectionStep 是驾车/步行/骑行路径中的一个导航步骤。
type DirectionStep struct {
	Instruction     StringOrEmpty `json:"instruction"`
	Road            StringOrEmpty `json:"road"`
	Distance        string        `json:"distance"`
	Orientation     StringOrEmpty `json:"orientation"`
	Duration        string        `json:"duration"`
	Polyline        string        `json:"polyline"`
	Action          StringOrEmpty `json:"action"`
	Navi            StringOrEmpty `json:"navi"`
	AssistantAction StringOrEmpty `json:"assistant_action"`
}

// DirectionPath 一条完整路径方案。
type DirectionPath struct {
	Path         string          `json:"path"`
	Distance     string          `json:"distance"`
	Duration     string          `json:"duration"`
	TollDistance string          `json:"toll_distance"`
	Steps        []DirectionStep `json:"steps"`
}

// DirectionResponse 驾车/步行/骑行的返回结构（三者一致）。
type DirectionResponse struct {
	baseResponse
	Route struct {
		Origin       string          `json:"origin"`
		Destination  string          `json:"destination"`
		Distance     string          `json:"distance"`
		TollDistance string          `json:"toll_distance"`
		Paths        []DirectionPath `json:"paths"`
	} `json:"route"`
}

// DirectionOption 定制路径规划请求。
type DirectionOption func(url.Values)

// DirectionWithCity 设置城市（城市内公交/限定范围）。
func DirectionWithCity(city string) DirectionOption {
	return func(v url.Values) { v.Set("city", city) }
}

// DirectionWithExtensions base（概要）或 all（含每一步 polyline）。
func DirectionWithExtensions(ext string) DirectionOption {
	return func(v url.Values) { v.Set("extensions", ext) }
}

// DirectionWithStrategy 驾车策略（0-2 等）；步行/骑行忽略。
func DirectionWithStrategy(strategy string) DirectionOption {
	return func(v url.Values) { v.Set("strategy", strategy) }
}

func directionQuery(origin, destination string, opts []DirectionOption) url.Values {
	q := url.Values{}
	q.Set("origin", origin)
	q.Set("destination", destination)
	for _, o := range opts {
		o(q)
	}
	return q
}

// DirectionDriving 驾车路径规划。origin/destination 为 "经度,纬度"。
func (c *Client) DirectionDriving(ctx context.Context, origin, destination string, opts ...DirectionOption) (*DirectionResponse, error) {
	out := &DirectionResponse{}
	if err := c.do(ctx, "/direction/driving", directionQuery(origin, destination, opts), out); err != nil {
		return nil, err
	}
	return out, nil
}

// DirectionWalking 步行路径规划。
func (c *Client) DirectionWalking(ctx context.Context, origin, destination string, opts ...DirectionOption) (*DirectionResponse, error) {
	out := &DirectionResponse{}
	if err := c.do(ctx, "/direction/walking", directionQuery(origin, destination, opts), out); err != nil {
		return nil, err
	}
	return out, nil
}

// DirectionBicycling 骑行路径规划。
func (c *Client) DirectionBicycling(ctx context.Context, origin, destination string, opts ...DirectionOption) (*DirectionResponse, error) {
	out := &DirectionResponse{}
	if err := c.do(ctx, "/direction/bicycling", directionQuery(origin, destination, opts), out); err != nil {
		return nil, err
	}
	return out, nil
}

// TransitBusLine 公交/地铁线路信息。
type TransitBusLine struct {
	Name          string `json:"name"`
	Distance      string `json:"distance"`
	Duration      string `json:"duration"`
	DepartureStop struct {
		Name string `json:"name"`
	} `json:"departure_stop"`
	ArrivalStop struct {
		Name string `json:"name"`
	} `json:"arrival_stop"`
	ViaStops []struct {
		Name string `json:"name"`
	} `json:"via_stops"`
}

// TransitWalking 公交换乘中的步行段。
type TransitWalking struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Distance    string `json:"distance"`
	Duration    string `json:"duration"`
	Steps       []struct {
		Instruction     StringOrEmpty `json:"instruction"`
		Road            StringOrEmpty `json:"road"`
		Distance        string        `json:"distance"`
		Action          StringOrEmpty `json:"action"`
		AssistantAction StringOrEmpty `json:"assistant_action"`
	} `json:"steps"`
}

// TransitRailway 火车/高铁段信息。
type TransitRailway struct {
	Name     string `json:"name"`
	Trip     string `json:"trip"`
	Distance string `json:"distance"`
	Duration string `json:"duration"`
}

// TransitBusGroup 公交段内的线路集合。
type TransitBusGroup struct {
	BusLines []TransitBusLine `json:"buslines"`
}

// TransitNode 出入口节点（名称）。
type TransitNode struct {
	Name string `json:"name"`
}

// TransitSegment 公交方案中的一个换乘段（步行+公交/地铁/火车）。
// 高德在缺失某段时会将该字段返回为空数组而非空对象，故相关字段用 MaybeEmpty 兼容。
type TransitSegment struct {
	Walking  MaybeEmpty[TransitWalking] `json:"walking"`
	Bus      MaybeEmpty[TransitBusGroup] `json:"bus"`
	Entrance MaybeEmpty[TransitNode]    `json:"entrance"`
	Exit     MaybeEmpty[TransitNode]    `json:"exit"`
	Railway  MaybeEmpty[TransitRailway] `json:"railway"`
}

// Transit 一条公交换乘方案。
type Transit struct {
	Duration       string           `json:"duration"`
	WalkingDistance StringOrEmpty   `json:"walking_distance"`
	Segments       []TransitSegment `json:"segments"`
}

// TransitResponse 公交路径规划的返回结构。
type TransitResponse struct {
	baseResponse
	Route struct {
		Origin      string    `json:"origin"`
		Destination string    `json:"destination"`
		Distance    string    `json:"distance"`
		Transits    []Transit `json:"transits"`
	} `json:"route"`
}

// TransitOption 定制公交路径规划请求。
type TransitOption func(url.Values)

// TransitWithStrategy 公交换乘策略（0-2）。
func TransitWithStrategy(strategy string) TransitOption {
	return func(v url.Values) { v.Set("strategy", strategy) }
}

// TransitWithNightFlag 是否包含夜班车（"0"/"1"）。
func TransitWithNightFlag(flag string) TransitOption {
	return func(v url.Values) { v.Set("nightflag", flag) }
}

// DirectionTransit 公交（含火车/地铁）路径规划，跨城须传 city/cityd。
// city/cityd 为起点/终点城市名或 adcode。
func (c *Client) DirectionTransit(ctx context.Context, origin, destination, city, cityd string, opts ...TransitOption) (*TransitResponse, error) {
	q := url.Values{}
	q.Set("origin", origin)
	q.Set("destination", destination)
	q.Set("city", city)
	q.Set("cityd", cityd)
	for _, o := range opts {
		o(q)
	}
	out := &TransitResponse{}
	if err := c.do(ctx, "/direction/transit/integrated", q, out); err != nil {
		return nil, err
	}
	return out, nil
}
