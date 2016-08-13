package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Version struct {
	Build, Version string
}

func (v Version) String() string {
	return fmt.Sprintf("%s version %s (%s)", app, v.Version, v.Build)
}

var versionCmd = &cobra.Command{

	Use:   "version",
	Short: fmt.Sprintf("Print the version number of %s", app),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(this)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
