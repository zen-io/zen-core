package target

import environs "github.com/zen-io/zen-core/environments"

type BuildFields struct {
	Srcs       []string `mapstructure:"srcs" desc:"Sources for the build"`
	Outs       []string `mapstructure:"outs" desc:"Outs for the build"`
	BaseFields `mapstructure:",squash"`
}

type BaseFields struct {
	Name        string            `mapstructure:"name" desc:"Name for the target"`
	Description string            `mapstructure:"desc" desc:"Target description"`
	Labels      []string          `mapstructure:"labels" desc:"Labels to apply to the targets"` //
	Deps        []string          `mapstructure:"deps" desc:"Build dependencies"`
	PassEnv     []string          `mapstructure:"pass_env" desc:"List of environment variable names that will be passed from the OS environment, they are part of the target hash"`
	SecretEnv   []string          `mapstructure:"secret_env" desc:"List of environment variable names that will be passed from the OS environment, they are not used to calculate the target hash"`
	Env         map[string]string `mapstructure:"env" desc:"Key-Value map of static environment variables to be used"`
	Tools       map[string]string `mapstructure:"tools" desc:"Key-Value map of tools to include when executing this target. Values can be references"`
	Visibility  []string          `mapstructure:"visibility" desc:"List of visibility for this target"`
}

func (bf *BuildFields) GetBuildMods() []TargetOption {
	return bf.GetBaseMods()
}

func (bf *BaseFields) GetBaseMods() []TargetOption {
	return []TargetOption{
		WithVisibility(bf.Visibility),
		WithDescription(bf.Description),
		WithLabels(bf.Labels),
		WithPassEnv(bf.PassEnv),
		WithSecretEnvVars(bf.SecretEnv),
		WithEnvVars(bf.Env),
		WithTools(bf.Tools),
	}
}

type DeployFields struct {
	Environments map[string]*environs.Environment `mapstructure:"environments" desc:"Deployment Environments"`
}

func (df *DeployFields) GetDeployMods() []TargetOption {
	return []TargetOption{
		WithEnvironments(df.Environments),
	}
}
