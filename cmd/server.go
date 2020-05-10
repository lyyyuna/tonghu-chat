package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "basic",
	Short: "Start a basic http proxy without mitm",
	Run:   runChatServer,
}

func runChatServer(cmd *cobra.Command, args []string) {
	r := gin.Default()
	r.Run()
}
