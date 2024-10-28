package gttools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddListWithPrefix(t *testing.T) {
	assert := assert.New(t)

	args := []string{}
	args = addStringList(args, []string{"a", "b", "c"}, addSettings{prefix: "--mysetting=", listSeparator: ","})

	assert.Equal([]string{"--mysetting=a,b,c"}, args)
}

func TestAddListWithPrepend(t *testing.T) {
	assert := assert.New(t)

	args := []string{}
	args = addStringList(args, []string{"a", "b", "c"}, addSettings{prependElements: []string{"--mysetting"}, listSeparator: ","})

	assert.Equal([]string{"--mysetting", "a,b,c"}, args)
}

func TestAddListSeparateWithPrefix(t *testing.T) {
	assert := assert.New(t)

	args := []string{}
	args = addStringList(args, []string{"a", "b", "c"}, addSettings{prefix: "--mysetting=", handleEachListItemSeparately: true})

	assert.Equal([]string{"--mysetting=a", "--mysetting=b", "--mysetting=c"}, args)
}

func TestAddListSeparateWitPrepend(t *testing.T) {
	assert := assert.New(t)

	args := []string{}
	args = addStringList(args, []string{"a", "b", "c"}, addSettings{prependElements: []string{"--mysetting"}, handleEachListItemSeparately: true})

	assert.Equal([]string{"--mysetting", "a", "--mysetting", "b", "--mysetting", "c"}, args)
}
