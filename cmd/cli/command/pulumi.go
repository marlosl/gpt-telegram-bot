package command

import (
	"github.com/marlosl/gpt-telegram-bot/cmd/cli/awsdeploy"

	"github.com/spf13/cobra"
)

var (
	awsCmd = &cobra.Command{
		Use:   "aws",
		Short: "Manage AWS infrastructure",
	}

	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy the GPT Talkerinfrastructure.",
		Run: func(cmd *cobra.Command, args []string) {
			awsdeploy.ExecuteCommand(awsdeploy.AWS_DEPLOY_COMMAND)
		},
	}

	destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the GPT Talker infrastructure.",
		Run: func(cmd *cobra.Command, args []string) {
			awsdeploy.ExecuteCommand(awsdeploy.AWS_DESTROY_COMMAND)
		},
	}
)

func init() {
	awsCmd.AddCommand(deployCmd)
	awsCmd.AddCommand(destroyCmd)

	rootCmd.AddCommand(awsCmd)
}
