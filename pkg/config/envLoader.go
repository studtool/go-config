package config

import (
	"os"
	"strconv"
	"sync"
)

type EnvLoader struct {
	mutex     *sync.Mutex
	variables map[string]interface{}
}

func NewEnvLoader() *EnvLoader {
	return &EnvLoader{
		mutex:     &sync.Mutex{},
		variables: make(map[string]interface{}),
	}
}

func (loader *EnvLoader) ParseString(name string) string {
	return loader.parseString(name, "", true)
}

func (loader *EnvLoader) ParseStringDefault(name string, def string) string {
	return loader.parseString(name, def, false)
}

func (loader *EnvLoader) ParseInt(name string) int {
	return loader.parseInt(name, 0, true)
}

func (loader *EnvLoader) ParseIntDefault(name string, def int) int {
	return loader.parseInt(name, def, false)
}

func (loader *EnvLoader) ParseFlag(name string) bool {
	return loader.parseFlag(name, false, true)
}

func (loader *EnvLoader) ParseFlagDefault(name string, def bool) bool {
	return loader.parseFlag(name, def, false)
}

func (loader *EnvLoader) parseString(
	name string, defVal string, isRequired bool,
) string {
	val := loader.parseValue(
		name,
		func(value string) (interface{}, error) {
			return value, nil
		},
		"[STRING]",
		defVal,
		isRequired,
	)
	return val.(string)
}

func (loader *EnvLoader) parseInt(
	name string, defVal int, isRequired bool,
) int {
	val := loader.parseValue(
		name,
		func(value string) (interface{}, error) {
			return strconv.Atoi(value)
		},
		"[INTEGER]",
		defVal,
		isRequired,
	)
	return val.(int)
}

func (loader *EnvLoader) parseFlag(
	name string, defVal bool, isRequired bool,
) bool {
	val := loader.parseValue(
		name,
		func(value string) (interface{}, error) {
			return strconv.ParseBool(value)
		},
		"[BOOLEAN]",
		defVal,
		isRequired,
	)
	return val.(bool)
}

type converter func(value string) (interface{}, error)

func (loader *EnvLoader) parseValue(
	name string,
	converter converter, format string,
	defVal interface{}, isRequired bool,
) interface{} {
	loader.mutex.Lock()
	defer loader.mutex.Unlock()

	if v, ok := loader.variables[name]; ok {
		return v
	}

	val := defVal

	strVal := os.Getenv(name)
	if strVal == "" {
		if isRequired {
			panicNotSet(name)
		}
	} else {
		var err error
		if val, err = converter(strVal); err != nil {
			panicInvalidFormat(name, format)
		}
	}

	loader.variables[name] = val
	return val
}
