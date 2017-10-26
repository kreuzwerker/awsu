package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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

		sess, err := newSession(rootFlags.workspace)

		if err != nil {
			return err
		}

		// TODO: make this work for CA roles
		if sess.SessionToken == "" || sess.ExternalID == "" {
			return fmt.Errorf("cannot marshal federation token without session token or external ID")
		}

		fep := map[string]string{
			"sessionId":    sess.AccessKeyID,
			"sessionKey":   sess.SessionToken,
			"sessionToken": sess.SessionToken,
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

		return nil

	},
}

func init() {

	consoleCmd.Flags().BoolVarP(&consoleFlags.open, "open", "o", true, "attempts to open the generated url")

	rootCmd.AddCommand(consoleCmd)

}
