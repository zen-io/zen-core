package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func InterpolateMapWithItself(toInterpolate map[string]string) (map[string]string, error) {
	re := regexp.MustCompile(`\{[A-Z\.\_]+\}`)
	var err error
	needsInterpolation := true

	for needsInterpolation {
		needsInterpolation = false

		for key, val := range toInterpolate {
			if re.MatchString(val) {
				needsInterpolation = true
			}

			toInterpolate[key] = re.ReplaceAllStringFunc(val, func(m string) string {
				if val, ok := toInterpolate[m[1:len(m)-1]]; !ok {
					err = fmt.Errorf("%s is not a valid interpolation var", m)
				} else {
					return val
				}

				return ""
			})
		}

	}

	return toInterpolate, err
}

func Interpolate(text string, vars map[string]string) (string, error) {
	re := regexp.MustCompile(`\{[A-Za-z\.\_]+\}`)
	var err error
	var needsInterpolation bool

	retString := text
	for needsInterpolation = true; needsInterpolation; needsInterpolation = re.MatchString(retString) {
		retString = re.ReplaceAllStringFunc(retString, func(m string) string {
			if val, ok := vars[m[1:len(m)-1]]; !ok {
				err = fmt.Errorf("%s is not a valid interpolation var", m)
			} else {
				return val
			}

			return ""
		})
	}

	return retString, err
}

func InterpolateSetVars(text string, vars map[string]string) string {
	keys := []string{}
	for key := range vars {
		keys = append(keys, strings.ToUpper(key))
	}
	if len(keys) == 0 {
		return text
	}

	re := regexp.MustCompile(`\{(?:` + strings.Join(keys, "|") + `)\}`)
	var needsInterpolation bool

	retString := text
	for needsInterpolation = true; needsInterpolation; needsInterpolation = re.MatchString(retString) {
		retString = re.ReplaceAllStringFunc(retString, func(m string) string {
			return vars[strings.ToLower(m)[1:len(m)-1]]
		})
	}

	return retString
}

func InterpolateSlice(texts []string, vars map[string]string) ([]string, error) {
	var err error

	retString := []string{}
	for _, t := range texts {
		interpolatedText, err := Interpolate(t, vars)
		if err != nil {
			return nil, fmt.Errorf("interpolating %s: %w", t, err)
		}

		retString = append(retString, interpolatedText)
	}
	return retString, err
}

func InterpolateMap(m map[string]string, vars map[string]string) (map[string]string, error) {
	ret := map[string]string{}
	for k, v := range m {
		interpolatedKey, err := Interpolate(k, vars)
		if err != nil {
			return nil, fmt.Errorf("interpolating key %s: %w", k, err)
		}

		interpolatedVal, err := Interpolate(v, vars)
		if err != nil {
			return nil, fmt.Errorf("interpolating val %s: %w", v, err)
		}

		ret[interpolatedKey] = interpolatedVal
	}

	return ret, nil
}

func MergeMaps(maps ...map[string]string) map[string]string {
	merged := map[string]string{}

	for _, m := range maps {
		if m == nil {
			continue
		}

		for k, v := range m {
			merged[k] = v
		}
	}

	return merged
}

func CopyWithInterpolate(from, to string, interpolateVars ...map[string]string) error {
	data, err := os.ReadFile(from)
	if err != nil {
		return fmt.Errorf("reading from %s: %w", from, err)
	}

	if CheckFileCanInterpolate(data) {
		interpolatedData, err := Interpolate(string(data), MergeMaps(interpolateVars...))
		if err != nil {
			return fmt.Errorf("interpolating data: %w", err)
		}
		data = []byte(interpolatedData)
	}

	if err := os.MkdirAll(filepath.Dir(to), os.ModePerm); err != nil {
		return fmt.Errorf("opening folder %s: %w", filepath.Dir(to), err)
	}

	if err := os.WriteFile(to, data, 0644); err != nil {
		return fmt.Errorf("writing to %s, %w", to, err)
	}

	return nil
}

func CheckFileCanInterpolate(data []byte) bool {
	// Check if the file is an archive
	if
	// empty file
	len(data) == 0 ||
		// Gzip archive
		(data[0] == 0x1f && data[1] == 0x8b) ||
		// // Rar archive
		(data[0] == 0x52 && data[1] == 0x61 && data[2] == 0x72 && data[3] == 0x21 && data[4] == 0x1a && data[5] == 0x07 && data[6] == 0x00) ||
		// ELF binary
		(data[0] == 0x7f && data[1] == 0x45 && data[2] == 0x4c && data[3] == 0x46) ||
		// DOS binary
		(data[0] == 0x4d && data[1] == 0x5a) ||
		// Mach-O Executable (32 bit)
		(data[0] == 0xFE && data[1] == 0xED && data[2] == 0xFA && data[3] == 0xCE) ||
		// Mach-O Executable (64 bit)
		(data[0] == 0xFE && data[1] == 0xED && data[2] == 0xFA && data[3] == 0xCF) ||
		// Zip archive
		(data[0] == 0x50 && data[1] == 0x4B && data[2] == 0x03 && data[3] == 0x04) {
		return false
	}

	return true
}
