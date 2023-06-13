package target

import (
	"io/ioutil"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

var srcs = map[string][]string{
	"hello": {"hello1", "hello2"},
	"bye":   {"bye*"},
	"ref":   {"//project/path/to/ref/pkg:ref"},
}

func MockTarget(t *testing.T) *Target {
	target := NewTarget(
		"name",
		WithSrcs(srcs),
	)

	target.SetOriginalPath(t.TempDir())
	target.SetFqn("project", "path/to/pkg")

	// CreateTestFiles(root string, files []string)
	return target
}

func MockTargetFull(t *testing.T) *Target {
	target := MockTarget(t)

	for _, srcs := range target.Srcs {
		for _, src := range srcs {
			if err := ioutil.WriteFile(target.Path(), []byte{}, os.ModePerm); err != nil {
				t.Errorf("cant create src: %s\n", src)
				return nil
			}
		}
	}

	return target
}

func TestNewTarget(t *testing.T) {
	target := NewTarget(
		"name",
		WithSrcs(srcs),
	)

	root := t.TempDir()
	target.SetOriginalPath(root)
	target.SetFqn("project", "path/to/pkg")

	assert.Equal(t, target.Project(), "project")
	assert.Equal(t, target.Package(), "path/to/pkg")
	assert.Equal(t, target.Name, "name")
	assert.DeepEqual(t, target.Srcs, srcs)
	assert.Equal(t, target.Path(), root)

}

// import (
// 	"fmt"
// 	"github.com/zen-io/zen-core/utils"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"testing"
// )

// var target = &Target{
// 	Name: "test",
// 	Srcs: map[string][]string{
// 		"base": {"hello*"},
// 		"test": {"bye*"},
// 	},
// }

// func CreateMappings(root string) error {
// 	if err := ioutil.WriteFile(filepath.Join(root, "hello1"), []byte{}, os.ModePerm); err != nil {
// 		return err
// 	}
// 	if err := ioutil.WriteFile(filepath.Join(root, "hello2"), []byte{}, os.ModePerm); err != nil {
// 		return err
// 	}
// 	if err := ioutil.WriteFile(filepath.Join(root, "bye"), []byte{}, os.ModePerm); err != nil {
// 		return err
// 	}

// 	var expectedMappings = map[string]map[string]string{
// 		"base": {
// 			"hello1": "hello1",
// 		},
// 		"test": {
// 			"test1": "test1",
// 		},
// 	}

// 	return nil
// }

// func getRefOuts(ref string) (map[string]string, error) {
// 	return make(map[string]string), nil
// }

// func TestExpandSrcs(t *testing.T) {
// 	root := t.TempDir()
// 	if err := CreateTempFiles(root); err != nil {
// 		t.Error(err)
// 	}

// 	target.SetOriginalPath(root)
// 	mappings, err := target.ExpandSrcs(getRefOuts)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if fmt.Sprint(mappings) != fmt.Sprint(expectedMappings) {
// 		t.Error("mappings are not the same")
// 		utils.PrettyPrint(mappings)
// 		fmt.Println()
// 		utils.PrettyPrint(expectedMappings)
// 	}
// }
