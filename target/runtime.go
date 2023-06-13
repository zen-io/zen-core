package target

import (
	"fmt"
	"os"

	environs "github.com/zen-io/zen-core/environments"
	"github.com/zen-io/zen-core/utils"

	"github.com/spf13/pflag"
)

type RuntimeContext struct {
	Variables    map[string]string
	Environments map[string]*environs.Environment

	DryRun   bool
	Debug    bool
	Clean    bool
	WithDeps bool
	Env      string
	Tag      string
	Shell    bool
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

func NewRuntimeContext(flags *pflag.FlagSet, envs map[string]*environs.Environment, path, hostOS, hostArch string) *RuntimeContext {
	var env, tag string
	var dryRun, debug, clean, withDeps bool

	env, _ = flags.GetString("env")
	tag, _ = flags.GetString("tag")
	dryRun, _ = flags.GetBool("dry-run")
	clean, _ = flags.GetBool("clean")
	debug, _ = flags.GetBool("debug")
	withDeps, _ = flags.GetBool("with-deps")
	shell, _ := flags.GetBool("shell")

	return &RuntimeContext{
		Env:          env,
		Tag:          tag,
		DryRun:       dryRun,
		Debug:        debug,
		Clean:        clean,
		WithDeps:     withDeps,
		Shell:        shell,
		Environments: envs,
		Variables: map[string]string{
			"USER":            os.Getenv("USER"),
			"HOME":            os.Getenv("HOME"),
			"SHLVL":           "1",
			"PATH":            fmt.Sprintf("%s:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin", path),
			"ENV":             env,
			"TAG":             tag,
			"TARGET.OS":       hostOS,
			"TARGET.ARCH":     hostArch,
			"CONFIG.HOSTOS":   hostOS,
			"CONFIG.HOSTARCH": hostArch,
		},
	}
}

func (target *Target) GetEnvironmentVariablesList(additionalVars ...map[string]string) []string {
	envVarList := []string{}
	for k, v := range utils.MergeMaps(append([]map[string]string{target.EnvVars()}, additionalVars...)...) {
		if v != "" && k != "ENV" { // ENV is a special variable in sh, that causes it to execute a script. We need to consider renaming.
			envVarList = append(envVarList, fmt.Sprintf("%s=%s", k, v))
		}
	}

	return envVarList
}
