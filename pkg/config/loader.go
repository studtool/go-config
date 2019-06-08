package config

import (
	"fmt"
)

type Loader interface {
	ParseString(name string) string
	ParseStringDefault(name string, def string) string

	ParseInt(name string) int
	ParseIntDefault(name string, def int) int

	ParseFlag(name string) bool
	ParseFlagDefault(name string, def bool) bool
}

func panicNotSet(name string) {
	panic(fmt.Sprintf("'%s' is required", name))
}

func panicInvalidFormat(name string, pattern string) {
	panic(fmt.Sprintf("'%s' format error; pattern - '%s'", name, pattern))
}
