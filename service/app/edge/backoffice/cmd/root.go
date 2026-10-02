package cmd

import (
	"fmt"
	"go-pkg/conf"
	"go-pkg/hook"
	"go-pkg/log"
	"os"

	"github.com/spf13/cobra"
)

func initService() {
	log.Info(func() string {
		return fmt.Sprintf("fileConfig used: %s", conf.GetConfigStatus())
	})
	hook.Init()
}

func init() {
	cobra.OnInitialize(initService)
}

// RootCmd .
var RootCmd = &cobra.Command{
	Short: "Compile many function of BackofficeService}",
}

// Execute ...
func Execute() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(serveCmd)
	if err := RootCmd.Execute(); err != nil {
		log.Error(func() string {
			return fmt.Sprintf("failed to start program, err: %v", err)
		})
		os.Exit(-1)
	}
}
