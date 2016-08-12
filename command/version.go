package command

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type Version struct {
	Build, Version string
}

func (v Version) Session() *string {

	var buf = make([]byte, 16)

	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		panic(err)
	}

	session := fmt.Sprintf("%s-%s-%s-%x",
		app,
		v.Version,
		v.Build,
		buf)

	return &session

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
