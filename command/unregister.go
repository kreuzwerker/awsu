package command

import (
	"github.com/spf13/cobra"

	"github.com/gesellix/awsu/source/yubikey"
	"github.com/gesellix/awsu/strategy"
	"github.com/gesellix/awsu/target/mfa"
)

var unregisterCmd = &cobra.Command{

	Use:   "unregister :username",
	Short: "Removes a device on AWS and Yubikey",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		username := args[0]

		creds, err := strategy.Apply(&conf)

		if err != nil {
			return err
		}

		source, err := yubikey.New()

		if err != nil {
			return err
		}

		target, err := mfa.New(creds.NewSession(), source)

		if err != nil {
			return err
		}

		return target.Delete(username)

	},
}

func init() {
	rootCmd.AddCommand(unregisterCmd)
}
