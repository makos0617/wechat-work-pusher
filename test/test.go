package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Push2Wechat(to, msg string) error {
	form := url.Values{}
	form.Set("to", to)
	form.Set("content", msg)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:9000/wechat-work-pusher/msg", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoicm9vdCIsInVwZGF0ZVRpbWUiOjE2MzcxNTk1OTZ9.ZdzKwwiwyxNDiOpAXyoMwKiyffgWC1sLsgdTD3wPqIw"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status: %d", resp.StatusCode)
	}
	return nil
}

// 与服务端签名一致 (to,title,description,url)
func Push2WechatCard(to, title, description, link string) error {
	form := url.Values{}
	form.Set("to", to)
	form.Set("title", title)
	form.Set("description", description)
	form.Set("url", link)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:9000/wechat-work-pusher/card", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoicm9vdCIsInVwZGF0ZVRpbWUiOjE2MzcxNTk1OTZ9.ZdzKwwiwyxNDiOpAXyoMwKiyffgWC1sLsgdTD3wPqIw"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status: %d", resp.StatusCode)
	}
	return nil
}

func main() {
	date := time.Now().Format("2006-01-02")
	title := "测试卡片"
	desc := fmt.Sprintf("<div class=\"gray\">%s</div><div class=\"normal\">这是一条测试卡片消息。</div><div class=\"highlight\">点击卡片查看更多详情。</div>", date)
	link := "https://www.baidu.com"

	if err := Push2WechatCard("@all", title, desc, link); err != nil {
		fmt.Println("Push2WechatCard error:", err)
	}
}
