package mock

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	zen_target "github.com/zen-io/zen-core/target"

	"gotest.tools/v3/assert"
)

func MockBasicTarget(t *testing.T) *zen_target.Target {
	target := zen_target.NewTarget(
		"basic",
		zen_target.WithSrcs(MockSrcs["basic"].Srcs),
		zen_target.WithOuts(MockSrcs["basic"].Outs),
	)

	target.SetOriginalPath(t.TempDir())
	target.SetFqn("project", "path/to/pkg")

	return target
}

func MockBasicTargetFull(t *testing.T) *zen_target.Target {
	target := MockBasicTarget(t)

	return target
}

func MockComplexTarget(t *testing.T) *zen_target.Target {
	target := zen_target.NewTarget(
		"complex",
		zen_target.WithSrcs(MockSrcs["complex"].Srcs),
		zen_target.WithOuts(MockSrcs["complex"].Outs),
	)

	target.SetOriginalPath(t.TempDir())
	target.SetFqn("project", "path/to/pkg")

	return target
}

func MockComplexTargetFull(t *testing.T) *zen_target.Target {
	target := MockComplexTarget(t)

	return target
}

type MockSrcsDef struct {
	Srcs         map[string][]string
	SrcsMappings map[string]map[string]string
	Outs         []string
	OutsMappings map[string]string
}

var MockSrcs = map[string]MockSrcsDef{
	"basic": {
		Srcs: map[string][]string{
			"hello": {"hello1", "hello2"},
			"bye":   {"bye*"},
		},
		SrcsMappings: map[string]map[string]string{
			"hello": {
				"hello1": "hello1",
				"hello2": "hello2",
			},
			"bye": {
				"bye1": "bye1",
			},
		},
		Outs: []string{"hello1", "bye1"},
		OutsMappings: map[string]string{
			"hello1": "hello1",
			"bye1":   "bye1",
		},
	},
	"complex": {
		Srcs: map[string][]string{
			"test": {"test1", "//project/path/to/pkg:complex"},
		},
		SrcsMappings: map[string]map[string]string{
			"test": {
				"test1":  "test1",
				"hello1": "hello1",
				"bye1":   "bye1",
			},
		},
		Outs: []string{"*"},
		OutsMappings: map[string]string{
			"hello1": "hello1",
			"bye1":   "bye1",
			"test1":  "test1",
		},
	},
}

func (msd MockSrcsDef) ExpandSrcsMappings(root string) map[string]map[string]string {
	mocked := map[string]map[string]string{}

	for cat, val := range msd.SrcsMappings {
		mocked[cat] = make(map[string]string)
		for k, v := range val {
			mocked[cat][k] = filepath.Join(root, v)
		}
	}
	return mocked
}

func (msd MockSrcsDef) ExpandOutsMappings(root string) map[string]string {
	mocked := map[string]string{}

	for k, v := range msd.OutsMappings {
		mocked[k] = filepath.Join(root, v)
	}
	return mocked
}

var ComplexSrcs = map[string][]string{
	"hello": {"hello1", "hello2"},
	"bye":   {"bye*"},
	"ref":   {"//project/path/to/ref/pkg:ref"},
}

func CreateFiles(t *testing.T, files []string) {
	for _, src := range files {
		err := ioutil.WriteFile(src, []byte{}, os.ModePerm)
		assert.NilError(t, err)
	}
}

func FileMapToSlice(t *testing.T, m map[string]map[string]string) []string {
	sl := []string{}
	for _, val := range m {
		for _, v := range val {
			sl = append(sl, v)
		}
	}
	return sl
}
