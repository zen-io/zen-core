package target

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	environs "github.com/zen-io/zen-core/environments"
	"github.com/zen-io/zen-core/utils"

	out_mgr "github.com/tiagoposse/go-tasklist-out"

	"golang.org/x/exp/slices"
)

type TargetOption func(*Target) error

type TargetCreator interface {
	GetTargets(*TargetConfigContext) ([]*Target, error)
}

type TargetCreatorMap map[string]TargetCreator

type TargetScript struct {
	Alias      []string
	Deps       []string
	Env        map[string]string
	Pre        func(target *Target, runCtx *RuntimeContext) error
	Post       func(target *Target, runCtx *RuntimeContext) error
	Run        func(target *Target, runCtx *RuntimeContext) error
	CheckCache func(target *Target) (bool, error)
}

type Target struct {
	Name         string
	Srcs         map[string][]string
	Outs         []string
	TransformOut func(target *Target, o string) (string, bool)
	Labels       []string
	Hashes       []string
	Visibility   []string
	Tools        map[string]string
	Environments map[string]*environs.Environment
	Env          map[string]string
	PassEnv      []string
	SecretEnv    []string
	Local        bool
	Description  string
	Scripts      map[string]*TargetScript
	Binary       bool
	External     bool
	Clean        bool

	NoCacheInterpolation bool
	flattenOuts     bool

	// This will be filled up by the engine
	*QualifiedTargetName
	_original_path string

	out_mgr.TaskLogger
	Cwd string
}

func NewTarget(name string, opts ...TargetOption) *Target {
	target := &Target{
		Name:            name,
		Srcs:            map[string][]string{},
		Outs:            []string{},
		Labels:          []string{},
		Visibility:      []string{},
		Tools:           map[string]string{},
		SecretEnv:       []string{},
		Env:             map[string]string{},
		Environments:    map[string]*environs.Environment{},
		PassEnv:         make([]string, 0),
		Local:           true,
		Description:     "",
		Binary:          false,
		NoCacheInterpolation: false,
		Scripts: map[string]*TargetScript{
			"build": {},
		},
		Clean:          false,
		_original_path: "",
		Cwd:            "",
	}

	target.Local = true

	for _, opt := range opts {
		opt(target)
	}

	return target
}

func (target *Target) EnsureValidTarget() error {
	buildDeps := []string{}

	// tools needs to happen before deps, because we add references to the deps
	tools := map[string]string{}
	for toolName, toolRef := range target.Tools {
		if IsTargetReference(toolRef) { // src is a reference
			if toolRefFqn, err := InferFqn(toolRef, target.Project(), target.Package(), "build"); err != nil {
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
	target.Tools = tools

	for scriptName, script := range target.Scripts {
		deps := []string{}
		for _, dep := range script.Deps {
			if depRefFqn, err := InferFqn(dep, target.Project(), target.Package(), scriptName); err != nil {
				return fmt.Errorf("%s deps \"%s\" format not correct: %w", scriptName, dep, err)
			} else {
				deps = append(deps, depRefFqn.Fqn())
			}
		}

		if scriptName == "build" {
			buildDeps = append(buildDeps, deps...)
			script.Deps = buildDeps
			sort.Strings(script.Deps)
		} else {
			deps = append(deps, target.BuildFqn())
			script.Deps = deps
			sort.Strings(script.Deps)
		}
	}

	srcs := map[string][]string{}
	for sName, sSrcs := range target.Srcs {
		srcs[sName] = []string{}

		for _, src := range sSrcs {
			if IsTargetReference(src) { // src is a reference
				if srcRefFqn, err := InferFqn(src, target.Project(), target.Package(), "build"); err != nil {
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
	target.Srcs = srcs

	visibility := []string{}
	for _, vis := range target.Visibility {
		if IsTargetReference(vis) { // src is a reference
			if visRefFqn, err := InferFqn(vis, target.Project(), target.Package(), "build"); err != nil {
				return fmt.Errorf("visibility ref %s not valid: %w", vis, err)
			} else {
				visibility = append(visibility, visRefFqn.Fqn())
			}
		} else {
			visibility = append(visibility, vis)
		}
	}

	if target.Env == nil {
		target.Env = map[string]string{}
	}

	target.Visibility = visibility

	if target.Scripts["build"].Run == nil {
		target.Scripts["build"].Run = func(target *Target, runCtx *RuntimeContext) error {
			for _, sSrcs := range target.Srcs {
				for _, src := range sSrcs {
					if err := target.Copy(src, src); err != nil {
						return err
					}
				}
			}

			return nil
		}
	}

	sort.Strings(target.Labels)

	if target.TransformOut == nil {
		target.TransformOut = func(target *Target, o string) (string, bool) {
			return o, true
		}
	}
	return nil
}

func (t *Target) Fqn() string {
	return t.Qn()
}

func (t *Target) Exec(command []string, errorMsg string) error {
	t.Debugln(strings.Join(command, " "))
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = t.Cwd
	cmd.Env = t.GetEnvironmentVariablesList()

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %w (%s)", errorMsg, err, out)
	}
	return nil
}

func (t *Target) Copy(from, to string, customVars ...map[string]string) error {
	if t.ShouldInterpolate() {
		if err := utils.CopyWithInterpolate(from, to, append([]map[string]string{t.EnvVars()}, customVars...)...); err != nil {
			return fmt.Errorf("interpolating while copying from %s to %s: %w", from, to, err)
		}
	} else if from != to {
		return utils.CopyFile(from, to)
	}

	return nil
}
