package command

import (
	"fmt"

	"github.com/kreuzwerker/awsu"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{

	Use:   "whoami",
	Short: "Info about the currently selected account",
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := awsu.NewClient()

		if err != nil {
			return err
		}

		// TODO: use flags to determine what to return

		res, err := client.CallerIdentity()

		if err != nil {
			return err
		}

		fmt.Println(res.ARN)

		return nil

	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
