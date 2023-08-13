package command

import (
	"fmt"

	"github.com/marlosl/gpt-telegram-bot/services/telegram"

	"github.com/spf13/cobra"
)

var (
	telegramWebhookCmd = &cobra.Command{
		Use:   "telegram-webhook",
		Short: "Set Telegram Webhook.",
	}

	setTextWebhookCmd = &cobra.Command{
		Use:   "text",
		Short: "Set Text Webhook.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide a webhook url.")
				cmd.Help()
				return
			}

			service := telegram.NewTextService()
			service.SetWebhook(args[0], args[1])
		},
	}

	setImageWebhookCmd = &cobra.Command{
		Use:   "image",
		Short: "Set Image Webhook.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide a webhook url.")
				cmd.Help()
				return
			}
			service := telegram.NewImageService()
			service.SetWebhook(args[0], args[1])
		},
	}
)

func init() {
	telegramWebhookCmd.AddCommand(setTextWebhookCmd)
	telegramWebhookCmd.AddCommand(setImageWebhookCmd)

	rootCmd.AddCommand(telegramWebhookCmd)
}
