package environments

import (
	"github.com/zen-io/zen-core/utils"

	"golang.org/x/exp/slices"
)

func MergeEnvironments(envs ...*Environment) *Environment {
	finalEnv := &Environment{
		Aws:        &AwsAuthenticationConfig{},
		Kubernetes: &K8sAuthenticationConfig{},
		Variables:  make(map[string]string),
	}

	for _, e := range envs {
		if e == nil {
			continue
		}

		if e.Aws != nil {
			finalEnv.Aws.Merge(e.Aws)
		}

		if e.Kubernetes != nil {
			finalEnv.Kubernetes.Merge(e.Kubernetes)
		}

		finalEnv.Variables = utils.MergeMaps(finalEnv.Variables, e.Variables)
	}

	return finalEnv
}

func MergeEnvironmentMaps(envMaps ...map[string]*Environment) map[string]*Environment {
	mergedEnvs := map[string]*Environment{}
	allKeys := []string{}
	for _, m := range envMaps {
		for env := range m {
			if !slices.Contains(allKeys, env) {
				allKeys = append(allKeys, env)
			}
		}
	}

	for _, key := range allKeys {
		toMerge := make([]*Environment, 0)
		for _, m := range envMaps {
			toMerge = append(toMerge, m[key])
		}

		mergedEnvs[key] = MergeEnvironments(toMerge...)
	}

	return mergedEnvs
}
