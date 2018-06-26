package generator

import (
	"github.com/kreuzwerker/awsu/generator/manual"
	"github.com/kreuzwerker/awsu/generator/yubikey"
)

const (
	Manual  Name = "manual"
	Yubikey Name = "yubikey"
)

var Generators = map[Name]Generator{
	Yubikey: yubikey.Generate,
	Manual:  manual.Generate,
}

type Name string

type Generator func(serial string) (string, error)
