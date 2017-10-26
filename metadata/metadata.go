package metadata

import (
	"reflect"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

type Metadata struct {
	ExternalID   string `json:",omitempty"`
	SerialNumber string `json:",omitempty"`
}

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
