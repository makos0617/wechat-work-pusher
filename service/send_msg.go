package service

import (
	"fmt"
	"net/http"
	"sync"
	"wechat-work-pusher/constant"
	"wechat-work-pusher/pkg/config"
	"wechat-work-pusher/pkg/httpclient"
	"wechat-work-pusher/service/model"

	"github.com/tidwall/gjson"
)

type agent struct {
	val  string
	once sync.Once
}

var agentParams agent

func SendMsg(to, content string) error {
	agentParams.once.Do(func() {
		agentParams.val = config.GetString(constant.ConfigKeyWorkAgentId)
	})
	msg := &model.TextMessage{}
	msg.AgentId = agentParams.val
	msg.MsgType = constant.MessageText
	msg.Text.Content = content
	if to != "" {
		msg.ToUser = to
	} else {
		msg.ToUser = config.GetString(constant.ConfigKeyDefaultReceiver)
	}
	// 获取 token，增加错误返回
	token, err := GetTokenFromWechat()
	if err != nil {
		return fmt.Errorf("get token failed: %w", err)
	}
	resp := httpclient.DoRequest(httpclient.Request{
		Method:      http.MethodPost,
		URL:         fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token),
		ContentType: "application/json",
		JSONData:    msg,
	})
	if resp.Error != nil {
		return fmt.Errorf("send text to wechat req err: %w", resp.Error)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send text to wechat req err, code:%d", resp.StatusCode)
	}
	if ec := gjson.Get(resp.Body, "errcode").Int(); ec != 0 {
		return fmt.Errorf("send text to wechat resp error, code:%d", ec)
	}
	return nil
}

func SendCardMsg(to, title, des, url string) error {
	agentParams.once.Do(func() {
		agentParams.val = config.GetString(constant.ConfigKeyWorkAgentId)
	})
	msg := &model.TextCardMessage{}
	msg.AgentId = agentParams.val
	msg.MsgType = constant.MessageCard
	msg.TextCard.Title = title
	msg.TextCard.Description = des
	msg.TextCard.Url = url
	msg.TextCard.Detail = "详情"
	if to != "" {
		msg.ToUser = to
	} else {
		msg.ToUser = config.GetString(constant.ConfigKeyDefaultReceiver)
	}
	// 获取 token，增加错误返回
	token, err := GetTokenFromWechat()
	if err != nil {
		return fmt.Errorf("get token failed: %w", err)
	}
	resp := httpclient.DoRequest(httpclient.Request{
		Method:      http.MethodPost,
		URL:         fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token),
		ContentType: "application/json",
		JSONData:    msg,
	})
	if resp.Error != nil {
		return fmt.Errorf("send textcard to wechat req err: %w", resp.Error)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send textcard to wechat req err, code:%d", resp.StatusCode)
	}
	if ec := gjson.Get(resp.Body, "errcode").Int(); ec != 0 {
		return fmt.Errorf("send textcard to wechat resp error, code:%d", ec)
	}
	return nil
}
