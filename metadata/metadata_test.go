package metadata

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stretchr/testify/assert"
)

type assumeRoleProvider struct {
	ExternalID   *string
	SerialNumber *string
}

func (a *assumeRoleProvider) IsExpired() bool {
	return false
}

func (a *assumeRoleProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{}, nil
}

func TestNew(t *testing.T) {

	assert := assert.New(t)

	m := New(credentials.NewEnvCredentials())

	assert.Empty(m.ExternalID)
	assert.Empty(m.SerialNumber)

	m = New(credentials.NewCredentials(&assumeRoleProvider{
		ExternalID:   aws.String("external-id"),
		SerialNumber: aws.String("serial-number"),
	}))

	assert.Equal("external-id", m.ExternalID)
	assert.Equal("serial-number", m.SerialNumber)

}
