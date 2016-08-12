package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{

	Use:   "ls",
	Short: "List all profile",
	RunE: func(cmd *cobra.Command, args []string) error {

		keys := config.config.SectionStrings()
		sort.Strings(keys)

		fmt.Println(strings.Join(keys, "\n"))

		return nil

	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
