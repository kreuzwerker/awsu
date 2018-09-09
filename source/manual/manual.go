package manual

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

// Manual is a generator that is based on manual input
type Manual struct {
	reader io.Reader
}

// New initializes a new manual generator
func New() *Manual {

	return &Manual{
		reader: os.Stdin,
	}

}

// Generate generates a new OTP by asking for it on the commandline
func (m *Manual) Generate(clock time.Time, name string) (string, error) {

	fmt.Printf("enter TOTP token for %q: ", name)

	scanner := bufio.NewScanner(m.reader)
	scanner.Scan()

	return scanner.Text(), nil

}
