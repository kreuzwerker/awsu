package mfa

import (
	"encoding/base32"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/kreuzwerker/awsu/log"
	"github.com/kreuzwerker/awsu/source"
	"github.com/pkg/errors"
)

const (
	errAddToSource               = "failed to add %q to source %q"
	errCalculateFirstOTP         = "failed to calculate first one-time password"
	errCalculateSecondOTP        = "failed to calculate second one-time password"
	errCreateVirtualDevice       = "failed to create virtual AWS MFA device"
	errDeactivateVirtualDevice   = "failed to deactivate virtual AWS MFA device with serial %q"
	errDecodeSecret              = "failed to decode virtual AWS MFA device secret"
	errDeleteSerialDetermination = "failed to determine serial number for device deletion"
	errDeleteVirtualDevice       = "failed to delete virtual AWS MFA device with serial %q"
	errEnableVirtualDevice       = "failed to enable virtual AWS MFA device"
	errRemoveFromSource          = "failed to remove %q from source %q"
	logAddSecretToSource         = "adding secret to source %q"
	logCreateVirtualDevice       = "create virtual AWS MFA device"
	logDeactivateVirtualDevice   = "disable virtual AWS MFA device with serial %q"
	logDeleteVirtualDevice       = "delete virtual AWS MFA device with serial %q"
	logEnableVirtualDevice       = "enable virtual AWS MFA devices with codes %q and %q"
)

// MFA implements AWS virtual MFA devices as target
type MFA struct {
	iam    *iam.IAM
	source source.Source
	sts    *sts.STS
}

// New initializes a AWS virtual MFA device as target
func New(sess *session.Session, s source.Source) (*MFA, error) {

	return &MFA{
		iam:    iam.New(sess),
		sts:    sts.New(sess),
		source: s,
	}, nil

}

// Add adds virtual MFA to the source and associates it with the given IAM
// username and returns MFA serial and TOTP secret
func (m *MFA) Add(username string) (*string, []byte, error) {

	serial, secret, err := m.create(username)

	if err != nil {
		return nil, nil, err
	}

	if err := m.enable(username, serial, secret); err != nil {
		return nil, nil, err
	}

	return serial, secret, nil

}

// Delete removes a virtual MFA from the source including it's association with
// the given IAM username
func (m *MFA) Delete(username string) error {

	res, err := m.sts.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		return errors.Wrapf(err, errDeleteSerialDetermination)
	}

	serial, err := CallerIdentityToSerial(res.Arn)

	if err != nil {
		return err
	}

	// ignore errors here
	m.deactivate(username, &serial)

	if err := m.delete(&serial); err != nil {
		return err
	}

	return nil

}

// create creates the virtual MFA device
func (m *MFA) create(username string) (*string, []byte, error) {

	log.Debug(logCreateVirtualDevice)

	res, err := m.iam.CreateVirtualMFADevice(&iam.CreateVirtualMFADeviceInput{
		VirtualMFADeviceName: &username,
	})

	if err != nil {
		return nil, nil, errors.Wrapf(err, errCreateVirtualDevice)
	}

	secret, err := base32.StdEncoding.DecodeString(string(res.VirtualMFADevice.Base32StringSeed))

	if err != nil {
		return nil, nil, errors.Wrapf(err, errDecodeSecret)
	}

	return res.VirtualMFADevice.SerialNumber, secret, nil

}

// deactivate deactivates the virtual MFA device and removes it from the source
func (m *MFA) deactivate(username string, serial *string) error {

	log.Debug(logDeactivateVirtualDevice, *serial)

	_, err := m.iam.DeactivateMFADevice(&iam.DeactivateMFADeviceInput{
		SerialNumber: serial,
		UserName:     &username,
	})

	if err != nil {
		return errors.Wrapf(err, errDeactivateVirtualDevice, *serial)
	}

	name, err := SerialToName(serial)

	if err != nil {
		return err
	}

	if err := m.source.Delete(name); err != nil {
		return errors.Wrapf(err, errRemoveFromSource, name, m.source.Name())
	}

	return nil

}

// delete deletes the virtual MFA device
func (m *MFA) delete(serial *string) error {

	log.Debug(logDeleteVirtualDevice, *serial)

	_, err := m.iam.DeleteVirtualMFADevice(&iam.DeleteVirtualMFADeviceInput{
		SerialNumber: serial,
	})

	if err != nil {
		return errors.Wrapf(err, errDeleteVirtualDevice, *serial)
	}

	return nil

}

// enable enables the virtual MFA device and adds it from the source
func (m *MFA) enable(username string, serial *string, secret []byte) error {

	log.Debug(logAddSecretToSource, m.source.Name())

	name, err := SerialToName(serial)

	if err != nil {
		return err
	}

	if err = m.source.Add(name, secret); err != nil {
		return errors.Wrapf(err, errAddToSource, name, m.source.Name())
	}

	otp1, err := m.source.Generate(time.Now(), name)

	if err != nil {
		return errors.Wrapf(err, errCalculateFirstOTP)
	}

	otp2, err := m.source.Generate(time.Now().Add(30*time.Second), name)

	if err != nil {
		return errors.Wrapf(err, errCalculateSecondOTP)
	}

	log.Debug(logEnableVirtualDevice, otp1, otp2)

	if _, err := m.iam.EnableMFADevice(&iam.EnableMFADeviceInput{
		AuthenticationCode1: &otp1,
		AuthenticationCode2: &otp2,
		SerialNumber:        serial,
		UserName:            &username,
	}); err != nil {
		return errors.Wrapf(err, errEnableVirtualDevice)
	}

	return nil

}
