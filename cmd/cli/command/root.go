package command

import (
	"os"

	"github.com/marlosl/gpt-telegram-bot/consts"
	"github.com/marlosl/gpt-telegram-bot/utils"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "A CLI utility to manage the GPT Talker",
}

func Execute() error {
	cobra.OnInitialize(utils.InitConfig)

	os.Setenv(consts.AwsRegion, "us-east-1")

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd.Execute()
}
