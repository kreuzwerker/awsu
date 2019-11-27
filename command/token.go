package command

import (
	"fmt"
	"time"

	"github.com/kreuzwerker/awsu/source/yubikey"
	"github.com/kreuzwerker/awsu/target/mfa"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const errTokenProfileNotFound = "no such profile or no direct MFA configured for profile %q"

var tokenCmd = &cobra.Command{

	Use:   "token",
	Short: "Generates one-time password from Yubikey",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.Unmarshal(&conf.Register)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		var serial = conf.MFASerial

		if serial == "" {

			profile, ok := conf.Profiles[conf.Profile]

			if !ok || profile.MFASerial == "" {
				return fmt.Errorf(errTokenProfileNotFound, conf.Profile)
			}

			serial = profile.MFASerial

		}

		name, err := mfa.SerialToName(&serial)

		if err != nil {
			return err
		}

		token, err := yubikey.New()

		if err != nil {
			return err
		}

		requireTouch, err := cmd.Flags().GetBool("require-touch")
		if err != nil {
			return err
		}

		otp, err := token.Generate(time.Now(), name, requireTouch)

		if err != nil {
			return err
		}

		fmt.Printf("%s", otp)

		return nil

	},
}

func init() {
	flag(tokenCmd.Flags(),
		false,
		"require-touch",
		"t",
		"AWSU_REQUIRE_TOUCH",
		"require touch to generate Yubikey token",
	)

	rootCmd.AddCommand(tokenCmd)
}
