package target

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zen-io/zen-core/utils"
)

// Checks if text is a reference to a target. This is true if the text starts with // or :
func IsTargetReference(text string) bool {
	return regexp.MustCompile(`^(\/\/|:)`).MatchString(text)
}

func (t *Target) StripCwd(s string) string {
	return strings.TrimPrefix(s, t.Cwd+"/")
}

func (target *Target) Interpolate(text string, custom ...map[string]string) (string, error) {
	interpolateVars := utils.MergeMaps(
		append([]map[string]string{target.Env}, custom...)...,
	)

	return utils.Interpolate(text, interpolateVars)
}

func InferArrayRefs(arr []string, proj, pkg, script string) ([]string, error) {
	result := make([]string, 0)

	for _, item := range arr {
		if IsTargetReference(item) { // src is a reference
			if refFqn, err := InferFqn(item, proj, pkg, script); err != nil {
				return nil, fmt.Errorf("ref %s not valid: %w", item, err)
			} else {
				result = append(result, refFqn.Fqn())
			}
		} else {
			result = append(result, item)
		}
	}
	return result, nil
}
