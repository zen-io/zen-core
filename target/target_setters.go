package target

import (
	"fmt"
	"reflect"

	environs "github.com/zen-io/zen-core/environments"
	"golang.org/x/exp/slices"
)

func (t *TargetBuilder) SetOriginalPath(path string) {
	if t._original_path == "" {
		t._original_path = path
	}
}

func (tb *TargetBuilder) SetFqn(project, pkg, name, script string) {
	tb.QualifiedTargetName = NewFqnFromParts(project, pkg, name, script)

	if tb.Env == nil {
		tb.Env = make(map[string]string)
	}

	tb.Env["NAME"] = name
	tb.Env["PROJECT"] = project
	tb.Env["PKG"] = pkg
	tb.Env["RESOURCE"] = fmt.Sprintf("%s/%s", pkg, tb.Name)
}

func ToTarget(in interface{}) *TargetBuilder {
	t := &TargetBuilder{
		Description:          "",
		Labels:               make([]string, 0),
		Hashes:               make([]string, 0),
		Visibility:           make([]string, 0),
		Local:                false,
		ExternalPath:         nil,
		Srcs:                 make(map[string][]string),
		Outs:                 make([]string, 0),
		Tools:                make(map[string]string),
		Binary:               false,
		Environments:         make(map[string]*environs.Environment),
		Deps:                 make([]string, 0),
		Env:                  make(map[string]string),
		PassEnv:              make([]string, 0),
		PassSecretEnv:        make([]string, 0),
		SecretEnv:            make(map[string]string),
		Scripts:              make(map[string]*TargetBuilderScript),
		NoCacheInterpolation: false,
	}

	_ = CopyTargetFields(in, t)

	return t
}

func CopyTargetFields(from, to interface{}) error {
	fromVal := reflect.ValueOf(from)
	toVal := reflect.ValueOf(to)

	if toVal.Kind() != reflect.Ptr || toVal.IsNil() {
		return fmt.Errorf("'to' needs to be a pointer")
	} else {
		toVal = toVal.Elem()
	}

	typeOfFrom := fromVal.Type()

	for i := 0; i < fromVal.NumField(); i++ {
		// Check if zen tag is set to "yes"
		if zenTag := typeOfFrom.Field(i).Tag.Get("zen"); zenTag == "yes" {
			// Check if field exists in struct B
			fieldTo := toVal.FieldByName(typeOfFrom.Field(i).Name)
			if fieldTo.IsValid() && fieldTo.CanSet() {
				if slices.Contains([]reflect.Kind{reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice}, fromVal.Field(i).Kind()) && fromVal.Field(i).IsNil() {
					continue
				}

				fieldTo.Set(fromVal.Field(i))
			}
		}
	}

	return nil
}

func (tb *TargetBuilder) KnownEnvironments() []string {
	envs := make([]string, 0)
	for e := range tb.Environments {
		envs = append(envs, e)
	}

	return envs
}
