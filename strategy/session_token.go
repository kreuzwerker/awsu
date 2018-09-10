package strategy

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/kreuzwerker/awsu/config"
	"github.com/kreuzwerker/awsu/log"
	"github.com/kreuzwerker/awsu/source"
	"github.com/kreuzwerker/awsu/source/manual"
	"github.com/kreuzwerker/awsu/source/yubikey"
	"github.com/kreuzwerker/awsu/strategy/credentials"
	"github.com/kreuzwerker/awsu/target/mfa"
)

const (
	GenManual  = "manual"
	GenYubikey = "yubikey"
)

const (
	errSessionTokenOnUnsuitableProfiles = "failed to get session token on unsuitable profiles: at least one long-term keypair must be configured"
	errSessionTokenWithoutMFA           = "failed to get session token on unsuitable profiles: at least one MFA must be configured"
	errUnknownGenerator                 = "unknown generator %q"
	logGettingSessionToken              = "getting session token for profile %q and serial %q"
	logSerialExplicit                   = "using explicitly supplied MFA serial %q"
	logSerialFromProfile                = "using %q profile for MFA serial"
)

type SessionToken struct {
	Duration  time.Duration
	Generator string
	Grace     time.Duration
	MFASerial string
	Profiles  []*config.Profile
	_serial   string
}

func (s *SessionToken) Credentials(sess *session.Session) (*credentials.Credentials, error) {

	var (
		client = sts.New(sess)
		lt     = s.Profile()
		serial = s.serial()
	)

	log.Debug(logGettingSessionToken, lt.Name, serial)

	token, err := s.generate(&serial)

	if err != nil {
		return nil, err
	}

	res, err := client.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(int64(s.Duration.Seconds())),
		SerialNumber:    &serial,
		TokenCode:       &token,
	})

	if err != nil {
		return nil, err
	}

	creds := credentials.NewShortTerm(
		lt.Name,
		*res.Credentials.AccessKeyId,
		*res.Credentials.SecretAccessKey,
		*res.Credentials.SessionToken,
		time.Now().Add(s.Duration).Add(s.Grace*-1),
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

		if profile != nil && profile.IsLongTerm() && s.serial() != "" {
			return profile
		}

	}

	return nil

}

func (s *SessionToken) generate(serial *string) (string, error) {

	var g source.Generator

	switch s.Generator {
	case GenYubikey:

		gen, err := yubikey.New()

		if err != nil {
			return "", err
		}

		g = gen

	case GenManual:
		g = manual.New()

	default:
		return "", fmt.Errorf(errUnknownGenerator, s.Generator)
	}

	name, err := mfa.SerialToName(serial)

	if err != nil {
		return "", err
	}

	return g.Generate(time.Now(), name)

}

func (s *SessionToken) serial() string {

	if s._serial != "" {
		return s._serial
	}

	if s._serial = s.MFASerial; s._serial != "" {
		log.Debug(logSerialExplicit, s._serial)
		return s._serial
	}

	// find the MFA
	for _, profile := range s.Profiles {

		if profile != nil && profile.MFASerial != "" {
			s._serial = profile.MFASerial
			log.Debug(logSerialFromProfile, profile.Name)
			return s._serial
		}

	}

	// TODO: try autodetection as a last resort OR just don't get a session token?

	return s._serial

}
