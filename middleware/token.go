package middleware

import (
	"strings"
	"wechat-work-pusher/constant"
	"wechat-work-pusher/pkg/config"
	"wechat-work-pusher/pkg/httpserver"
)

var token string

func AuthToken() httpserver.Middleware {
	if token == "" {
		token = config.GetString(constant.ConfigKeyToken)
	}
	return httpserver.AuthMiddleware(func(ctx *httpserver.Context) bool {
		// 从 Authorization: Bearer <token> 读取鉴权令牌
		authHeader := ctx.Request.Header.Get("Authorization")
		var getToken string
		if strings.HasPrefix(authHeader, "Bearer ") {
			getToken = strings.TrimSpace(authHeader[len("Bearer "):])
		}
		if getToken == "" || getToken != token {
			ctx.Json(httpserver.RestRet{
				Result: httpserver.ResultAuthErr,
				Message: httpserver.String{
					String: "鉴权失败，缺少或非法Authorization",
					Valid:  true,
				},
			})
			ctx.StopExecution()
			return false
		}
		return true
	})
}
