package aquirer

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kreuzwerker/awsu/config"
)

type LongTerm struct {
	Profiles []*config.Profile
}

func (l *LongTerm) Credentials(sess *session.Session) (*Credentials, error) {

	p := l.Profile()

	return newLongTermCredentials(p.Name,
			p.AccessKeyID,
			p.SecretAccessKey),
		nil

}

func (l *LongTerm) IsCacheable() bool {
	return false
}

func (l *LongTerm) Name() string {
	return "long_term"
}

func (l *LongTerm) Profile() *config.Profile {

	for _, profile := range l.Profiles {

		if profile != nil && profile.IsLongTerm() {
			return profile
		}

	}

	return nil

}
