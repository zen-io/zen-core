package target

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	environs "github.com/zen-io/zen-core/environments"
	"github.com/zen-io/zen-core/utils"
	"golang.org/x/exp/slices"
)

type TargetCreator interface {
	GetTargets(*TargetConfigContext) ([]*TargetBuilder, error)
}

type TargetCreatorMap map[string]TargetCreator

type TargetBuilderScript struct {
	Alias         []string
	Deps          []string
	Env           map[string]string
	PassEnv       []string
	SecretEnv     map[string]string
	PassSecretEnv []string
	Pre           func(target *Target, runCtx *RuntimeContext) error
	Post          func(target *Target, runCtx *RuntimeContext) error
	Run           func(target *Target, runCtx *RuntimeContext) error
	CheckCache    func(target *Target) (bool, error)
	TransformOut  func(target *Target, o string) (string, bool)
	Outs          []string
	Local         bool
}

type TargetBuilder struct {
	Name                 string                           `mapstructure:"name"`
	Description          string                           `mapstructure:"desc" zen:"yes"`
	Labels               []string                         `mapstructure:"labels" zen:"yes"`
	Hashes               []string                         `mapstructure:"hashes" zen:"yes"`
	Visibility           []string                         `mapstructure:"visibility" zen:"yes"`
	Local                bool                             `mapstructure:"local"`
	ExternalPath         *string                          `mapstructure:"external_path"`
	Srcs                 map[string][]string              `mapstructure:"srcs"`
	Outs                 []string                         `mapstructure:"outs"`
	Tools                map[string]string                `mapstructure:"tools" zen:"yes"`
	Binary               bool                             `mapstructure:"binary"`
	Environments         map[string]*environs.Environment `mapstructure:"environments"`
	Deps                 []string                         `mapstructure:"deps"`
	Env                  map[string]string                `mapstructure:"env"`
	PassEnv              []string                         `mapstructure:"pass_env"`
	SecretEnv            map[string]string                `mapstructure:"secret_env"`
	PassSecretEnv        []string                         `mapstructure:"pass_secret_env"`
	Scripts              map[string]*TargetBuilderScript  `mapstructure:"scripts"`
	NoCacheInterpolation bool
	*QualifiedTargetName

	_original_path string
}

func (tb *TargetBuilder) ToTarget(script string, baseEnv map[string]string) *Target {
	// Qn is valid here and script is not empty, so this won't throw an error
	fqn, _ := NewFqnFromStr(fmt.Sprintf("%s:%s", tb.Qn(), script))
	scriptBuilder := tb.Scripts[script]

	for _, e := range scriptBuilder.PassEnv {
		scriptBuilder.Env[e] = os.Getenv(e)
	}

	t := &Target{
		QualifiedTargetName: fqn,
		Labels:              tb.Labels,
		Hashes:              tb.Hashes,
		Tools:               tb.Tools,
		Env:                 utils.MergeMaps(baseEnv, scriptBuilder.Env),
		shouldInterpolate:   !tb.NoCacheInterpolation,
	}

	if script == "build" {
		t.Srcs = tb.Srcs
		t.Outs = tb.Outs
	} else {
		t.Srcs = map[string][]string{"_src": {fqn.BuildFqn()}}
		t.Outs = scriptBuilder.Outs
	}

	return t
}

func (tb *TargetBuilder) EnsureValidTarget() (err error) {
	if tb.Name == "" {
		return fmt.Errorf("Target needs a name")
	}

	tb.Visibility, err = InferArrayRefs(tb.Visibility, tb.Project(), tb.Package(), "build")
	if err != nil {
		return fmt.Errorf("ensuring visibility: %w", err)
	}

	buildDeps := []string{}
	for _, dep := range tb.Deps {
		fqn, err := InferFqn(dep, tb.Project(), tb.Package(), "build")
		if err != nil {
			return err
		}

		buildDeps = append(buildDeps, fqn.Fqn())
	}

	// tools needs to happen before deps, because we add references to the deps
	tools := map[string]string{}
	for toolName, toolRef := range tb.Tools {
		if IsTargetReference(toolRef) { // src is a reference
			if toolRefFqn, err := InferFqn(toolRef, tb.Project(), tb.Package(), "build"); err != nil {
				return fmt.Errorf("tool ref %s not valid: %w", toolName, err)
			} else {
				tools[toolName] = toolRefFqn.Fqn()
				if !slices.Contains(buildDeps, toolRefFqn.Fqn()) {
					buildDeps = append(buildDeps, toolRefFqn.Fqn())
				}
			}
		} else {
			tools[toolName] = toolRef
		}
	}
	tb.Tools = tools

	for scriptName, script := range tb.Scripts {
		deps, err := InferArrayRefs(script.Deps, tb.Project(), tb.Package(), scriptName)
		if err != nil {
			return fmt.Errorf("ensuring deps for script %s: %w", scriptName, err)
		}
		if scriptName != "build" {
			deps = append(deps, tb.BuildFqn())
		} else {
			deps = append(buildDeps, deps...)
		}
		sort.Strings(deps)
		script.Deps = deps

		if script.Outs == nil {
			script.Outs = make([]string, 0)
		} else {
			sort.Strings(script.Outs)
		}
		if script.PassSecretEnv == nil {
			script.PassSecretEnv = make([]string, 0)
		}

		if script.PassEnv == nil {
			script.PassEnv = make([]string, 0)
		}

		if script.Env == nil {
			script.Env = make(map[string]string)
		}
		if script.SecretEnv == nil {
			script.SecretEnv = make(map[string]string)
		}

		if script.TransformOut == nil {
			script.TransformOut = func(target *Target, o string) (string, bool) {
				return o, true
			}
		}
	}

	srcs := map[string][]string{}
	for sName, sSrcs := range tb.Srcs {
		srcs[sName] = []string{}

		for _, src := range sSrcs {
			if IsTargetReference(src) { // src is a reference
				if srcRefFqn, err := InferFqn(src, tb.Project(), tb.Package(), "build"); err != nil {
					return fmt.Errorf("src \"%s\" ref format not correct: %w", src, err)
				} else if !slices.Contains(buildDeps, srcRefFqn.Fqn()) {
					return fmt.Errorf("%s is a src but not a dependency, exiting", src)
				} else {
					srcs[sName] = append(srcs[sName], srcRefFqn.Fqn())
				}
			} else {
				srcs[sName] = append(srcs[sName], src)
			}
		}
	}
	tb.Srcs = srcs

	if tb.Scripts["build"] == nil {
		tb.Scripts["build"] = &TargetBuilderScript{}
	}

	if tb.Scripts["build"].Run == nil {
		tb.Scripts["build"].Run = func(target *Target, runCtx *RuntimeContext) error {
			// fmt.Println(target.Srcs)
			// for _, sSrcs := range tb.Srcs {
			// 	for _, src := range sSrcs {
			// 		if err := target.Copy(src, src); err != nil {
			// 			return err
			// 		}
			// 	}
			// }

			return nil
		}
	}

	if len(tb.Environments) > 0 {
		env_names := make([]string, 0)
		for e := range tb.Environments {
			env_names = append(env_names, e)
		}
		sort.Strings(env_names)
		tb.Labels = append(tb.Labels, "environments="+strings.Join(env_names, ","))
	}

	for _, e := range tb.PassEnv {
		tb.Env[e] = os.Getenv(e)
	}

	sort.Strings(tb.Labels)
	sort.Strings(tb.Outs)

	return nil
}

func (tb *TargetBuilder) ExpandEnvironments(envs ...map[string]*environs.Environment) {
	mergedEnvs := map[string]*environs.Environment{}
	for k, v := range tb.Environments {
		envsToMerge := []*environs.Environment{}
		for _, e := range envs {
			if val, ok := e[k]; ok {
				envsToMerge = append(envsToMerge, val)
			}
		}
		envsToMerge = append(envsToMerge, v)

		mergedEnvs[k] = environs.MergeEnvironments(envsToMerge...)
	}

	tb.Environments = mergedEnvs
}

func (target *Target) ExpandTools(getRefOuts func(ref string) (map[string]string, error)) error {
	tools := map[string]string{}

	tool_keys := make([]string, 0)
	for toolName := range target.Tools {
		tool_keys = append(tool_keys, toolName)
	}
	sort.Strings(tool_keys)

	for _, toolName := range tool_keys {
		tool := target.Tools[toolName]
		if IsTargetReference(tool) {
			outs, err := getRefOuts(tool)
			if err != nil {
				return err
			}

			for _, v := range outs {
				tools[toolName] = v
				break
			}
		} else {
			tools[toolName] = tool
		}

		if val, ok := target.Env["PATH"]; ok {
			target.Env["PATH"] = strings.Join([]string{filepath.Dir(tools[toolName]), val}, ":")
		} else {
			target.Env["PATH"] = filepath.Dir(tools[toolName])
		}
		target.Env["ZEN_TOOL_"+toolName] = tools[toolName]
	}

	target.Tools = tools

	return nil
}

func (target *Target) InterpolateMyself() error {
	srcs := map[string][]string{}
	for sName, sSrcs := range target.Srcs {
		srcs[sName] = make([]string, 0)

		for _, src := range sSrcs {
			if interpolatedSrc, err := target.Interpolate(src); err != nil {
				return err
			} else {
				srcs[sName] = append(srcs[sName], interpolatedSrc)
			}
		}
	}

	target.Srcs = srcs

	outs := make([]string, 0)
	for _, o := range target.Outs {
		if interpolatedOut, err := target.Interpolate(o); err != nil {
			return err
		} else {
			outs = append(outs, interpolatedOut)
		}
	}
	target.Outs = outs

	tools := map[string]string{}
	for toolName, toolValue := range target.Tools {
		if interpolatedTool, err := target.Interpolate(toolValue); err != nil {
			return err
		} else {
			tools[toolName] = interpolatedTool
		}
	}
	target.Tools = tools

	return nil
}

func (tb *TargetBuilder) Interpolate(text string) (string, error) {
	return utils.Interpolate(text, tb.Env)
}
