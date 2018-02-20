package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/spf13/cobra"
)

type listOutput struct {
	Account string `json:"account,omitempty"`
	Class   string `json:"class"`
	MFA     bool   `json:"mfa"`
}

var listCmd = &cobra.Command{

	Use:   "list",
	Short: "List all configured profiles",
	RunE: func(cmd *cobra.Command, args []string) error {

		reports := make(map[string]listOutput, len(rootFlags.profiles))

		for _, profile := range rootFlags.profiles {

			var report = listOutput{
				MFA: profile.MFASerial != "",
			}

			if profile.RoleARN == "" {
				report.Class = "long-term"
			} else {

				a, err := arn.Parse(profile.RoleARN)

				if err != nil {
					return err
				}

				report.Account = a.AccountID

				if profile.ExternalID != "" {
					report.Class = "external-cross-account"
				} else {
					report.Class = "cross-account"
				}

			}

			reports[profile.Name] = report

		}

		out, err := json.Marshal(reports)

		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, string(out))

		return nil

	},
}

func init() {

	rootCmd.AddCommand(listCmd)

}
