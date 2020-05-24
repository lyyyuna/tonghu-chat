package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var serverCmd = &cobra.Command{
	Use:   "basic",
	Short: "Start a basic http proxy without mitm",
	Run:   runChatServer,
}

var (
	natsClusterId string
	natsClientId  string
	natsHost      string
	natsPort      int
	redisHost     string
	redisPort     int
	chatPort      int
	chatHost      string
)

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	serverCmd.Flags().StringVarP(&natsClusterId, "natsclusterid", "", "0.0.0.0", "")
	serverCmd.Flags().StringVarP(&natsClientId, "natsclientid", "", "", "")
	serverCmd.Flags().StringVarP(&natsHost, "natshost", "", "127.0.0.1", "")
	serverCmd.Flags().IntVarP(&natsPort, "natsport", "", 1, "")
	serverCmd.Flags().StringVarP(&redisHost, "redishost", "", "127.0.0.1", "")
	serverCmd.Flags().IntVarP(&redisPort, "redisport", "", 6380, "")
	serverCmd.Flags().IntVarP(&chatPort, "chatport", "", 8080, "")
	serverCmd.Flags().StringVarP(&chatHost, "chathost", "", "0.0.0.0", "")
}

func runChatServer(cmd *cobra.Command, args []string) {
	r := gin.Default()
	r.Run()
}
