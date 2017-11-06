package command

import (
	"fmt"
	"os"

	"github.com/kreuzwerker/awsu/yubikey"
	"github.com/spf13/cobra"
)

var registerFlags = struct {
}{}

var registerCmd = &cobra.Command{

	Use:   "register",
	Short: "Initializes an device on AWS and Yubikey",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		sess, err := newSession(rootFlags.workspace)

		if err != nil {
			return err
		}

		arn, err := yubikey.NewMFA(
			sess,
			os.Stderr,
			args[0],
		)

		if err != nil {
			return err
		}

		fmt.Println(*arn)

		return nil

	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
