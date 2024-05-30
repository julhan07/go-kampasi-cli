package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fiber-cli",
	Short: "Fiber CLI is a tool for managing your Fiber application",
	Long: `Fiber CLI is a command line tool that helps you manage your
Fiber application efficiently. You can use it to generate models,
repositories, services, interfaces, and handlers.`,
}

func Execute(app func()) {
	generatorExecute()
	executeServer(app)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
