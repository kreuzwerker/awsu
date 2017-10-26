package metadata

import (
	"reflect"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

// Metadata is additional provider metadata like the serial number of the MFA
// and the external ID of the non-cross-account role
type Metadata struct {
	ExternalID   string `json:",omitempty"`
	SerialNumber string `json:",omitempty"`
}

// New extracts additional provider metadata (if applicable) from the given
// credentials
func New(c *credentials.Credentials) Metadata {

	v1 := reflect.Indirect(reflect.ValueOf(c))
	v2 := reflect.Indirect(v1.FieldByName("provider").Elem())
	v3 := reflect.Indirect(v2.FieldByName("ExternalID"))
	v4 := reflect.Indirect(v2.FieldByName("SerialNumber"))

	metadata := Metadata{}

	if v3.IsValid() {
		metadata.ExternalID = v3.String()
	}

	if v4.IsValid() {
		metadata.SerialNumber = v4.String()
	}

	return metadata

}
