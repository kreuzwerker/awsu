package config

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type Profile struct {
	AccessKeyID     string `ini:"aws_access_key_id"`
	SecretAccessKey string `ini:"aws_secret_access_key"`
	ExternalID      string `ini:"external_id"`
	MFASerial       string `ini:"mfa_serial"`
	Name            string
	RoleARN         string `ini:"role_arn"`
	SourceProfile   string `ini:"source_profile"`
}

func (p *Profile) IsLongTerm() bool {
	return p.RoleARN == ""
}

func (p *Profile) Value() credentials.Value {

	return credentials.Value{
		AccessKeyID:     p.AccessKeyID,
		SecretAccessKey: p.SecretAccessKey,
	}

}
