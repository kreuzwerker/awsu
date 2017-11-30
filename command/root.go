package command

import (
	"fmt"
	"os"

	"github.com/kreuzwerker/awsu/log"
	"github.com/spf13/cobra"
	"github.com/yawn/doubledash"
)

var rootFlags = struct {
	verbose   bool
	workspace string
}{}

var rootCmd = &cobra.Command{
	Use: app,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		if rootFlags.verbose {
			log.Debug = true
		}

		if os.Getenv("AWSU_VERBOSE") != "" {
			log.Debug = true
		}

		return nil

	},
	RunE: func(cmd *cobra.Command, args []string) error {

		sess, err := newSession(rootFlags.workspace)

		if err != nil {
			return err
		}

		if len(doubledash.Xtra) > 0 {
			return sess.Exec(doubledash.Xtra[0], doubledash.Xtra)
		}

		fmt.Println(sess.String())

		return nil

	},
}

func init() {

	os.Args = doubledash.Args

	rootCmd.PersistentFlags().StringVarP(&rootFlags.workspace, "workspace", "w", "", "set the currently used workspace, default to Terraform settings")
	rootCmd.PersistentFlags().BoolVarP(&rootFlags.verbose, "verbose", "v", false, "enable verbose operations")

}
