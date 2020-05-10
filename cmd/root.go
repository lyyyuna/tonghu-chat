package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "tonghu-chat",
	Short: "Start a chat server/client",
	Long:  "",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		zap.L().Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	rootCmd.AddCommand(serverCmd)
}
