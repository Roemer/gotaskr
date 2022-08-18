package argparse

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleLongEqual(t *testing.T) {
	assert := assert.New(t)

	argsMap := ParseArgString(strings.Fields("--name=test"))

	assert.Equal(1, len(argsMap))
	assert.Contains(argsMap, "name")
	assert.Equal("test", argsMap["name"])
}

func TestSingleLongSpace(t *testing.T) {
	assert := assert.New(t)

	argsMap := ParseArgString(strings.Fields("--name test"))

	assert.Equal(1, len(argsMap))
	assert.Contains(argsMap, "name")
	assert.Equal("test", argsMap["name"])
}

func TestSingleShort(t *testing.T) {
	assert := assert.New(t)

	argsMap := ParseArgString(strings.Fields("-c"))

	assert.Equal(1, len(argsMap))
	assert.Contains(argsMap, "c")
	assert.Empty(argsMap["c"])
}

func TestSingleShortWithValue(t *testing.T) {
	assert := assert.New(t)

	argsMap := ParseArgString(strings.Fields("-l verbose"))

	assert.Equal(1, len(argsMap))
	assert.Contains(argsMap, "l")
	assert.Equal("verbose", argsMap["l"])
}

func TestCombinedShort(t *testing.T) {
	assert := assert.New(t)

	argsMap := ParseArgString(strings.Fields("-abc"))

	assert.Equal(3, len(argsMap))
	assert.Contains(argsMap, "a")
	assert.Contains(argsMap, "b")
	assert.Contains(argsMap, "c")
	assert.Empty(argsMap["a"])
	assert.Empty(argsMap["b"])
	assert.Empty(argsMap["c"])
}

func TestComplex(t *testing.T) {
	assert := assert.New(t)

	argsMap := ParseArgString(strings.Fields("--loglevel debug -x --path=/root/tmp -cf"))

	assert.Equal(5, len(argsMap))
	assert.Contains(argsMap, "loglevel")
	assert.Contains(argsMap, "path")
	assert.Contains(argsMap, "x")
	assert.Contains(argsMap, "c")
	assert.Contains(argsMap, "f")
	assert.Equal("debug", argsMap["loglevel"])
	assert.Equal("/root/tmp", argsMap["path"])
	assert.Empty(argsMap["x"])
	assert.Empty(argsMap["c"])
	assert.Empty(argsMap["f"])
}
