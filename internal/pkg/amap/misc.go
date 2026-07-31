package amap

import (
	"context"
	"net/url"
	"strconv"
)

// IPResponse IP 定位返回。
type IPResponse struct {
	baseResponse
	IP        string        `json:"ip"`
	Province  StringOrEmpty `json:"province"`
	City      StringOrEmpty `json:"city"`
	Adcode    StringOrEmpty `json:"adcode"`
	Rectangle StringOrEmpty `json:"rectangle"`
}

// IP 根据 IP 定位所在城市。ip 为空时取请求端公网 IP。
func (c *Client) IP(ctx context.Context, ip string) (*IPResponse, error) {
	q := url.Values{}
	if ip != "" {
		q.Set("ip", ip)
	}
	out := &IPResponse{}
	if err := c.do(ctx, "/ip", q, out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeatherLive 实时天气。
type WeatherLive struct {
	Province      StringOrEmpty `json:"province"`
	City          StringOrEmpty `json:"city"`
	Adcode        string `json:"adcode"`
	Weather       string `json:"weather"`
	Temperature   string `json:"temperature"`
	WindDirection string `json:"winddirection"`
	WindPower     string `json:"windpower"`
	Humidity      string `json:"humidity"`
	ReportTime    string `json:"reporttime"`
}

// WeatherCast 天气预报（单日）。
type WeatherCast struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	DayWeather   string `json:"dayweather"`
	NightWeather string `json:"nightweather"`
	DayTemp      string `json:"daytemp"`
	NightTemp    string `json:"nighttemp"`
	DayWind      string `json:"daywind"`
	NightWind    string `json:"nightwind"`
	DayPower     string `json:"daypower"`
	NightPower   string `json:"nightpower"`
}

// WeatherResponse 天气查询返回（实时或预报）。
type WeatherResponse struct {
	baseResponse
	Lives    []WeatherLive `json:"lives"`
	Forecasts []struct {
		City       StringOrEmpty `json:"city"`
		Adcode     string        `json:"adcode"`
		Province   StringOrEmpty `json:"province"`
		ReportTime string        `json:"reporttime"`
		Casts      []WeatherCast `json:"casts"`
	} `json:"forecasts"`
}

// Weather 查询天气。city 为城市 adcode 或城市名；extensions="base" 实时、"all" 预报。
func (c *Client) Weather(ctx context.Context, city, extensions string) (*WeatherResponse, error) {
	q := url.Values{}
	q.Set("city", city)
	q.Set("extensions", extensions)
	out := &WeatherResponse{}
	if err := c.do(ctx, "/weather/weatherInfo", q, out); err != nil {
		return nil, err
	}
	return out, nil
}

// POI 兴趣点。
type POI struct {
	ID           string        `json:"id"`
	Name         StringOrEmpty `json:"name"`
	Address      StringOrEmpty `json:"address"`
	Location     string        `json:"location"`
	CityCode     string        `json:"citycode"`
	CityName     StringOrEmpty `json:"cityname"`
	Adcode       string        `json:"adcode"`
	ProvinceName StringOrEmpty `json:"provincename"`
	Tel          StringOrEmpty `json:"tel"`
	Type         StringOrEmpty `json:"type"`
	TypeCode     StringOrEmpty `json:"typecode"`
}

// PlaceResponse POI 搜索返回。
type PlaceResponse struct {
	baseResponse
	Count string `json:"count"`
	POIs  []POI  `json:"pois"`
}

// PlaceOption 定制 POI 搜索请求。
type PlaceOption func(url.Values)

// PlaceWithCity 限定城市（城市名或 adcode）。
func PlaceWithCity(city string) PlaceOption {
	return func(v url.Values) { v.Set("city", city) }
}

// PlaceWithPage 分页（page 从 1 开始，offset 每页条数）。
func PlaceWithPage(page, offset int) PlaceOption {
	return func(v url.Values) {
		v.Set("page", strconv.Itoa(page))
		v.Set("offset", strconv.Itoa(offset))
	}
}

// PlaceText 关键字搜索 POI。keywords 为关键词，city 可空。
func (c *Client) PlaceText(ctx context.Context, keywords string, opts ...PlaceOption) (*PlaceResponse, error) {
	q := url.Values{}
	q.Set("keywords", keywords)
	for _, o := range opts {
		o(q)
	}
	out := &PlaceResponse{}
	if err := c.do(ctx, "/place/text", q, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PlaceAround 周边搜索 POI。location 为 "经度,纬度"。
func (c *Client) PlaceAround(ctx context.Context, location, keywords string, opts ...PlaceOption) (*PlaceResponse, error) {
	q := url.Values{}
	q.Set("location", location)
	if keywords != "" {
		q.Set("keywords", keywords)
	}
	for _, o := range opts {
		o(q)
	}
	out := &PlaceResponse{}
	if err := c.do(ctx, "/place/around", q, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DistanceResponse 距离测量返回。
type DistanceResponse struct {
	baseResponse
	Results []struct {
		Origin      string `json:"origin"`
		Destination string `json:"destination"`
		Distance    string `json:"distance"`
		Duration    string `json:"duration"`
		Status      string `json:"status"`
	} `json:"results"`
}

// Distance 测量距离。origins 可为多个，用 ";" 分隔（每个 "经度,纬度"）。
// typ: 1=驾车, 2=步行, 3=骑行。
func (c *Client) Distance(ctx context.Context, origins, destination string, typ int) (*DistanceResponse, error) {
	q := url.Values{}
	q.Set("origins", origins)
	q.Set("destination", destination)
	q.Set("type", strconv.Itoa(typ))
	out := &DistanceResponse{}
	if err := c.do(ctx, "/distance", q, out); err != nil {
		return nil, err
	}
	return out, nil
}
