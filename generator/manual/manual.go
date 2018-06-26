package manual

import (
	"bufio"
	"fmt"
	"os"
)

var reader = os.Stdin

// Generate will ask a commandline user for input via stdin
func Generate(serial string) (string, error) {

	fmt.Printf("enter TOTP token for %q: ", serial)

	scanner := bufio.NewScanner(reader)
	scanner.Scan()

	return scanner.Text(), nil

}
