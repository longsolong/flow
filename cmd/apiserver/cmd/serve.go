package cmd

import (
	"os"
	"github.com/spf13/viper"
	"net/http"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/infra"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		viper.AutomaticEnv()
		run()
	},
}

func run() {
	logger := infra.InitializeLogger()

	srv, err := InitializeServer()
	if err != nil {
		logger.SugarZap.Errorf("InitializeServer error %v", err)
		os.Exit(1)
	}
	bootstrap(srv)
	logger.SugarZap.Info("Starting Application")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.SugarZap.Fatalf("listen: %s\n", err)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
