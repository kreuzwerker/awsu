package log

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {

	assert := assert.New(t)

	buf := bytes.NewBuffer(nil)
	logger.SetOutput(buf)

	defer func() {
		logger.SetOutput(os.Stderr)
		Debug = false
	}()

	Log("hello")

	assert.Empty(buf.String())

	Debug = true

	Log("goodbye")

	assert.Regexp(`\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2} \[DEBUG\] goodbye\n`, buf.String())

}
