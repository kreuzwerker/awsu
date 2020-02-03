package provided

import (
	"os"
	"time"
)

// Provided is a generator that is based on input from the environment
type Provided struct {
}

// New initializes a new provided generator
func New() *Provided {

	return &Provided{}
}

// Generate generates a new OTP by from the environment
func (m *Provided) Generate(clock time.Time, name string) (string, error) {

	return os.Getenv("AWSU_PROVIDED_OTP"), nil
}

// Name returns the name of this generator
func (m *Provided) Name() string {
	return "provided"
}
