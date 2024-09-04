package console

import (
	"os"

	"github.com/kodinggo/payment-service-gb1/internal/config"

	"github.com/spf13/cobra"
)

var rootCMd = &cobra.Command{
	Use:   "console",
	Short: "Console commands",
	Long:  "Console commands",
}

func Execute() {
	err := rootCMd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCMd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	config.SetupLogger()
}
