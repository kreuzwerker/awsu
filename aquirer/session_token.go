package aquirer

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/kreuzwerker/awsu/config"
	"github.com/kreuzwerker/awsu/log"
	"github.com/kreuzwerker/awsu/yubikey"
)

const (
	errSessionTokenWithoutMFA           = "failed to get session token on unsuitable profiles: at least one MFA must be configured"
	errSessionTokenOnUnsuitableProfiles = "failed to get session token on unsuitable profiles: at least one long-term keypair must be configured"
)

var tokenSource = yubikey.Generate

type SessionToken struct {
	Duration time.Duration
	Grace    time.Duration
	Profiles []*config.Profile
}

func (s *SessionToken) Credentials(sess *session.Session) (*Credentials, error) {

	var (
		client       = sts.New(sess)
		serialNumber string
		lt           = s.Profile()
	)

	// find the MFA
	for _, profile := range s.Profiles {

		if profile != nil && profile.MFASerial != "" {
			serialNumber = profile.MFASerial
			log.Log("using %q for MFA serial", profile.Name)
			break
		}

	}

	// TODO: try autodetection as a last resort
	if serialNumber == "" {
		return nil, errors.New(errSessionTokenWithoutMFA)
	}

	log.Log("getting session token for profile %q and serial %q", lt.Name, serialNumber)

	token, err := tokenSource(serialNumber)

	if err != nil {
		return nil, err
	}

	res, err := client.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(int64(s.Duration.Seconds())),
		SerialNumber:    &serialNumber,
		TokenCode:       &token,
	})

	if err != nil {
		return nil, err
	}

	creds := newShortTermCredentials(
		lt.Name,
		*res.Credentials.AccessKeyId,
		*res.Credentials.SecretAccessKey,
		*res.Credentials.SessionToken,
		time.Now().Add(s.Duration),
	)

	return creds, nil

}

func (s *SessionToken) IsCacheable() bool {
	return true
}

func (s *SessionToken) Name() string {
	return "session_token"
}

func (s *SessionToken) Profile() *config.Profile {

	for _, profile := range s.Profiles {

		if profile != nil && profile.IsLongTerm() {
			return profile
		}

	}

	return nil

}
