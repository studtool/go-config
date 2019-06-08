package config

import (
	"testing"

	"github.com/franela/goblin"

	"fmt"
	"os"
)

func TestEnvLoader_Interface(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("EnvLoader", func() {
		g.It("implements Loader interface", func() {
			var _ Loader = NewEnvLoader()
		})
	})
}

func TestEnvLoader_ParseString(t *testing.T) {
	loader := NewEnvLoader()

	checkParser(
		t,
		func(name string) interface{} {
			return loader.ParseString(name)
		},
		"ParseString()",
		"value_1",
		"value_2",
		nil,
	)
}

func TestEnvLoader_ParseStringDefault(t *testing.T) {
	loader := NewEnvLoader()

	checkParserDefault(
		t,
		func(name string, defaultValue interface{}) interface{} {
			return loader.ParseStringDefault(name, defaultValue.(string))
		},
		"ParseStringDefault()",
		"value_1",
		"value_2",
		"default_value",
		nil,
	)
}

const (
	commonInvalidValue = "string"
)

func TestEnvLoader_ParseInt(t *testing.T) {
	loader := NewEnvLoader()
	invalidValue := commonInvalidValue

	checkParser(
		t,
		func(name string) interface{} {
			return loader.ParseInt(name)
		},
		"ParseInt()",
		1,
		2,
		&invalidValue,
	)
}

func TestEnvLoader_ParseIntDefault(t *testing.T) {
	loader := NewEnvLoader()
	invalidValue := commonInvalidValue

	checkParserDefault(
		t,
		func(name string, defaultValue interface{}) interface{} {
			return loader.ParseIntDefault(name, defaultValue.(int))
		},
		"ParseIntDefault()",
		1,
		2,
		0,
		&invalidValue,
	)
}

func TestEnvLoader_ParseFlag(t *testing.T) {
	loader := NewEnvLoader()
	invalidValue := commonInvalidValue

	checkParser(
		t,
		func(name string) interface{} {
			return loader.ParseFlag(name)
		},
		"ParseInt()",
		true,
		false,
		&invalidValue,
	)
}

func TestEnvLoader_ParseFlagDefault(t *testing.T) {
	loader := NewEnvLoader()
	invalidValue := commonInvalidValue

	checkParserDefault(
		t,
		func(name string, defaultValue interface{}) interface{} {
			return loader.ParseFlagDefault(name, defaultValue.(bool))
		},
		"ParseFlagDefault()",
		true,
		false,
		false,
		&invalidValue,
	)
}

func makeEnvVarName() string {
	name := "ENV_VAR"
	for os.Getenv(name) != "" {
		name += "1"
	}
	return name
}

func checkPanicHappens(f func()) (panicHappened bool) {
	panicHappened = false

	defer func() {
		if r := recover(); r != nil {
			panicHappened = true
		}
	}()

	f()
	return
}

func checkParser(
	t *testing.T,
	parser func(string) interface{},
	parserName string,
	firstValue interface{},
	secondValue interface{},
	invalidValue *string,
) {
	g := goblin.Goblin(t)
	g.Describe("EnvLoader."+parserName, func() {
		checkParserCommon(
			g,
			parser,
			firstValue,
			secondValue,
			invalidValue,
		)

		g.It("panics if variable is not set", func() {
			name := makeEnvVarName()

			g.Assert(checkPanicHappens(func() {
				parser(name)
			})).IsTrue()
		})
	})
}

func checkParserDefault(
	t *testing.T,
	parser func(string, interface{}) interface{},
	parserName string,
	firstValue interface{},
	secondValue interface{},
	defaultValue interface{},
	invalidValue *string,
) {
	g := goblin.Goblin(t)
	g.Describe("EnvLoader."+parserName, func() {
		checkParserCommon(
			g,
			func(name string) interface{} {
				return parser(name, defaultValue)
			},
			firstValue,
			secondValue,
			invalidValue,
		)

		g.It("returns default value if variable is not set", func() {
			name := makeEnvVarName()

			g.Assert(parser(name, defaultValue)).Equal(defaultValue)
		})
	})
}

func checkParserCommon(
	g *goblin.G,
	parser func(string) interface{},
	firstValue interface{},
	secondValue interface{},
	invalidValue *string,
) {
	g.It("returns value if variable is set", func() {
		name := makeEnvVarName()

		strValue := fmt.Sprintf("%v", firstValue)
		if err := os.Setenv(name, strValue); err != nil {
			panic(err)
		}

		g.Assert(parser(name)).Equal(firstValue)
	})
	g.It("returns previously parsed value if called second time", func() {
		name := makeEnvVarName()

		strValue := fmt.Sprintf("%v", firstValue)
		if err := os.Setenv(name, strValue); err != nil {
			panic(err)
		}

		_ = parser(name)

		strValue = fmt.Sprintf("%v", secondValue)
		if err := os.Setenv(name, strValue); err != nil {
			panic(err)
		}

		g.Assert(parser(name)).Equal(firstValue)
	})
	if invalidValue != nil {
		g.It("panics if value is is invalid", func() {
			name := makeEnvVarName()

			if err := os.Setenv(name, *invalidValue); err != nil {
				panic(err)
			}

			g.Assert(checkPanicHappens(func() {
				parser(name)
			})).IsTrue()
		})
	}
}
