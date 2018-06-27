package yubikey

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
)

// ARN extends the AWS ARN to functions that map to elements of a Yubikey OATH identifier
type ARN arn.ARN

// Name is the "name" argument of an OATH key registration
func (a ARN) Name() string {
	return strings.TrimPrefix(a.Resource, "mfa/")
}

// Issuer is the "--issuer" parameter of an OATH key registration
func (a ARN) Issuer() string {
	return fmt.Sprintf("aws/iam/%s", a.AccountID)
}

// Query is the argument that the ykmango package expects for code generation
func (a ARN) Query() string {
	return strings.Join([]string{
		a.Issuer(),
		a.Name(),
	}, ":")
}
