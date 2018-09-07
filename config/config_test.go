package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitLoad(t *testing.T) {

	var (
		assert  = assert.New(t)
		require = require.New(t)
	)

	split := &Config{
		ConfigFile:            "testdata/config",
		Duration:              5 * time.Minute,
		Grace:                 1 * time.Minute,
		SharedCredentialsFile: "testdata/credentials",
	}

	merged := &Config{
		ConfigFile:            "testdata/config-merged",
		Duration:              5 * time.Minute,
		Grace:                 1 * time.Minute,
		SharedCredentialsFile: "testdata/merged",
	}

	for _, config := range []*Config{split, merged} {

		err := config.Init()
		require.NoError(err)

		source, ok := config.Profiles["default"]
		require.True(ok)

		assert.Equal("default", source.Name)

		assert.Equal("AKIAIOSFODNN7EXAMPLE", source.AccessKeyID)
		assert.Equal("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", source.SecretAccessKey)
		assert.Equal(source.AccessKeyID, source.Value().AccessKeyID)
		assert.Equal(source.SecretAccessKey, source.Value().SecretAccessKey)

		profile, ok := config.Profiles["foo"]
		require.True(ok)

		assert.Equal("foo", profile.Name)

		assert.Equal("default", profile.SourceProfile)
		assert.Equal("123456", profile.ExternalID)
		assert.Equal("arn:aws:iam::123456789012:mfa/jonsmith", profile.MFASerial)
		assert.Equal("arn:aws:iam::123456789012:role/marketingadmin", profile.RoleARN)
		assert.Equal(profile.AccessKeyID, profile.Value().AccessKeyID)
		assert.Equal(profile.SecretAccessKey, profile.Value().SecretAccessKey)

	}

}

func TestInitValidate(t *testing.T) {

	var assert = assert.New(t)

	valid := &Config{
		Duration: 5 * time.Minute,
		Grace:    1 * time.Minute,
	}

	assert.NoError(valid.Init())

	invalid := &Config{
		Duration: 1 * time.Minute,
		Grace:    5 * time.Minute,
	}

	assert.Errorf(invalid.Init(), "invalid grace (f) for duration (d)")

}
