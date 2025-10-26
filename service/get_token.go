package service

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
	"wechat-work-pusher/constant"
	"wechat-work-pusher/pkg/config"
	"wechat-work-pusher/pkg/httpclient"

	"github.com/tidwall/gjson"
)

type getToken struct {
	cropId string
	secret string
	once   sync.Once
}

type respToken struct {
	val    string
	expire time.Time
	lock   sync.RWMutex
}

var getTokenParams getToken
var respTokenParams respToken

// 并发刷新合并控制
var refreshMu sync.Mutex
var refreshing bool
var refreshCond *sync.Cond

func (th *respToken) set(val string) {
	th.lock.Lock()
	defer th.lock.Unlock()
	th.val = val
	th.expire = time.Now().Add(time.Minute * constant.RespTokenExpireMin)
}

func (th *respToken) isValid() bool {
	th.lock.RLock()
	defer th.lock.RUnlock()
	return !th.expire.IsZero() && time.Now().Before(th.expire) && th.val != ""
}

func (th *respToken) get() string {
	th.lock.RLock()
	defer th.lock.RUnlock()
	return th.val
}

// 带有限重试的获取企业微信token，并发刷新合并
func GetTokenFromWechat() (string, error) {
	if respTokenParams.isValid() {
		return respTokenParams.get(), nil
	}
	getTokenParams.once.Do(func() {
		getTokenParams.cropId = config.GetString(constant.ConfigKeyWorkCorpId)
		getTokenParams.secret = config.GetString(constant.ConfigKeyWorkCorpSecret)
	})
	refreshMu.Lock()
	if refreshCond == nil {
		refreshCond = sync.NewCond(&refreshMu)
	}
	if refreshing {
		// 等待其他协程刷新完成
		for refreshing {
			refreshCond.Wait()
		}
		val := respTokenParams.get()
		refreshMu.Unlock()
		if val == "" {
			return "", errors.New("token refresh failed")
		}
		return val, nil
	}
	// 当前协程负责刷新
	refreshing = true
	refreshMu.Unlock()

	var lastErr error
	// 简单退避重试：3次，200ms, 400ms, 800ms
	for attempt := 0; attempt < 3; attempt++ {
		resp := httpclient.DoRequest(httpclient.Request{
			Method: http.MethodGet,
			URL:    fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", getTokenParams.cropId, getTokenParams.secret),
		})
		code := resp.StatusCode
		if resp.Error != nil {
			lastErr = resp.Error
		} else if code == http.StatusOK {
			if gjson.Get(resp.Body, "errcode").Int() == 0 {
				token := gjson.Get(resp.Body, "access_token").String()
				if token != "" {
					respTokenParams.set(token)
					lastErr = nil
					break
				} else {
					lastErr = errors.New("empty access_token")
				}
			} else {
				ec := gjson.Get(resp.Body, "errcode").Int()
				lastErr = fmt.Errorf("wechat errcode:%d", ec)
			}
		} else {
			lastErr = fmt.Errorf("http code:%d", code)
		}
		time.Sleep(time.Duration(200*(1<<attempt)) * time.Millisecond)
	}

	refreshMu.Lock()
	refreshing = false
	refreshCond.Broadcast()
	val := respTokenParams.get()
	refreshMu.Unlock()
	if val == "" {
		if lastErr == nil {
			lastErr = errors.New("token refresh unknown error")
		}
		return "", lastErr
	}
	return val, nil
}
