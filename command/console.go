package command

import (
	"fmt"

	"github.com/kreuzwerker/awsu/console"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var consoleCmd = &cobra.Command{

	Use:   "console",
	Short: "Generates link to or opens AWS console",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.Unmarshal(&conf.Console)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		cons, err := console.New(&conf)

		if err != nil {
			return err
		}

		link, err := cons.Link()

		if err != nil {
			return err
		}

		if conf.Console.Open {
			return open.Run(link)
		}

		fmt.Println(link)

		return nil

	},
}

func init() {

	flag(consoleCmd.Flags(),
		true,
		"open",
		"o",
		"",
		"attempts to open the generated url in a browser",
	)

	rootCmd.AddCommand(consoleCmd)

}
