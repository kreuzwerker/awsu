package mfa

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallerIdentityToSerial(t *testing.T) {

	var (
		assert  = assert.New(t)
		require = require.New(t)
	)

	name, err := CallerIdentityToSerial(aws.String("arn:aws:iam::1234567890:user/foo"))

	require.NoError(err)
	assert.Equal("arn:aws:iam::1234567890:mfa/foo", name)

	name, err = CallerIdentityToSerial(aws.String("foo"))

	assert.EqualError(err, `failed to parse "foo" as ARN: arn: invalid prefix`)
	assert.Empty(name)

}

func TestSerialToName(t *testing.T) {

	var (
		assert  = assert.New(t)
		require = require.New(t)
	)

	name, err := SerialToName(aws.String("arn:aws:iam::1234567890:mfa/foo"))

	require.NoError(err)
	assert.Equal("aws/iam/1234567890:foo", name)

	name, err = SerialToName(aws.String("foo"))

	assert.EqualError(err, `failed to parse "foo" as ARN: arn: invalid prefix`)
	assert.Empty(name)

}
