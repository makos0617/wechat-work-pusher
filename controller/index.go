package controller

import (
	"wechat-work-pusher/middleware"
	"wechat-work-pusher/pkg/httpserver"
	"wechat-work-pusher/service"
)

func Init(router *httpserver.Router) {
	r := router.Use(middleware.AuthToken())
	r.Post("/msg", msg)
	r.Post("/card", card)
}

func msg(ctx *httpserver.Context) {
	to := ctx.Request.FormValue("to")
	content := ctx.Request.FormValue("content")
	
	if err := service.SendMsg(to, content); err != nil {
		ctx.Json(httpserver.RestRet{
			Result:  httpserver.ResultErr,
			Message: httpserver.String{String: err.Error(), Valid: true},
		})
		return
	}
	ctx.JsonSuccess("ok")
}

func card(ctx *httpserver.Context) {
	to := ctx.Request.FormValue("to")
	title := ctx.Request.FormValue("title")
	description := ctx.Request.FormValue("description")
	url := ctx.Request.FormValue("url")
	
	if err := service.SendCardMsg(to, title, description, url); err != nil {
		ctx.Json(httpserver.RestRet{
			Result:  httpserver.ResultErr,
			Message: httpserver.String{String: err.Error(), Valid: true},
		})
		return
	}
	ctx.JsonSuccess("ok")
}
