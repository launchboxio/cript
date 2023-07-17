package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cript",
		Short: "Container Risk Inspection & Protection Tool",
	}
	logger *zap.SugaredLogger
)

func init() {
	rootCmd.AddCommand(operatorCmd)
	rootCmd.AddCommand(scanCmd)

	baseLogger, _ := zap.NewProduction()
	defer baseLogger.Sync()
	logger = baseLogger.Sugar()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
