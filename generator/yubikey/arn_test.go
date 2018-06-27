package yubikey

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/stretchr/testify/assert"
)

func TestARN(t *testing.T) {

	assert := assert.New(t)

	a, err := arn.Parse("arn:aws:iam::1234567890:mfa/foo")

	assert.NoError(err)

	assert.Equal("aws/iam/1234567890", ARN(a).Issuer())
	assert.Equal("foo", ARN(a).Name())
	assert.Equal("aws/iam/1234567890:foo", ARN(a).Query())

}
