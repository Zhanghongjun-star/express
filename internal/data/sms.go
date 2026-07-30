package data

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// SendSMS 通过互亿无线发送短信验证码。
// content 如 "您的验证码是：123456。请不要把验证码泄露给其他人。"
func SendSMS(ctx context.Context, mobile, content string) error {
	if SMSURL == "" || SMSAcc == "" || SMSPwd == "" {
		return fmt.Errorf("sms not configured")
	}
	v := url.Values{}
	v.Set("account", SMSAcc)
	v.Set("password", SMSPwd)
	v.Set("mobile", mobile)
	v.Set("content", content)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, SMSURL, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "SubmitResult") {
		return fmt.Errorf("sms send failed: %s", string(body))
	}
	return nil
}
