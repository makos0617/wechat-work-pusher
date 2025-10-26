package main

import (
	"fmt"
	"wechat-work-pusher/cmd"
	"wechat-work-pusher/pkg/config"
)

var (
	version   string
	date      string
	goVersion string
)

func main() {
	info := fmt.Sprintf("***Wechat-Work-Pusher %s***\n***BuildDate %s***\n***%s***\n", version, date, goVersion)
	fmt.Print(info)
	config.LoadConfig()
	cmd.Execute()
}
