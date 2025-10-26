package cmd

import (
	"fmt"
	"wechat-work-pusher/constant"
	"wechat-work-pusher/controller"
	"wechat-work-pusher/pkg/config"
	"wechat-work-pusher/pkg/httpserver"

	"github.com/spf13/cobra"
)

func init() {
	DefFlags(rootCmd)
}

var rootCmd = &cobra.Command{
	Use: "Wechat-Work-Pusher",
	Run: func(cmd *cobra.Command, args []string) {
		// 覆盖配置中的值（如果提供了命令行参数）
		if v, _ := cmd.Flags().GetString(constant.ConfigKeyWorkCorpId); v != "" {
			config.SetString(constant.ConfigKeyWorkCorpId, v)
		}
		if v, _ := cmd.Flags().GetString(constant.ConfigKeyWorkAgentId); v != "" {
			config.SetString(constant.ConfigKeyWorkAgentId, v)
		}
		if v, _ := cmd.Flags().GetString(constant.ConfigKeyWorkCorpSecret); v != "" {
			config.SetString(constant.ConfigKeyWorkCorpSecret, v)
		}
		if v, _ := cmd.Flags().GetString(constant.ConfigKeyDefaultReceiver); v != "" {
			config.SetString(constant.ConfigKeyDefaultReceiver, v)
		}
		if v, _ := cmd.Flags().GetString(constant.ConfigKeyToken); v != "" {
			config.SetString(constant.ConfigKeyToken, v)
		}

		// 启动 HTTP 服务器
		cfg := config.GetConfig()
		srv := httpserver.NewServer(cfg.Rest.Port)
		srv.AddRoutes(func(r *httpserver.Router) {
			group := r.Group(cfg.Rest.Base)
			controller.Init(group)
		})
		if err := srv.Run(); err != nil {
			fmt.Println(err)
		}
	},
}

func Execute() {
	_ = rootCmd.Execute()
}

func DefFlags(cmd *cobra.Command) {
	cmd.Flags().String(constant.ConfigKeyWorkCorpId, "", "企业ID")
	cmd.Flags().String(constant.ConfigKeyWorkAgentId, "", "应用AgentID")
	cmd.Flags().String(constant.ConfigKeyWorkCorpSecret, "", "应用SecretID")
	cmd.Flags().String(constant.ConfigKeyDefaultReceiver, "", "默认被推送ID")
	cmd.Flags().String(constant.ConfigKeyToken, "", "接口调用Token")
}
