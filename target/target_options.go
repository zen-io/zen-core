package target

import (
	"fmt"

	environs "github.com/zen-io/zen-core/environments"
)

func WithPublicVisibility() func(*Target) error {
	return func(s *Target) error {
		s.Visibility = []string{"PUBLIC"}
		return nil
	}
}

func WithNoInterpolation() func(*Target) error {
	return func(s *Target) error {
		s.noInterpolation = true
		return nil
	}
}

func WithSrcs(srcs map[string][]string) func(*Target) error {
	return func(self *Target) error {
		for sName, srcValues := range srcs {
			if self.Srcs[sName] == nil {
				self.Srcs[sName] = []string{}
			}

			for _, s := range srcValues {
				if s == "" {
					return fmt.Errorf("source is empty")
				}

				self.Srcs[sName] = append(self.Srcs[sName], s)
			}
		}

		return nil
	}
}

func WithHashes(hashes []string) func(*Target) error {
	return func(s *Target) error {
		s.Hashes = hashes
		return nil
	}
}

func WithTools(tools map[string]string) func(*Target) error {
	return func(s *Target) error {
		s.Tools = tools
		for _, t := range tools {
			if t == "" {
				return fmt.Errorf("empty tool detected")
			}

			s.Scripts["build"].Deps = append(s.Scripts["build"].Deps, t)
		}

		return nil
	}
}

func WithLabels(labels []string) func(*Target) error {
	return func(s *Target) error {
		s.Labels = labels
		return nil
	}
}

func WithDescription(desc string) func(*Target) error {
	return func(s *Target) error {
		s.Description = desc
		return nil
	}
}

func WithFlattenOuts() func(*Target) error {
	return func(s *Target) error {
		s.flattenOuts = true
		return nil
	}
}

func WithVisibility(vis []string) func(*Target) error {
	return func(s *Target) error {
		s.Visibility = vis
		return nil
	}
}

func WithRemoteExecution() func(*Target) error {
	return func(s *Target) error {
		s.Local = false
		return nil
	}
}

func WithOuts(outs []string) func(*Target) error {
	return func(s *Target) error {
		s.Outs = outs
		return nil
	}
}

func WithBinary() func(*Target) error {
	return func(s *Target) error {
		s.Binary = true
		return nil
	}
}

func WithEnvironments(env map[string]*environs.Environment) func(*Target) error {
	return func(s *Target) error {
		s.Environments = env
		return nil
	}
}

func WithSecretEnvVars(env []string) func(*Target) error {
	return func(s *Target) error {
		s.SecretEnv = env
		return nil
	}
}

func WithEnvVars(env map[string]string) func(*Target) error {
	return func(s *Target) error {
		s.Env = env
		return nil
	}
}

func WithPassEnv(passEnv []string) func(*Target) error {
	return func(s *Target) error {
		s.PassEnv = passEnv
		return nil
	}
}

func WithExternalPath(p string) func(*Target) error {
	return func(s *Target) error {
		s.External = true
		s._original_path = p
		return nil
	}
}

func WithTargetScript(name string, script *TargetScript) func(*Target) error {
	if script.Deps == nil {
		script.Deps = []string{}
	}
	if script.Alias == nil {
		script.Alias = []string{}
	}

	return func(s *Target) error {
		s.Scripts[name] = script
		return nil
	}
}
