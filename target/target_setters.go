package target

import (
	"fmt"
	"os"
	"reflect"

	environs "github.com/zen-io/zen-core/environments"
	"github.com/zen-io/zen-core/utils"
)

func (t *Target) SetOriginalPath(path string) {
	if t._original_path == "" {
		t._original_path = path
	}
}

func (target *Target) SetFqn(project, pkg string) {
	target.QualifiedTargetName = NewFqnFromParts(project, pkg, target.Name, "")

	if target.Env == nil {
		target.Env = make(map[string]string)
	}

	target.Env["NAME"] = target.Name
	target.Env["PROJECT"] = project
	target.Env["PKG"] = pkg
	target.Env["RESOURCE"] = fmt.Sprintf("%s/%s", pkg, target.Name)
}

func (target *Target) SetBuildVariables(vars map[string]string) (err error) {
	passedEnv := map[string]string{}
	for _, e := range append(target.PassEnv, target.SecretEnv...) {
		passedEnv[e] = os.Getenv(e)
	}

	target.Env, err = utils.InterpolateMapWithItself(
		utils.MergeMaps(vars, passedEnv, target.Env),
	)

	return
}

func (target *Target) SetDeployVariables(env string, proj, cli map[string]string) (err error) {
	targetDeployEnv, err := target.Environments[env].EnvVarsForEnv()
	if err != nil {
		return fmt.Errorf("loading environment: %w", err)
	}

	target.Env, err = utils.InterpolateMapWithItself(utils.MergeMaps(proj, cli, target.Env, targetDeployEnv, map[string]string{"DEPLOY_ENV": env, "ENV": env}))

	return
}

func (target *Target) ExpandEnvironments(envs ...map[string]*environs.Environment) {
	mergedEnvs := map[string]*environs.Environment{}
	for k, v := range target.Environments {
		envsToMerge := []*environs.Environment{}
		for _, e := range envs {
			if val, ok := e[k]; ok {
				envsToMerge = append(envsToMerge, val)
			}
		}
		envsToMerge = append(envsToMerge, v)

		mergedEnvs[k] = environs.MergeEnvironments(envsToMerge...)
	}

	target.Environments = mergedEnvs
}

func ToTarget(m TargetCreator) (*Target) {
	t := &Target{}

	_ = CopyTargetFields(&m, t)

	return t
}

func CopyTargetFields(from, to interface{}) error {
	fromVal := reflect.ValueOf(from)
	toVal := reflect.ValueOf(to)

	if fromVal.Kind() != reflect.Ptr || fromVal.IsNil() {
		return fmt.Errorf("'from' needs to be a pointer")
	} else {
		fromVal = fromVal.Elem()
	}
	if toVal.Kind() != reflect.Ptr || toVal.IsNil() {
		return fmt.Errorf("'to' needs to be a pointer")
	} else {
		toVal = toVal.Elem()
	}

	typeOfA := fromVal.Type()

	for i := 0; i < fromVal.NumField(); i++ {
		// Check if zen tag is set to "yes"
		if zenTag := typeOfA.Field(i).Tag.Get("zen"); zenTag == "yes" {
			// Check if field exists in struct B
			fieldB := toVal.FieldByName(typeOfA.Field(i).Name)
			if fieldB.IsValid() && fieldB.CanSet() {
				fieldB.Set(fromVal.Field(i))
			}
		}
	}

	return nil
}
