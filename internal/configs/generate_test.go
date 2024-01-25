package configs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	err := Generate("./test-data")
	assert.NoError(t, err)

	data, err := os.ReadFile("./test-data/config_gen.go")
	assert.NoError(t, err)

	refData, err := os.ReadFile("./test-data/ref/config_gen.go")
	assert.Equal(t, string(refData), string(data))
}
