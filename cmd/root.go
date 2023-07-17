package cmd

import (
  "github.com/launchboxio/cript/internal/config"
  "github.com/spf13/cobra"
  "go.uber.org/zap"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cript",
		Short: "Container Risk Inspection & Protection Tool",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			configPath, _ := cmd.Flags().GetString("config")
			LoadConfig(configPath)
		},
	}
	logger *zap.SugaredLogger
	conf   *config.Config
)

func init() {
	rootCmd.PersistentFlags().String("config", "/etc/cript", "Path to a configuration file")
	rootCmd.AddCommand(operatorCmd)
	rootCmd.AddCommand(scanCmd)

	baseLogger, _ := zap.NewProduction()
	defer baseLogger.Sync()
	logger = baseLogger.Sugar()

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func LoadConfig(configPath string) {
	var err error
	conf, err = config.Load(configPath)
	if err != nil {
		logger.Fatal(err)
	}
}
