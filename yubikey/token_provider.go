package yubikey

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kreuzwerker/awsu/metadata"
)

type TokenProvider struct {
	Session *session.Session
}

func (t *TokenProvider) Provide() (string, error) {
	metadata := metadata.New(t.Session.Config.Credentials)
	return Generate(metadata.SerialNumber)
}
