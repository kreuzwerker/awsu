package command

import (
	"encoding/base32"
	"fmt"
	"os"

	"github.com/kreuzwerker/awsu/log"
	"github.com/kreuzwerker/awsu/source/yubikey"
	"github.com/kreuzwerker/awsu/strategy"
	"github.com/kreuzwerker/awsu/target/mfa"
	qr "github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var registerCmd = &cobra.Command{

	Use:   "register :username",
	Short: "Initializes an device on AWS and Yubikey",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.Unmarshal(&conf.Register)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		username := args[0]

		creds, err := strategy.Apply(&conf)

		if err != nil {
			return err
		}

		source, err := yubikey.New()

		if err != nil {
			return err
		}

		target, err := mfa.New(creds.NewSession(), source)

		if err != nil {
			return err
		}

		serial, secret, err := target.Add(username)

		if err != nil {
			return err
		}

		if conf.Register.QR {

			uri := fmt.Sprintf("otpauth://totp/%s@%s?secret=%s&issuer=%s",
				username,
				creds.Profile,
				base32.StdEncoding.EncodeToString(secret),
				conf.Register.Issuer,
			)

			qr.Generate(uri, qr.L, os.Stderr)

		}

		log.Info("MFA %q serial successfully registered", *serial)

		return nil

	},
}

func init() {

	flag(registerCmd.Flags(),
		"Amazon",
		"issuer",
		"i",
		"AWSU_QR_ISSUER",
		"issuer parameter in the QR key uri",
	)

	flag(registerCmd.Flags(),
		true,
		"qr",
		"q",
		"AWSU_QR",
		"generate a QR barcode as backup for soft tokens",
	)

	rootCmd.AddCommand(registerCmd)

}
