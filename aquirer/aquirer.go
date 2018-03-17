package aquirer

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kreuzwerker/awsu/config"
)

type Aquirer interface {
	// Credentials aquires actual credentials
	Credentials(*session.Session) (*Credentials, error)

	// Name returns the name of this aquirer
	Name() string

	// IsCacheable indicates the output of this aquirer can be cached
	IsCacheable() bool

	// Profiles returns the name of the profile used (if applicable)
	Profile() *config.Profile
}
