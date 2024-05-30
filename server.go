package main

import (
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the Fiber web server",
	Long:  "Start the Fiber web server and serve your application",
}

func executeServer(app func()) {
	serveCmd.Run = func(cmd *cobra.Command, args []string) {
		app()
	}
	rootCmd.AddCommand(serveCmd)
}
