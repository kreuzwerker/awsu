package command

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{

	Use:     "whoami",
	Short:   "Info about the currently selected account",
	PreRunE: stsClient.PreRun,
	RunE: func(cmd *cobra.Command, args []string) error {

		res, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})

		if err != nil {
			return err
		}

		fmt.Println(*res.Arn)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
