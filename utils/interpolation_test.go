package utils

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestInterpolateMyself(t *testing.T) {
	original := map[string]string{
		"TEST":              "{TEST_REPLACE}",
		"TEST_REPLACE_DEEP": "{TEST}",
		"TEST_REPLACE":      "1",
	}

	expected := map[string]string{
		"TEST":              "1",
		"TEST_REPLACE_DEEP": "1",
		"TEST_REPLACE":      "1",
	}

	received, err := InterpolateMapWithItself(original)
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, received)

	errored := map[string]string{
		"TEST":         "{TEST_NO_EXIST}",
		"TEST_REPLACE": "1",
	}
	_, err = InterpolateMapWithItself(errored)
	assert.Error(t, err, "{TEST_NO_EXIST} is not a valid interpolation var")
}
