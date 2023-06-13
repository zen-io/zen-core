package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func CreateTestFiles(root string, files []string) (map[string]string, error) {
	mappings := map[string]string{}

	for _, f := range files {
		mappings[f] = filepath.Join(root, f)
		if err := ioutil.WriteFile(mappings[f], []byte{}, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return mappings, nil
}

func MockIsRef(path string) bool {
	return strings.HasPrefix(path, "ref_")
}

func MockGetRefMap(refMap map[string]string) func(ref string) (map[string]string, error) {
	return func(ref string) (map[string]string, error) {
		return refMap, nil
	}
}

func TestGlobPath(t *testing.T) {
	root := t.TempDir()
	expectedMappings, err := CreateTestFiles(root, []string{"hello1", "hello2"})
	assert.NilError(t, err)

	mappings, err := GlobPath(
		root,
		"hello*",
	)
	assert.NilError(t, err)
	assert.DeepEqual(t, mappings, expectedMappings)
}

func TestReadExclusionFile(t *testing.T) {
	root := t.TempDir()
	exclusionFilePath := filepath.Join(root, "exclusion")
	err := ioutil.WriteFile(exclusionFilePath, []byte("test\ntest2"), os.ModePerm)
	assert.NilError(t, err)

	exclusions, err := ReadExclusionFile(exclusionFilePath)
	assert.NilError(t, err)

	assert.DeepEqual(t, exclusions, []string{"test", "test2"})
}
