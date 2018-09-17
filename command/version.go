package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version describe the version of awsu
type Version struct {
	Build, Time, Version string
}

// String returns a string representation of the version
func (v Version) String() string {
	return fmt.Sprintf("%s version %s (%s, %s)", app, v.Version, v.Build, v.Time)
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
