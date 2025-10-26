package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Request 请求参数
type Request struct {
	Method      string
	URL         string
	Headers     map[string]string
	FormData    map[string]string
	JSONData    interface{}
	ContentType string
	Timeout     time.Duration
}

// Response 响应结果
type Response struct {
	Body       string
	StatusCode int
	Error      error
}

// DefaultClient 默认HTTP客户端，带超时设置
var DefaultClient = &http.Client{
	Timeout: 30 * time.Second,
}

// DoRequest 执行HTTP请求
func DoRequest(req Request) Response {
	client := DefaultClient
	if req.Timeout > 0 {
		client = &http.Client{Timeout: req.Timeout}
	}

	var body io.Reader
	contentType := req.ContentType

	// 处理请求体
	if req.JSONData != nil {
		jsonBytes, err := json.Marshal(req.JSONData)
		if err != nil {
			return Response{Error: err}
		}
		body = bytes.NewReader(jsonBytes)
		if contentType == "" {
			contentType = "application/json"
		}
	} else if req.FormData != nil {
		formValues := url.Values{}
		for k, v := range req.FormData {
			formValues.Set(k, v)
		}
		body = strings.NewReader(formValues.Encode())
		if contentType == "" {
			contentType = "application/x-www-form-urlencoded"
		}
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest(req.Method, req.URL, body)
	if err != nil {
		return Response{Error: err}
	}

	// 设置Content-Type
	if contentType != "" {
		httpReq.Header.Set("Content-Type", contentType)
	}

	// 设置自定义头部
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// 执行请求
	resp, err := client.Do(httpReq)
	if err != nil {
		return Response{Error: err}
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{Error: err}
	}

	return Response{
		Body:       string(respBody),
		StatusCode: resp.StatusCode,
		Error:      nil,
	}
}