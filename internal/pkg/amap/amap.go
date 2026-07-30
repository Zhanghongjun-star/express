// Package amap 封装高德地图开放平台 Web 服务 API（https://lbs.amap.com/）。
// 支持地理/逆地理编码、路径规划（驾车/步行/骑行/公交）、IP 定位、天气、POI 搜索与距离测量。
//
// 使用示例：
//
//	client := amap.NewClient("你的高德Key")
//	resp, err := client.Geocode(ctx, "北京市朝阳区望京")
//	if err != nil {
//		// 处理错误（业务错误为 *amap.APIError）
//	}
//	fmt.Println(resp.Geocodes[0].Location)
package amap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultBaseURL = "https://restapi.amap.com/v3"
	defaultTimeout = 5 * time.Second
)

// Client 是高德地图 API 的 HTTP 客户端。
type Client struct {
	key        string
	baseURL    string
	httpClient *http.Client
}

// Option 用于定制 Client。
type Option func(*Client)

// WithBaseURL 自定义 API 基地址（默认 https://restapi.amap.com/v3）。
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = baseURL }
}

// WithHTTPClient 自定义底层 http.Client（用于注入超时、拦截器等）。
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// NewClient 使用高德 Key 创建一个 API 客户端。
func NewClient(key string, opts ...Option) *Client {
	c := &Client{
		key:        key,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// APIError 表示高德返回的业务错误（status=0）。
type APIError struct {
	Infocode string
	Info     string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("amap: infocode=%s info=%s", e.Infocode, e.Info)
}

// StringOrEmpty 兼容高德把空字段返回为空数组而非空字符串的情况，
// 反序列化时若为数组则取首个元素（或留空）。
type StringOrEmpty string

func (s *StringOrEmpty) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		*s = StringOrEmpty(str)
		return nil
	}
	var arr []string
	if err := json.Unmarshal(b, &arr); err == nil {
		if len(arr) > 0 {
			*s = StringOrEmpty(arr[0])
		}
		return nil
	}
	*s = ""
	return nil
}

// MaybeEmpty 兼容高德把空字段返回为空数组而非空对象的情况。
// 解码后若原值为数组（空），Data 为 nil；否则为指向解析结果的非 nil 指针。
type MaybeEmpty[T any] struct {
	Data *T
}

func (m *MaybeEmpty[T]) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '[' {
		m.Data = nil
		return nil
	}
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	m.Data = &v
	return nil
}

// baseResponse 是所有响应的公共字段。
type baseResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
}

func (b baseResponse) status() string   { return b.Status }
func (b baseResponse) info() string     { return b.Info }
func (b baseResponse) infocode() string { return b.Infocode }

type response interface {
	status() string
	info() string
	infocode() string
}

func (c *Client) do(ctx context.Context, path string, query url.Values, out any) error {
	query.Set("key", c.key)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path+"?"+query.Encode(), nil)
	if err != nil {
		return fmt.Errorf("amap: build request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("amap: http: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("amap: unexpected http status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("amap: read body: %w", err)
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("amap: decode: %w", err)
	}
	if r, ok := out.(response); ok && r.status() != "1" {
		return &APIError{Infocode: r.infocode(), Info: r.info()}
	}
	return nil
}

func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func itoa(i int) string { return strconv.Itoa(i) }
