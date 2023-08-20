package environments

import (
	"github.com/zen-io/zen-core/utils"
)

type EnvironmentConfig[T any] interface {
	EnvVars() map[string]string
	Merge(T)
}

type Environment struct {
	Aws        *AwsAuthenticationConfig `mapstructure:"aws"`
	Kubernetes *K8sAuthenticationConfig `mapstructure:"kubernetes"`
	Variables  map[string]string        `mapstructure:"variables"`
}

func (envConfig *Environment) Env() map[string]string {
	awsEnv := make(map[string]string)
	k8sEnv := make(map[string]string)

	if envConfig.Aws != nil {
		awsEnv = envConfig.Aws.EnvVars()
	}

	if envConfig.Kubernetes != nil {
		k8sEnv = envConfig.Kubernetes.EnvVars()
	}

	return utils.MergeMaps(awsEnv, k8sEnv, envConfig.Variables)
}
