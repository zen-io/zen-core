import { EnvironmentConfig } from './index'

export class K8sAuthenticationConfig implements EnvironmentConfig<K8sAuthenticationConfig> {
  context: string | undefined
  kubeconfig: string | undefined

  constructor(config: {
    context?: string, kubeconfig?: string
  }) {
    this.context = config.context
    this.kubeconfig = config.kubeconfig
  }

  env(): { [key: string]: string } {
    var env: { [key: string]: string } = {}
    if (this.context !== undefined) {
      env["HELM_KUBECONTEXT"] = this.context
      env["KUBECTX"] = this.context
    }

    if (this.kubeconfig !== undefined) {
      env["KUBECONFIG"] = this.kubeconfig
      env["KUBE_CONFIG_PATH"] = this.kubeconfig
    }

    return env
  }

  merge(...envs: K8sAuthenticationConfig[]): void {
    for (let env of envs) {
      if (env.context !== undefined) {
        this.context = env.context
      }

      if (env.kubeconfig !== undefined) {
        this.kubeconfig = env.kubeconfig
      }
    }
  }
}
