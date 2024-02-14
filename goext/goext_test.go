package goext

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunWithEnv(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("PRE_VAR_1", "baz")
	os.Setenv("PRE_VAR_2", "hello")
	assertEnvEquals(assert, "PRE_VAR_1", "baz")
	assertEnvEquals(assert, "PRE_VAR_2", "hello")
	assertEnvIsUnset(assert, "VAR_1")
	assertEnvIsUnset(assert, "VAR_2")

	err := RunWithEnv(map[string]string{
		"VAR_1":     "foo",
		"VAR_2":     "bar",
		"PRE_VAR_2": "world",
	}, func() error {
		assertEnvEquals(assert, "PRE_VAR_1", "baz")
		assertEnvEquals(assert, "PRE_VAR_2", "world")
		assertEnvEquals(assert, "VAR_1", "foo")
		assertEnvEquals(assert, "VAR_2", "bar")
		return nil
	})

	assert.NoError(err)
	assertEnvEquals(assert, "PRE_VAR_1", "baz")
	assertEnvEquals(assert, "PRE_VAR_2", "hello")
	assertEnvIsUnset(assert, "VAR_1")
	assertEnvIsUnset(assert, "VAR_2")
}

func assertEnvEquals(assert *assert.Assertions, name string, expected string) {
	actual, ok := os.LookupEnv(name)
	assert.True(ok)
	assert.Equal(expected, actual)
}

func assertEnvIsUnset(assert *assert.Assertions, name string) {
	_, ok := os.LookupEnv(name)
	assert.False(ok)
}
