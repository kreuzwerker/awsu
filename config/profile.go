package config

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// Profile is a long- or short-time credential profile managed in a shared config
type Profile struct {
	AccessKeyID     string `ini:"aws_access_key_id"`
	ExternalID      string `ini:"external_id"`
	MFASerial       string `ini:"mfa_serial"`
	Name            string
	Region          string `ini:"region"`
	RoleARN         string `ini:"role_arn"`
	SecretAccessKey string `ini:"aws_secret_access_key"`
	SourceProfile   string `ini:"source_profile"`
}

// IsLongTerm identifies a profile that does not assume a role using a source profile
func (p *Profile) IsLongTerm() bool {
	return p.RoleARN == ""
}

// Value returns the credentials associated with the profile (if any) - only long
// term profiles have credentials
func (p *Profile) Value() credentials.Value {

	return credentials.Value{
		AccessKeyID:     p.AccessKeyID,
		SecretAccessKey: p.SecretAccessKey,
	}

}
