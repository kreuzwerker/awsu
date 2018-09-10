package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/kreuzwerker/awsu/strategy"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var consoleFlags = struct {
	open bool
}{}

var consoleCmd = &cobra.Command{

	Use:   "console",
	Short: "Generates link to or opens AWS console",
	RunE: func(cmd *cobra.Command, args []string) error {

		creds, err := strategy.Apply(rootFlags)

		if err != nil {
			return err
		}

		profile := rootFlags.Profiles[rootFlags.Profile]

		// TODO: move the whole logic to a "console" package
		if profile.ExternalID == "" {

			a, err := arn.Parse(profile.RoleARN)

			if err != nil {
				return err
			}

			url := fmt.Sprintf("https://signin.aws.amazon.com/switchrole?account=%s&roleName=%s&displayName=%s",
				a.AccountID,
				strings.TrimPrefix(a.Resource, "role/"),
				profile.Name)

			if consoleFlags.open {
				return open.Run(url)
			}

			fmt.Println(url)

		} else {

			fep := map[string]string{
				"sessionId":    creds.AccessKeyID,
				"sessionKey":   creds.SessionToken,
				"sessionToken": creds.SessionToken,
			}

			enc, err := json.Marshal(fep)

			if err != nil {
				return fmt.Errorf("error while marshaling federation session: %s", err)
			}

			url := fmt.Sprintf("https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=%s", string(url.QueryEscape(string(enc))))

			var buf = bytes.NewBuffer(nil)

			res, err := http.Get(url)

			if err != nil {
				return fmt.Errorf("error while requesting federation: %s", err)
			}

			defer res.Body.Close()

			if _, err := io.Copy(buf, res.Body); err != nil {
				return fmt.Errorf("error while receiving federation response body: %s", err)
			}

			var body map[string]string

			if err := json.Unmarshal(buf.Bytes(), &body); err != nil {
				return fmt.Errorf("error while unmarshaling sign-in token: %s", err)
			}

			var (
				destination = "https://console.aws.amazon.com/"
				issuer      = ""
				token       = body["SigninToken"]
			)

			url = fmt.Sprintf("https://signin.aws.amazon.com/federation?Action=login&Issuer=%s&Destination=%s&SigninToken=%s\n",
				issuer,
				destination,
				token)

			if consoleFlags.open {
				return open.Run(url)
			}

			fmt.Println(url)

		}

		return nil

	},
}

func init() {

	flag(consoleCmd.Flags(),
		&consoleFlags.open,
		true,
		"open",
		"o",
		"",
		"attempts to open the generated url in a browser",
	)

	rootCmd.AddCommand(consoleCmd)

}
