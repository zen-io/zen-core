package target

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"
	environs "github.com/zen-io/zen-core/environments"
	"github.com/zen-io/zen-core/utils"
)

type RuntimeContext struct {
	Context         context.Context
	DryRun          bool
	Env             string
	Tag             string
	WithDeps        bool
	UseEnvironments bool
}

type TargetConfigContext struct {
	KnownToolchains map[string]string
	Variables       map[string]string
	Environments    map[string]*environs.Environment
}

func (tcc *TargetConfigContext) Interpolate(text string, custom ...map[string]string) (string, error) {
	return utils.Interpolate(text, utils.MergeMaps(append([]map[string]string{tcc.Variables}, custom...)...))
}

func (tcc *TargetConfigContext) ResolveToolchain(provided *string, name string, targetTools map[string]string) (tool string, err error) {
	if provided != nil {
		tool = *provided
	} else if val, ok := tcc.KnownToolchains[name]; ok {
		tool = val
	} else if _, ok := targetTools[name]; !ok {
		err = fmt.Errorf("%s toolchain is not configured", name)
	}

	return
}

func NewRuntimeContext(flags *pflag.FlagSet) *RuntimeContext {
	var env, tag string
	var dryRun bool

	env, _ = flags.GetString("env")
	tag, _ = flags.GetString("tag")
	dryRun, _ = flags.GetBool("dry-run")
	noDeps, _ := flags.GetBool("no-deps")

	return &RuntimeContext{
		Context: context.Background(),
		Env:     env,
		Tag:     tag,
		DryRun:  dryRun,
		WithDeps: !noDeps,
	}
}

func (target *Target) GetEnvironmentVariablesList(additionalVars ...map[string]string) []string {
	envVarList := []string{}
	for k, v := range utils.MergeMaps(append([]map[string]string{target.Env}, additionalVars...)...) {
		if v != "" && k != "ENV" { // ENV is a special variable in sh, that causes it to execute a script. We need to consider renaming.
			envVarList = append(envVarList, fmt.Sprintf("%s=%s", k, v))
		}
	}

	return envVarList
}
