package source

import "time"

// Generator is the interface defined by read-only sources of OTPs
type Generator interface {

	// Generate generates a named one-time password using the given reference time
	Generate(clock time.Time, name string, requireTouch bool) (string, error)

	// Name returns the name of this generator
	Name() string
}

// Source is the interface defined by r/w sources of OTPS
type Source interface {
	Generator

	// Add adds a new named TOTP secret
	Add(name string, secret []byte, requireTouch bool) error

	// Delete removes a named TOTP secret
	Delete(name string) error
}
