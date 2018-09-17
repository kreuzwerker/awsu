package strategy

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/kreuzwerker/awsu/config"
	"github.com/kreuzwerker/awsu/log"
	"github.com/kreuzwerker/awsu/strategy/credentials"
	"github.com/pkg/errors"
)

const (
	errAssumeRoleFailed = "failed to assume role %q"
	logAssumingRole     = "assuming role %q using profile %s (sid %s)"
)

// AssumeRole is a strategy that assumes IAM roles
type AssumeRole struct {
	Duration time.Duration
	Grace    time.Duration
	Profiles []*config.Profile
}

// Credentials aquires actual credentials
func (a *AssumeRole) Credentials(sess *session.Session) (*credentials.Credentials, error) {

	var (
		client  = sts.New(sess)
		profile = a.Profile()
		sid     = a.sessionName()
	)

	log.Debug(logAssumingRole, profile.RoleARN, profile.Name, sid)

	req := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(int64(a.Duration.Seconds())),
		RoleArn:         &profile.RoleARN,
		RoleSessionName: &sid,
	}

	if profile.ExternalID != "" {
		req.ExternalId = &profile.ExternalID
	}

	res, err := client.AssumeRole(req)

	if err != nil {
		return nil, errors.Wrapf(err, errAssumeRoleFailed, profile.RoleARN)
	}

	creds := credentials.NewShortTerm(
		profile.Name,
		*res.Credentials.AccessKeyId,
		*res.Credentials.SecretAccessKey,
		*res.Credentials.SessionToken,
		time.Now().Add(a.Duration).Add(a.Grace*-1),
	)

	return creds, nil

}

// IsCacheable indicates the output of this strategy can be cached (always true)
func (a *AssumeRole) IsCacheable() bool {
	return true
}

// Name returns the name of this strategy
func (a *AssumeRole) Name() string {
	return "assume_role"
}

// Profile returns the name of the profile used (if applicable, otherwise nil)
func (a *AssumeRole) Profile() *config.Profile {

	for _, profile := range a.Profiles {

		if profile != nil && !profile.IsLongTerm() {
			return profile
		}

	}

	return nil

}

// sessionName returns the name give to the assumed role sessions
func (a *AssumeRole) sessionName() string {

	var rid [16]byte

	io.ReadFull(rand.Reader, rid[:])

	// TODO: maybe add escaped profile name
	return fmt.Sprintf("awsu-%x", rid[:])

}
