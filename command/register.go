package command

import (
	"encoding/base32"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kreuzwerker/awsu/log"
	"github.com/kreuzwerker/awsu/source/yubikey"
	"github.com/kreuzwerker/awsu/strategy"
	"github.com/kreuzwerker/awsu/target/mfa"
	qr "github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
)

var registerFlags = struct {
	issuer string
	qr     bool
}{}

var registerCmd = &cobra.Command{

	Use:   "register :username",
	Short: "Initializes an device on AWS and Yubikey",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		username := args[0]

		creds, err := strategy.Apply(config)

		if err != nil {
			return err
		}

		source, err := yubikey.New()

		if err != nil {
			return err
		}

		sess := session.Must(session.NewSession())

		target, err := mfa.New(creds.UpdateSession(sess), source)

		if err != nil {
			return err
		}

		serial, secret, err := target.Add(username)

		if err != nil {
			log.Info("failed to register %q, attemping to delete from target", username)
			return target.Delete(username)
		}

		if registerFlags.qr {

			uri := fmt.Sprintf("otpauth://totp/%s@%s?secret=%s&issuer=%s",
				username,
				creds.Profile,
				base32.StdEncoding.EncodeToString(secret),
				registerFlags.issuer,
			)

			qr.Generate(uri, qr.L, os.Stderr)

		}

		log.Info("MFA %q serial successfully registered", *serial)

		return nil

	},
}

func init() {

	flag(registerCmd.Flags(),
		&registerFlags.issuer,
		"Amazon",
		"issuer",
		"i",
		"AWSU_QR_ISSUER",
		"issuer parameter in the QR key uri",
	)

	flag(registerCmd.Flags(),
		&registerFlags.qr,
		true,
		"qr",
		"q",
		"AWSU_QR",
		"generate a QR barcode as backup for soft tokens",
	)

	rootCmd.AddCommand(registerCmd)

}
