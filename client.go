package awsu

import (
	"crypto/rand"
	"fmt"
	"io"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/yawn/envmap"
)

type Whoami struct {
	Account  string
	ARN      string
	Username string
}

var arn = regexp.MustCompile(`^arn:aws:iam::(?P<account>\d+):(?P<username>.+)$`)

type Client struct {
	*sts.STS
}

type Options struct {
	ARN        string
	ExternalID string
}

func NewClient() (*Client, error) {

	sess, err := session.NewSession()

	if err != nil {
		return nil, fmt.Errorf("error while opening new session: %s", err)
	}

	return &Client{sts.New(sess)}, nil

}

func (c *Client) AssumeRole(o *Options) (envmap.Envmap, error) {

	const (
		empty = ""
		max   = 3600
	)

	params := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(max),
		RoleArn:         aws.String(o.ARN),
		RoleSessionName: aws.String(c.sessionName()),
	}

	if o.ExternalID != empty {
		params.ExternalId = aws.String(o.ExternalID)
	}

	res, err := c.STS.AssumeRole(params)

	if err != nil {
		return nil, fmt.Errorf("error during assume role call: %s", err)
	}

	return map[string]string{
		AccessKeyID:     *res.Credentials.AccessKeyId,
		SecretAccessKey: *res.Credentials.SecretAccessKey,
		SessionToken:    *res.Credentials.SessionToken,
	}, nil

}

func (c *Client) CallerIdentity() (*Whoami, error) {

	res, err := c.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		return nil, fmt.Errorf("error while getting caller identity: %s", err)
	}

	match := arn.FindStringSubmatch(*res.Arn)

	return &Whoami{
		ARN:      match[0],
		Account:  match[1],
		Username: match[2],
	}, nil

}

func (c *Client) sessionName() string {

	var buf = make([]byte, 16)

	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		panic(err)
	}

	return fmt.Sprintf("awssu-%x", buf)

}
