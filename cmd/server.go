package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lyyyuna/tonghu-chat/pkg/agent"
	"github.com/lyyyuna/tonghu-chat/pkg/chat"
	"github.com/lyyyuna/tonghu-chat/pkg/config"
	"github.com/lyyyuna/tonghu-chat/pkg/nats"
	"github.com/lyyyuna/tonghu-chat/pkg/redis"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var serverCmd = &cobra.Command{
	Use:   "runchat",
	Short: "Start a basic chat server",
	Run:   runChatServer,
}

var (
	configPath string
)

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	serverCmd.Flags().StringVarP(&configPath, "config", "", "chat.toml", "")
}

func runChatServer(cmd *cobra.Command, args []string) {
	r := gin.Default()

	chatConfig := config.ReadConfig(configPath)

	natsUri := fmt.Sprintf("%v:%v", chatConfig.Nats.Host, chatConfig.Nats.Port)
	nc := nats.NewNatsClient(chatConfig.Nats.ClusterId, chatConfig.Nats.ClientId, natsUri)
	// redisUri := fmt.Sprintf("%v:%v", redisHost, redisPort)
	rc := redis.NewRedisClient(chatConfig.Redis.Host, "", chatConfig.Redis.Port)

	agent.NewWssServer(r, nc, rc)
	chat.NewChatServer(r, rc)
	chatUri := fmt.Sprintf("%v:%v", chatConfig.ChatHost, chatConfig.ChatPort)
	r.Run(chatUri)
}
