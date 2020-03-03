package config

import (
	"testing"
	"time"

	assertmod "github.com/stretchr/testify/assert"
	requiremod "github.com/stretchr/testify/require"
)

func TestConfigInitLoad(t *testing.T) {

	var (
		assert  = assertmod.New(t)
		require = requiremod.New(t)
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

func TestConfigInitValidate(t *testing.T) {

	var assert = assertmod.New(t)

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

func TestMergeConfigAndCredentials(t *testing.T) {

	var (
		assert  = assertmod.New(t)
		require = requiremod.New(t)
	)

	config := &Config{
		ConfigFile:            "testdata/config",
		Duration:              5 * time.Minute,
		Grace:                 1 * time.Minute,
		SharedCredentialsFile: "testdata/credentials",
	}

	err := config.Init()
	require.NoError(err)

	mfaProfile, ok := config.Profiles["test-mfa"]
	require.True(ok)

	assert.Equal("test-mfa", mfaProfile.Name)

	assert.Equal("arn:aws:iam::123456789012:mfa/jondoe", mfaProfile.MFASerial)
	assert.Equal("AKIAIOSFODNN7EXAMPLE", mfaProfile.Value().AccessKeyID)
	assert.Equal("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", mfaProfile.Value().SecretAccessKey)
}
