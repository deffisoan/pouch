package config

import (
	"fmt"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestConfigValidate(t *testing.T) (t *testing.T) {
	assert := assert.New(t)
	origin := map[string]interface{}{
		"a": "a",
		"b": "b",
		"c": "c",
		"iter1": map[string]interface{}{
			"i1": "i1",
			"i2": "i2",
		},
		"iter11": map[string]interface{}{
			"ii1": map[string]interface{}{
				"iii1": "iii1",
				"iii2": "iii2",
			},
		},
	}

	expect := map[string]interface{}{
		"a":    "a",
		"b":    "b",
		"c":    "c",
		"i1":   "i1",
		"i2":   "i2",
		"iii1": "iii1",
		"iii2": "iii2",
	}

	config := make(map[string]interface{})
	iterateConfig(origin, config)
	assert.Equal(config, expect)

	// test nil map will not cause panic
	config = make(map[string]interface{})
	iterateConfig(nil, config)
	assert.Equal(config, map[string]interface{}{})
}

func TestConfigValidate(t *testing.T) {
	assert := assert.New(t)
	config := Config{Labels: []string{"a=b", "c=d"}}
	origin := config.Validate()
	assert.Nil(origin)
	config = Config{Labels: []string{"a=b", "cd"}}
	origin = config.Validate()
	assert.Equal(origin, fmt.Errorf("daemon label cd must be in format of key=value"))
	config = Config{Labels: []string{"a="}}
	origin = config.Validate()
	assert.Equal(origin, fmt.Errorf("key and value in daemon label a= cannot be empty"))
}

func TestGetConflictConfigurations(t *testing.T) {
	assert := assert.New(t)
	flagSet := pflag.NewFlagSet("d", 1)

	fileFlags := map[string]interface{}{
		"c": "d",
		"e": "a",
	}
	origin := getConflictConfigurations(flagSet, fileFlags)
	assert.Nil(origin)
}

func TestGetUnknownFlags(t *testing.T) {
	assert := assert.New(t)
	fileFlags := map[string]interface{}{
		"a": "a",
		"b": "b",
		"c": "c",
	}

	expect := fmt.Errorf("unknown flags: a, b, c")

	flagSet := pflag.NewFlagSet("test", 0)
	flagSet.String("d", "d", "d")
	// test if it works,can not found
	assert.Equal(expect, getUnknownFlags(flagSet, fileFlags))

	flagSet1 := pflag.NewFlagSet("test1", 0)
	flagSet1.String("a", "a", "a")
	expect = fmt.Errorf("unknown flags: b, c")
	assert.Equal(expect, getUnknownFlags(flagSet1, fileFlags))
}
