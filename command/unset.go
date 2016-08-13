package command

import (
	"fmt"

	"github.com/kreuzwerker/awsu"
	"github.com/spf13/cobra"
)

var unsetCmd = &cobra.Command{

	Use:   "unset",
	Short: "Unset shell variables",
	RunE: func(cmd *cobra.Command, args []string) error {

		for _, e := range awsu.AllKeys {
			fmt.Printf("unset %s\n", e)
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(unsetCmd)
}
