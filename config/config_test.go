package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectAndGet(t *testing.T) {

	assert := assert.New(t)

	cfg := Detect()

	assert.Nil(cfg)

	os.Setenv("AWSU_WORKSPACE_DEFAULT", "kreuzwerker")
	os.Setenv("AWSU_WORKSPACE_STAGING", "bar")
	os.Setenv("AWSU_WORKSPACE_PRODUCTION", "foo")

	cfg = Detect()

	assert.NotNil(cfg)

	var tt = []struct {
		in  string
		out string
	}{
		{"default", "kreuzwerker"},
		{"development", ""},
		{"staging", "bar"},
		{"production", "foo"},
	}

	for _, e := range tt {

		out, err := cfg.Get(e.in)

		assert.EqualValues(e.out, out)

		if e.out == "" {
			assert.Error(err)
		}

	}

}

func TestLoadAndGet(t *testing.T) {

	assert := assert.New(t)

	cfg, err := Load("no-such-path")

	assert.Nil(cfg)
	assert.Error(err)

	cfg, err = Load(".")

	assert.NotNil(cfg)
	assert.NoError(err)

	var tt = []struct {
		in  string
		out string
	}{
		{"default", "kreuzwerker"},
		{"development", ""},
		{"staging", "bar"},
		{"production", "foo"},
	}

	for _, e := range tt {

		out, err := cfg.Get(e.in)

		assert.EqualValues(e.out, out)

		if e.out == "" {
			assert.Error(err)
		}

	}

}
