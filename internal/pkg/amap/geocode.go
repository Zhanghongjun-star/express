package amap

import (
	"context"
	"net/url"
)

// Geocode 地址编码结果。
type Geocode struct {
	FormattedAddress StringOrEmpty `json:"formatted_address"`
	Country          StringOrEmpty `json:"country"`
	Province         StringOrEmpty `json:"province"`
	CityCode         string        `json:"citycode"`
	City             StringOrEmpty `json:"city"`
	District         StringOrEmpty `json:"district"`
	Adcode           string        `json:"adcode"`
	Street           StringOrEmpty `json:"street"`
	Number           StringOrEmpty `json:"number"`
	Location         string        `json:"location"`
	Level            StringOrEmpty `json:"level"`
}

// GeoResponse 地理编码返回。
type GeoResponse struct {
	baseResponse
	Count    string    `json:"count"`
	Geocodes []Geocode `json:"geocodes"`
}

// GeocodeOption 定制地理编码请求。
type GeocodeOption func(url.Values)

// GeocodeWithCity 限定城市（城市名或 adcode），提高精度。
func GeocodeWithCity(city string) GeocodeOption {
	return func(v url.Values) { v.Set("city", city) }
}

// GeocodeWithBatch 是否批量（多地址用 "|" 分隔）。
func GeocodeWithBatch(batch bool) GeocodeOption {
	return func(v url.Values) { v.Set("batch", boolToStr(batch)) }
}

// Geocode 将地址文本转换为经纬度坐标。
func (c *Client) Geocode(ctx context.Context, address string, opts ...GeocodeOption) (*GeoResponse, error) {
	q := url.Values{}
	q.Set("address", address)
	for _, o := range opts {
		o(q)
	}
	out := &GeoResponse{}
	if err := c.do(ctx, "/geocode/geo", q, out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddressComponent 逆地理编码的地址组成。
type AddressComponent struct {
	Country  StringOrEmpty `json:"country"`
	Province StringOrEmpty `json:"province"`
	City     StringOrEmpty `json:"city"`
	CityCode string        `json:"citycode"`
	District StringOrEmpty `json:"district"`
	Adcode   string        `json:"adcode"`
	Township StringOrEmpty `json:"township"`
	TownCode string        `json:"towncode"`
	Street   StringOrEmpty `json:"street"`
	Number   StringOrEmpty `json:"number"`
	Location string        `json:"location"`
}

// ReGeoResponse 逆地理编码返回。
type ReGeoResponse struct {
	baseResponse
	ReGeocode struct {
		FormattedAddress string          `json:"formatted_address"`
		AddressComponent AddressComponent `json:"addressComponent"`
	} `json:"regeocode"`
}

// ReGeoOption 定制逆地理编码请求。
type ReGeoOption func(url.Values)

// ReGeoWithRadius 搜索半径（米）。
func ReGeoWithRadius(radius int) ReGeoOption {
	return func(v url.Values) { v.Set("radius", itoa(radius)) }
}

// ReGeoWithPOIType 返回 POI 的类型筛选。
func ReGeoWithPOIType(t string) ReGeoOption {
	return func(v url.Values) { v.Set("poitype", t) }
}

// ReGeoWithExtensions 是否返回周边 POI（base/all）。
func ReGeoWithExtensions(ext string) ReGeoOption {
	return func(v url.Values) { v.Set("extensions", ext) }
}

// ReGeocode 将经纬度坐标转换为地址描述。location 为 "经度,纬度"。
func (c *Client) ReGeocode(ctx context.Context, location string, opts ...ReGeoOption) (*ReGeoResponse, error) {
	q := url.Values{}
	q.Set("location", location)
	for _, o := range opts {
		o(q)
	}
	out := &ReGeoResponse{}
	if err := c.do(ctx, "/geocode/regeo", q, out); err != nil {
		return nil, err
	}
	return out, nil
}
