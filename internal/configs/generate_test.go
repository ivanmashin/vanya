package configs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestGenerate_SingleObj(t *testing.T) {
	rootDir := "./test-data/single-obj"

	err := Generate(rootDir)
	assert.NoError(t, err)

	data, err := os.ReadFile(path.Join(rootDir, "config_gen.go"))
	assert.NoError(t, err)

	refData, err := os.ReadFile(path.Join(rootDir, "ref/config_gen.go"))
	assert.Equal(t, string(refData), string(data))
}

func TestGenerate_MultipleObj(t *testing.T) {
	rootDir := "./test-data/multiple-objs"

	err := Generate(rootDir)
	assert.NoError(t, err)

	data, err := os.ReadFile(path.Join(rootDir, "config_gen.go"))
	assert.NoError(t, err)

	refData, err := os.ReadFile(path.Join(rootDir, "ref/config_gen.go"))
	assert.Equal(t, string(refData), string(data))
}
