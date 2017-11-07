package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
