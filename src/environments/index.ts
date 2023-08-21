import { mergeMaps } from "../utils/maps";
import { AwsAuthenticationConfig } from "./aws";
import { K8sAuthenticationConfig } from "./k8s";

export interface EnvironmentConfig<T> {
  env: () => { [key: string]: string };
  merge: (...envs: T[]) => void;
}

export class Environment implements EnvironmentConfig<Environment> {
  aws: AwsAuthenticationConfig | undefined
  k8s: K8sAuthenticationConfig | undefined
  variables: { [key: string]: string }

  constructor(config: {
    aws?: AwsAuthenticationConfig,
    k8s?: K8sAuthenticationConfig,
    variables?: { [key: string]: string }
  }) {
    this.aws = config.aws
    this.k8s = config.k8s

    if (config.variables === undefined) {
      config.variables = {}
    }
    this.variables = config.variables
  }

  env(): { [key: string]: string } {
    var awsEnv: { [key: string]: string } = {}
    var k8sEnv: { [key: string]: string } = {}

    if (this.aws !== undefined) {
      awsEnv = this.aws.env()
    }

    if (this.k8s !== undefined) {
      k8sEnv = this.k8s.env()
    }

    return mergeMaps(awsEnv, k8sEnv, this.variables)
  }

  merge(...environments: Environment[]): void {
    for (let env of environments) {
      if (env === undefined) {
        continue
      }

      if (env.aws !== undefined) {
        if (this.aws === undefined) {
          this.aws = env.aws
        } else {
          this.aws.merge(env.aws)
        }
      }

      if (env.k8s !== undefined) {
        if (this.k8s === undefined) {
          this.k8s = env.k8s
        } else {
          this.k8s.merge(env.k8s)
        }
      }

      this.variables = mergeMaps(this.variables, env.variables)
    }
  }
}
