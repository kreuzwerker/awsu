package yubikey

import (
	"github.com/aws/aws-sdk-go/aws/arn"
	ykman "github.com/joshdk/ykmango"
	"github.com/kreuzwerker/awsu/log"
	"github.com/pkg/errors"
)

func Generate(mfa string) (string, error) {

	arn, err := arn.Parse(mfa)

	if err != nil {
		return "", errors.Wrapf(err, "cannot parse ARN from %q", mfa)
	}

	a := ARN(arn)

	log.Log("asking for yubikey OATH slot with issuer %q and name %q", a.Issuer(), a.Name())

	// TODO: check for touch-bit in metadata

	code, err := ykman.Generate(a.Query())

	if err != nil {
		switch err {
		case ykman.ErrorSlotNameUnknown:
			return "", errors.Wrapf(err, "cannot find registered MFA for issuer %q and name %q", a.Issuer(), a.Name())
		case ykman.ErrorYkmanNotFound, ykman.ErrorYubikeyNotDetected, ykman.ErrorYubikeyTimeout:
			return "", errors.Wrapf(err, "cannot connect to yubikey")
		default:
			return "", errors.Wrapf(err, "cannot communicate with yubikey")
		}
	}

	log.Log("received %q as code", code)

	return code, nil

}
