package strategy

import (
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/gesellix/awsu/config"
	"github.com/gesellix/awsu/strategy/credentials"
)

// LongTerm is a strategy that uses long-term credentials (IAM user keypairs)
type LongTerm struct {
	Profiles []*config.Profile
}

// Credentials aquires actual credentials
func (l *LongTerm) Credentials(sess *session.Session) (*credentials.Credentials, error) {

	p := l.Profile()

	return credentials.NewLongTerm(p.Name,
			p.AccessKeyID,
			p.SecretAccessKey),
		nil

}

// IsCacheable indicates the output of this strategy can be cached (always false)
func (l *LongTerm) IsCacheable() bool {
	return false
}

// Name returns the name of this strategy
func (l *LongTerm) Name() string {
	return "long_term"
}

// Profile returns the name of the profile used (if applicable, otherwise nil)
func (l *LongTerm) Profile() *config.Profile {

	for _, profile := range l.Profiles {

		if profile != nil && profile.IsLongTerm() {
			return profile
		}

	}

	return nil

}
