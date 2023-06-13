package environments

type K8sAuthenticationConfig struct {
	Context *string `mapstructure:"context"`
	Config  *string `mapstructure:"config"`
}

func (k8s *K8sAuthenticationConfig) EnvVars() map[string]string {
	envVars := map[string]string{}

	if k8s.Context != nil {
		envVars["HELM_KUBECONTEXT"] = *k8s.Context
		envVars["KUBECTX"] = *k8s.Context
	}

	if k8s.Config != nil {
		envVars["KUBECONFIG"] = *k8s.Config
		envVars["KUBE_CONFIG_PATH"] = *k8s.Config
	}

	return envVars
}

func (dest *K8sAuthenticationConfig) Merge(src *K8sAuthenticationConfig) {
	if src.Context != nil {
		dest.Context = src.Context
	}
	if src.Config != nil {
		dest.Config = src.Config
	}
}
