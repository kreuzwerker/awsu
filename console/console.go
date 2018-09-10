package console

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/kreuzwerker/awsu/config"
	"github.com/kreuzwerker/awsu/strategy"
)

// Console is a helper for opening links to the AWS console
type Console struct {
	conf    *config.Config
	profile *config.Profile
}

// New instantiates a new console helper
func New(conf *config.Config) (*Console, error) {

	profile, ok := conf.Profiles[conf.Profile]

	if !ok {
		return nil, fmt.Errorf("no such profile %q configured", conf.Profile)
	}

	return &Console{
		conf:    conf,
		profile: profile,
	}, nil

}

// Link returns a link to the AWS console
func (c *Console) Link() (string, error) {

	var f = c.linkInternal

	if c.profile.ExternalID != "" {
		f = c.linkExternal
	}

	return f()

}

func (c *Console) linkInternal() (string, error) {

	a, err := arn.Parse(c.profile.RoleARN)

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://signin.aws.amazon.com/switchrole?account=%s&roleName=%s&displayName=%s",
		a.AccountID,
		strings.TrimPrefix(a.Resource, "role/"),
		c.profile.Name)

	return url, nil

}

func (c *Console) linkExternal() (string, error) {

	creds, err := strategy.Apply(c.conf)

	if err != nil {
		return "", err
	}

	fep := map[string]string{
		"sessionId":    creds.AccessKeyID,
		"sessionKey":   creds.SessionToken,
		"sessionToken": creds.SessionToken,
	}

	enc, err := json.Marshal(fep)

	if err != nil {
		return "", fmt.Errorf("error while marshaling federation session: %s", err)
	}

	url := fmt.Sprintf("https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=%s", string(url.QueryEscape(string(enc))))

	var buf = bytes.NewBuffer(nil)

	res, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("error while requesting federation: %s", err)
	}

	defer res.Body.Close()

	if _, err := io.Copy(buf, res.Body); err != nil {
		return "", fmt.Errorf("error while receiving federation response body: %s", err)
	}

	var body map[string]string

	if err := json.Unmarshal(buf.Bytes(), &body); err != nil {
		return "", fmt.Errorf("error while unmarshaling sign-in token: %s", err)
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

	return url, nil

}
