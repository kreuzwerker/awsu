package command

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/gesellix/awsu/source/yubikey"
	"github.com/gesellix/awsu/target/mfa"
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

		otp, err := token.Generate(time.Now(), name)

		if err != nil {
			return err
		}

		fmt.Printf("%s", otp)

		return nil

	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
