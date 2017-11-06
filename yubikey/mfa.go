package yubikey

import (
	"crypto/rand"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	otp "github.com/hgfischer/go-otp"
	"github.com/kreuzwerker/awsu/log"
	"github.com/pkg/errors"
)

// NewMFA generates a new virtual MFA device and associates it with a Yubikey
func NewMFA(sess *session.Session, qr func(string), username string) (*arn.ARN, error) {

	client := iam.New(sess)

	device, err := createDevice(client, qr)

	if err != nil {
		return nil, err
	}

	a, err := arn.Parse(*device.SerialNumber)

	if err != nil {
		return nil, err
	}

	if err = associateDevice(ARN(a), device.Base32StringSeed); err != nil {
		return nil, errors.Wrapf(err, "failed to associate device with Yubikey, consider deleting obsolete MFA device %q", *device.SerialNumber)
	}

	if err = enableDevice(client, device, username); err != nil {
		return nil, errors.Wrapf(err, "failed to enable device in IAM, consider deleting obsolete MFA device %q", *device.SerialNumber)
	}

	return &a, nil

}

// associateDevice will provision the MFA secret into a present Yubikey
func associateDevice(a ARN, secret []byte) error {

	// TODO: maybe port this to ykmango

	out, err := exec.Command("ykman",
		"oath",
		"add",
		"--oath-type", "TOTP",
		"--digits", "6",
		"--algorithm", "SHA1",
		"--period", "30",
		"--issuer", a.Issuer(),
		a.Name(),
		string(secret),
	).CombinedOutput()

	if err != nil {
		return errors.Wrapf(err, "failed to associated device with yubikey: %s", string(out))
	}

	return nil

}

// createDevice creates a new virtual MFA device and calls the qr function with it's secret
func createDevice(client *iam.IAM, qr func(string)) (*iam.VirtualMFADevice, error) {

	id, err := newDeviceID()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate new device ID")
	}

	req := &iam.CreateVirtualMFADeviceInput{
		VirtualMFADeviceName: id,
	}

	res, err := client.CreateVirtualMFADevice(req)

	if err != nil {
		return nil, err
	}

	qr(string(res.VirtualMFADevice.Base32StringSeed))

	return res.VirtualMFADevice, nil

}

// enabledDevice associates the device with a user by completing two consecutive TOTP challenges
func enableDevice(client *iam.IAM, device *iam.VirtualMFADevice, username string) error {

	oath := &otp.TOTP{
		IsBase32Secret: true,
		Length:         6,
		Period:         30,
		Secret:         string(device.Base32StringSeed),
		Time:           time.Now(),
	}

	var (
		code1 = oath.Get()
		code2 string
	)

	// advance to the next period
	oath.Time = oath.Time.Add(30 * time.Second)

	code2 = oath.Get()

	log.Log("trying to enable MFA devices with codes %q and %q (secret %s)", code1, code2, oath.Secret)

	req := &iam.EnableMFADeviceInput{
		AuthenticationCode1: &code1,
		AuthenticationCode2: &code2,
		SerialNumber:        device.SerialNumber,
		UserName:            &username,
	}

	_, err := client.EnableMFADevice(req)

	return err

}

// newDeviceID generates a new random ID for the virtual MFA device
func newDeviceID() (*string, error) {

	buf := make([]byte, 16)

	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return nil, err
	}

	id := fmt.Sprintf("awsu-%x", buf)

	return &id, nil

}
