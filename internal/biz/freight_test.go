package biz

import (
	"testing"
)

func TestIsSameCity(t *testing.T) {
	uc := &FreightUsecase{}

	tests := []struct {
		city1, city2 string
		expected     bool
	}{
		{"深圳市", "深圳市", true},
		{"深圳市", "北京市", false},
		{"广州市", "深圳市", false},
		{"上海市", "杭州市", false},
		{"广州", "深圳", false},
		{"北京市", "北京市", true},
	}

	for _, tt := range tests {
		got := uc.isSameCity(tt.city1, tt.city2)
		if got != tt.expected {
			t.Errorf("isSameCity(%q, %q) = %v, want %v", tt.city1, tt.city2, got, tt.expected)
		}
	}
}

func TestDetermineZone(t *testing.T) {
	uc := &FreightUsecase{}

	tests := []struct {
		name                               string
		sendProvince, sendCity             string
		recvProvince, recvCity             string
		expected                           Zone
	}{
		{"同城", "广东省", "深圳市", "广东省", "深圳市", ZoneSameCity},
		{"同省不同市", "广东省", "广州市", "广东省", "深圳市", ZoneSameProvince},
		{"跨省-广东到北京", "广东省", "深圳市", "北京市", "北京市", ZoneDistant},
		{"同经济圈-广东到广西", "广东省", "广州市", "广西壮族自治区", "南宁市", ZoneSameCircle},
		{"同经济圈-上海到杭州", "上海市", "上海市", "浙江省", "杭州市", ZoneSameCircle},
		{"偏远-广东到新疆", "广东省", "深圳市", "新疆维吾尔自治区", "乌鲁木齐市", ZoneRemote},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uc.determineZone(tt.sendProvince, tt.sendCity, tt.recvProvince, tt.recvCity)
			if got != tt.expected {
				t.Errorf("determineZone = %d, want %d", got, tt.expected)
			}
		})
	}
}
