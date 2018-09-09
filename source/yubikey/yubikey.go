package yubikey

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/yawn/ykoath"
)

// Yubikey is a source based on a Yubikey
type Yubikey struct {
	client *ykoath.OATH
	sync.Mutex
}

const (
	errInitializeSmartcard   = "failed to initialize Yubikey"
	errSelectOathApplication = "failed to select OATH application in Yubikey"
)

// New initializes a new Yubikey source
func New() (*Yubikey, error) {

	oath, err := ykoath.New()

	if err != nil {
		return nil, errors.Wrapf(err, errInitializeSmartcard)
	}

	_, err = oath.Select()

	if err != nil {
		return nil, errors.Wrapf(err, errSelectOathApplication)
	}

	return &Yubikey{
		client: oath,
	}, nil

}

// Add adds / overwrites a credential to a Yubikey
func (y *Yubikey) Add(name string, secret []byte) error {
	return y.client.Put(name, ykoath.HmacSha1, ykoath.Totp, 6, secret, false)
}

// Delete deletes a credential from a Yubikey
func (y *Yubikey) Delete(name string) error {
	return y.client.Delete(name)
}

// Generate generates a new OTP with a Yubikey
func (y *Yubikey) Generate(clock time.Time, name string) (string, error) {

	y.Lock()
	defer y.Unlock()

	defer func(prev func() time.Time) {
		y.client.Clock = prev
	}(y.client.Clock)

	y.client.Clock = func() time.Time {
		return clock
	}

	return y.client.Calculate(name, nil)

}

// Name returns the name of this source
func (y *Yubikey) Name() string {
	return "yubikey (ykoath)"
}
