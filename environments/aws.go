package environments

type AwsAuthenticationConfig struct {
	Profile    *string `mapstructure:"profile"`
	AssumeRole *string `mapstructure:"assume_role"`
	Account    string  `mapstructure:"account"`
	Region     string  `mapstructure:"region"`
}

func (aws *AwsAuthenticationConfig) EnvVars() map[string]string {
	envVars := map[string]string{}

	envVars["AWS_REGION"] = aws.Region
	envVars["AWS_DEFAULT_REGION"] = aws.Region
	envVars["AWS_ACCOUNT_ID"] = aws.Account

	if aws.Profile != nil {
		envVars["AWS_PROFILE"] = *aws.Profile
		envVars["AWS_DEFAULT_PROFILE"] = *aws.Profile
	}

	return envVars
}

func (dest *AwsAuthenticationConfig) Merge(src *AwsAuthenticationConfig) {
	if src.Profile != nil {
		dest.Profile = src.Profile
	}
	if src.AssumeRole != nil {
		dest.AssumeRole = src.AssumeRole
	}
	dest.Account = src.Account
	dest.Region = src.Region
}
