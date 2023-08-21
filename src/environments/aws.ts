import { EnvironmentConfig } from './index'

export class AwsAuthenticationConfig implements EnvironmentConfig<AwsAuthenticationConfig> {
  profile: string | undefined
  assumeRole: string | undefined
  account: string = ""
  region: string = ""

  constructor(config: {
    account: string,
    region: string,
    profile: string | undefined,
    assumeRole: string | undefined
  }) {
    Object.assign(this, config)
  }

  env(): { [key: string]: string } {
    var env: { [key: string]: string } = {
      AWS_REGION: this.region,
      AWS_DEFAULT_REGION: this.region,
      AWS_ACCOUNT_ID: this.account
    }

    if (this.profile !== undefined) {
      env["AWS_PROFILE"] = this.profile
      env["AWS_DEFAULT_PROFILE"] = this.profile
    }

    return env
  }

  merge(...envs: AwsAuthenticationConfig[]): void {
    for (let env of envs) {
      this.account = env.account
      this.region = env.region

      if (env.profile !== undefined) {
        this.profile = env.profile
      }

      if (env.assumeRole !== undefined) {
        this.assumeRole = env.assumeRole
      }
    }
  }
}
