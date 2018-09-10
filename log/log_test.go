package log

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {

	assert := assert.New(t)

	buf := bytes.NewBuffer(nil)
	logger.SetOutput(buf)

	defer func() {
		logger.SetOutput(os.Stderr)
		Verbose = false
	}()

	Debug("hello")

	assert.Empty(buf.String())

	Verbose = true

	Debug("goodbye")

	assert.Regexp(`\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2} \[DEBUG\] goodbye\n`, buf.String())

}

func TestInfo(t *testing.T) {

	assert := assert.New(t)

	buf := bytes.NewBuffer(nil)
	logger.SetOutput(buf)

	defer func() {
		logger.SetOutput(os.Stderr)
		Verbose = false
	}()

	Info("hello")

	assert.Regexp(`\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2} hello\n`, buf.String())
	buf.Reset()

	Verbose = true

	Info("goodbye")

	assert.Regexp(`\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2} goodbye\n`, buf.String())

}
